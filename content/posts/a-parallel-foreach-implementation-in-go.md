---
date: 2025-04-27
title: "A Parallel ForEach Implementation in Go"
slug: "a-parallel-foreach-implementation-in-go"
tags: ["Go"]
---

Go's concurrency makes it very easy to fun multiple pieces of code in parallel.
As I've [written about](/posts/limiting-concurrency-with-semaphores/) before, running many tasks at the same time (and waiting for them all to finish) is achievable using little more than the `go` keyword and a [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup).

However, a problem that you can run into is actually being _too_ parallel.
For example, if you needed to download hundreds or thousands of files, initiating the requests for all of them at once could overwhelm the server.
Instead, it'd be nice if we could say something like: "download these files concurrently but only up to N at once".

In my aforementioned post, I showed how this could be achieved quite simply by using a [semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore).
One thing I failed to include, though, was a nice way to abstract this behavior in a reusable helper function that **also handles errors**.
Well, time to fix that!

## ForEach

Below is a generic implementation of this "perform a task concurrently up to N at a time" behavior.
The function `ForEach` accepts a list of items, a task function to process each item, and a concurrency limit (expressed as a number of goroutines).

```go
func ForEach[T any](concurrency int, items []T, taskFn func(T) error) error {
	// Use an error group to limit concurrency.
	var g errgroup.Group
	g.SetLimit(concurrency)

	// Perform tasks in parallel (up to "concurrency" at once).
	for _, item := range items {
		// Ensure the proper item value is captured for Go versions earlier than 1.22.
		item := item
		// For each item, attempt to acquire a "task slot" from the error group.
		// If the maximum number of tasks are already running, this will block.
		g.Go(func() error {
			return taskFn(item)
		})
	}

	// Wait for all tasks to finish and return the first non-nil error (if any).
	return g.Wait()
}
```

Let's see how this helper can be used in a basic "download some flags" example:

```go
func main() {
	// List of flags to be downloaded.
	urls := []string{
		"https://flagsapi.com/US/flat/64.png",
		"https://flagsapi.com/GB/flat/64.png",
		"https://flagsapi.com/AU/flat/64.png",
		"https://flagsapi.com/AQ/flat/64.png",
		"https://flagsapi.com/BE/flat/64.png",
	}

	// Download all flags two at a time.
	err := ForEach(2, urls, func(url string) error {
		// Fetch the data for a single flag.
		body, err := GetURL(url)
		if err != nil {
			return err
		}

		// Simulate processing the response.
		time.Sleep(1 * time.Second)
		fmt.Printf("Fetched %d bytes from %s\n", len(body), url)
		return nil
	})

	// If an error occurred while fetching a flag, print it.
	if err != nil {
		fmt.Printf("Failed to download flag: %v\n", err)
	}
}
```

Pulling this behavior into a dedicated function definitely makes the calling code (which used to be a complex mixture of ideas) easier to understand.
More specifically, it de-couples the "how do I perform the task" logic from the "how do I run the tasks concurrently" machinery.

Additionally, any errors that arise (well, just the first error, techincally) will be propagated up to the caller.
Note that despite only returning a single error, the `ForEach` helper will still wait for all tasks to finish before returning ([reference](https://pkg.go.dev/golang.org/x/sync/errgroup#Group.Wait)).
This behavior is similar to JavaScript's [Promise.allSettled](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled) utility.

Anyhow, I hope you can find this helper and example useful.
At the very least, maybe you found it informative.
Thanks for reading!