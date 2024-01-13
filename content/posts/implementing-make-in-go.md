---
date: 2024-01-11
title: "Implementing Make in Go"
slug: "implementing-make-in-go"
draft: true
---

https://github.com/theandrew168/make

A while back, my buddy and I went on a weekend trip to a remote cabin in eastern Iowa.
We've done this before: with plans to fish, hike, and just enjoy the quietness of nature and the bliss of having no responsibilities for a few days.
Despite being an outdoorsy trip, I still brought my laptop because I love to program!
Internet is usually spotty / non-existent in the areas we stay so I still consider it to be a nice disconnect.

Part of my prep for this trip was coming up with a narrowly-scoped, interesting project that I could complete (or nearly complete) within its bounds.
As a long time fan of [Make](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html), I thought it'd be interesting to implement some useful subset of its behavior.
Additionally, there isn't much Make support on Windows so writing it in an easily cross-platform language would be a boon.
So, I decided to write Make in [Go](https://go.dev/)!

# What is Make?

What is (POSIX) Make?
Why do I like Make so much?
Some only use it for C/C++ but I think it works great for any software project with tasks and deps.

Here is a small example Makefile:

```make
hello:
	echo 'hello'

world: hello
	echo 'world'
```

This file defines two targets: `hello` and `world`.
The `hello` target has no dependencies while the `world` targets depends on `hello`.
If you just want to run the `hello` target, nothing extra is printed:

```
$ make -s hello
hello
```

However, if you run the `world` target, Make will run `hello` first because it is a dependency:

```
$ make -s world
hello
world
```

From this simple foundation, an incredibly useful toolbox of targets can be built.
When expressed with proper dependencies, you never have to worry about running commands in the wrong order.
Combine that value with the speed of executing non-dependent targets in parallel and you've got a perfect project assistant.

# The Code

All of the code is written in a [single file](https://github.com/theandrew168/make/blob/main/make.go).
It really wasn't complex enough to warrant more than that (though I didn't know that at the start).
I simply started writing and got it working before things got out of hand.

### Data Structures

The first important piece of data is a `Target`.
This is a task defined with a Makefile that has dependencies (other targets ref's by name) and commands (things to execute on the command line).
Since I'm writing in Go, I defined `Target` as a struct:

```go
type Target struct {
	// more about this later...
	sync.Once

	Prerequisites []string
	Commands      []string
}
```

The second is the `Graph` which isn't really as much of a graph as it is a mapping of names to targets:

```go
type Graph map[string]*Target
```

### Lexing / Parsing

The code starts by reading the specified file (defaults to `Makefile`) line by line.
I used Go's [bufio.Scanner](https://pkg.go.dev/bufio#Scanner) for this.
Comments, dot-directives, and empty lines are ignored.
Then, anything line that _doesn't_ start with a TAB is assumed to be the start of a new target definition.
Subsequent lines that start with a TAB are added to the actions of the current target.

### Execution

I originally thought I'd need to topologically sort the targets in order to determine a correct ordering.
In practice, however, a simple recursion that enforces each target's deps run before itself accomplishes the same thing.
Once all of a Target's deps have ran successfully, its own actions will run.

### Concurrency

The execution of each target uses a `sync.WaitGroup` to ensure all dependencies have ran before moving onto its own actions.
If any of the dependencies return an error, then execution stops and the error is propagated without running this target's actions.
Since one of the deps hit an error, it'd be invalid to execute the current target.

# Missing Features

When comparing this small program to the full featureset of POSIX make, I'd say _most_ features are missing.
Many examples come to mind: variables, non-phony targets, default rules.
And that is only for POSIX make.
Other Make imples (like GNU-Make) greatly example the available features to include additional variable-related features, conditional logic, and more.

# Lessons Learned

This was one of those rare occurrences where a project ends up being _less_ work than you expect.
Granted, I only implemented a minimal subset of Make and achieving full POSIX compat would take much longer.
As mentioned earlier, I really did expect to need to build and topo sort the DAG.
I learned that a simple depth-first recursion plus ensuring each target only runs once does the trick.

Go is awesome, especially when it comes to concurrency and cross-platform support.
I guess this isn't a "lesson learned" as much as it is a "lesson reinforced".
Overall, this was a great mini-project and I enjoyed building it.
I'd like to continue interating on it sometime (maybe I need to take more cabin trips).
Achieving full POSIX-compat would be a solid milestone to work toward.

Thanks for reading!
