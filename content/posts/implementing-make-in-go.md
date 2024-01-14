---
date: 2024-01-14
title: "Implementing Make in Go"
slug: "implementing-make-in-go"
draft: true
tags: ["make", "go"]
---

https://github.com/theandrew168/make

https://en.wikipedia.org/wiki/Make_(software)

https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html

A while back, my buddy [Wes](https://brue.land/) and I took a weekend trip to a remote cabin in eastern Iowa.
We try to do this once a year with plans of fishing, hiking, and simply enjoying the quietness of nature.
Not to mention the bliss of having no responsibilities for a few days!
Despite being an outdoorsy trip, I still brought my laptop because programming is my favorite hobby.
The internet in such remote locations is typically unreliable so I still consider it to be a nice disconnect.

Part of my prep for this trip was coming up with a narrowly-scoped, interesting project that I could complete (in some minimally-viable form) within its duration.
As a longtime fan of [Make](<https://en.wikipedia.org/wiki/Make_(software)>), I thought it'd be interesting to implement some useful subset of its behavior.
Additionally, there isn't much Make support on Windows so writing it in an easily cross-platform language would be a boon.
So, I decided to implement Make in [Go](https://go.dev/)!

# What is Make?

[Make](<https://en.wikipedia.org/wiki/Make_(software)>) is a simple build automation tool for expressing arbitrary command-line actions and their dependencies on other commands.
While Make has historically had a close relationship with C/C++ development, I have found that it adds value to any software project that has various build steps and tasks (running tests, formatting code, etc).
In practice, this ends up being most projects!
I've personally used Make to manage Go-based tools, Python projects, and polyglot web applications.

Here is a small example of using Make to manage building and running a Go program:

```make
build:
	go build -o main main.go
run: build
	./main
```

This file defines two targets: `build` and `run`.
The `build` target has no dependencies while the `run` targets depends on `build`.
Make can be invoked on the CLI by typing `make` followed by the target you wish to execute.

If you want to run the `build` target, then Make will run the command provided to build the Go program:

```sh
$ make build
go build -o main main.go
```

However, if you choose the `run` target, Make will execute `build` first (because it is declared as a dependency) and _then_ run the program:

```sh
$ make run
go build -o main main.go
./main
Hello World!
```

From this simple foundation, an incredibly useful toolbox of targets can be built.
When expressed with proper dependencies, you never have to worry about running commands in the wrong order.
Combine that value with the speed of executing non-dependent targets in parallel and you've got a perfect project assistant.

# The Code

All of the code is written in a [single file](https://github.com/theandrew168/make/blob/main/make.go).
It really wasn't complex enough to warrant more than that (though I didn't know that at the start).
I simply started writing and got it working before things got out of hand.

### Parsing

The code starts by reading the specified file (defaults to `Makefile`) line by line.
I used Go's [bufio.Scanner](https://pkg.go.dev/bufio#Scanner) for this.
Comments, dot-directives, and empty lines are ignored.
Then, anything line that _doesn't_ start with a TAB is assumed to be the start of a new target definition.
Subsequent lines that start with a TAB are added to the actions of the current target.

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

### Execution

I originally thought I'd need to topologically sort the targets in order to determine a correct ordering.
In practice, however, a simple recursion that enforces each target's deps run before itself accomplishes the same thing.
Once all of a Target's deps have ran successfully, its own actions will run.

### Concurrency

The execution of each target uses a `sync.WaitGroup` to ensure all dependencies have ran before moving onto its own actions.
If any of the dependencies return an error, then execution stops and the error is propagated without running this target's actions.
Since one of the deps hit an error, it'd be invalid to execute the current target.

# Missing Features

Talk about POSIX make vs GNU make.
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
