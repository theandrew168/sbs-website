---
date: 2024-01-14
title: "Implementing Make in Go"
slug: "implementing-make-in-go"
tags: ["make", "go"]
draft: true
---

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
The commands for a given target must be indented with a hard tab character (this a requirement of Make's syntax).

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

The source code for this project can be found on [GitHub](https://github.com/theandrew168/make).

### Parsing

While not a full-blown lexer/parser, the subset of Make's syntax that I chose to support can be handled with line-based processing.
The code starts by reading the user-supplied Make file line by line (using Go's [bufio.Scanner](https://pkg.go.dev/bufio#Scanner)).
Any non-empty line that is a comment (starts with `#`) or a dot-directive (starts with `.`) is ignored.
Then, anything line that _doesn't_ start with a tab is assumed to be the start of a new target definition.
Subsequent lines that start with a tab are added to the list of commands for the current target.

### Data Structures

What is a "target", though?
A target is an entry within a Makefile that has commands (things to execute on the command line) and dependencies (other targets referenced by name).
Remember that the example used earlier had two targets: `build` (no dependencies) and `run` (depends on `build`).

Since I'm using Go, I represented a target as a struct:

```go
type Target struct {
	Dependencies []string
	Commands     []string
}
```

The second data type worthy of explanation is the `Graph`.
In practice, though, this isn't really much of a graph but instead a mapping of names to targets.
Dependency info is stored within each target but we still need a way to lookup targets by name.

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

Talk about [POSIX make](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html) vs GNU make.
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
