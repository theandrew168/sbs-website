---
date: 2024-04-21
title: "Limiting Concurrency with Semaphores"
slug: "limiting-concurrency-with-semaphores"
tags: ["Go"]
draft: true
---

I recently encountered a scenario where the serial execution of tasks within a program was performing slower than I'd like.
Specifically, I was working on how [Bloggulus](https://bloggulus.com/) syncs all of the blogs that it tracks (there are only 40 or so right now but I expect the number to grow).
The program evolved through a few different approaches before arriving at a balanced solution.
This post outlines and explains that journey.

# Serial

As a baseline, consider this simple program that executes multiple tasks in sequence.
Each job must be performed one after another until all have completed.
Instead of bogging you down the details of reading RSS feeds, I've simulated some work by sleeping for a quarter second before printing "job done!":

```go
package main

import (
	"fmt"
	"time"
)

func doWork() {
	time.Sleep(250 * time.Millisecond)
	fmt.Println("job done!")
}

func main() {
	for i := 0; i < 16; i++ {
		doWork()
	}
}
```

TODO: Add vis

Notice how the program takes roughly four seconds to execute (16 jobs \* 0.25 second per job).
There is nothing fancy going on here: each job executes linearly.

# Concurrent

Surely, we can do better than that, right?
This is Go, after all, and Go has goroutines!
We can just throw some `go` keywords in front of our `doWork` function and we'll be off to the races!

```go
package main

import (
	"fmt"
	"time"
)

func doWork() {
	time.Sleep(250 * time.Millisecond)
	fmt.Println("job done!")
}

func main() {
	for i := 0; i < 16; i++ {
		go doWork()
	}
}
```

TODO: Add vis

Wait, that doesn't look right.
The program didn't print anything!
What happened to our jobs?
The problem here is that Go's runtime doesn't wait for all goroutines to finish before exiting the program.
This means that our jobs didn't even get a chance to run because the program terminated before they could even get going.
How can we tell the program to wait for our jobs to complete?

# WaitGroup

Thankfully, Go's standard library holds the solution: [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup).
From the docs:

> A WaitGroup waits for a collection of goroutines to finish.

That sounds perfect!
How does it work?

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func doWork() {
	time.Sleep(250 * time.Millisecond)
	fmt.Println("job done!")
}

func main() {
	// initialize a WaitGroup
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		// add a counter to the WaitGroup
		wg.Add(1)

		go func() {
			// remove a counter from the WaitGroup after doWork completes
			defer wg.Done()
			doWork()
		}()
	}

	// wait for all jobs to finish
	wg.Wait()
}
```

TODO: Add vis

There we go, all the jobs executed and completed in roughly a quarter second!
In this example, we've added a `sync.WaitGroup` and ensured that each job increments the group's counter when it starts and decrements the group's counter when it finishes.
Then, at the end of the program, we call `wg.Wait()` to wait for all running jobs to finish.
Pretty neat!

There is one small problem, though, in the scenario I was facing.
I don't actually want to sync every blog at the same exact time.
With this approach, if [bloggulus](https://bloggulus.com/) was tracking 500 blogs, the sync process would blast the network with 500 outgoing requests at once!
This might not be an issue in practice but I'd rather find a way to smooth out the network traffic.
Maybe there is a way to put an upper limit on the number of simultaneous syncs?

# Semaphore

This time, Go's _extended_ standard library holds the solution: [semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore).

```go
package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/semaphore"
)

// limit the number of concurrent workers
const MaxWorkers = 4

func doWork() {
	time.Sleep(250 * time.Millisecond)
	fmt.Println("job done!")
}

func main() {
	// initialize a semaphore
	sem := semaphore.NewWeighted(MaxWorkers)
	for i := 0; i < 16; i++ {
		// acquire a single counter from the semaphore
		// (this blocks if all counters are in use)
		sem.Acquire(context.Background(), 1)

		go func() {
			// release a counter to the semaphore after doWork completes
			defer sem.Release(1)
			doWork()
		}()
	}

	// wait for all jobs to finish
	sem.Acquire(context.Background(), MaxWorkers)
}
```

TODO: Add vis

This is exactly what I'm after: the best of both worlds!
The jobs execute concurrently but only `MaxWorkers` are able to run at the same time.
With this approach, I can limit how many simultaneous requests Bloggulus makes and prevent clogging up the network.
Check out the final implementation [on GitHub](https://github.com/theandrew168/bloggulus/blob/981424b37cee14a13f4caec556bcc3042260ab37/backend/service/sync.go#L89-L116).

# Conclusion

This post walked through a few basic examples of how Go's concurrency can be used to speedup a program's execution while limiting the number of active goroutines.
Overall, I'm happy with how readable the final example is despite utilizing moderately-complex concurrency ideas.
It goes to show how well designed the Go programming language is!
Thanks for reading.
