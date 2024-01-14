---
date: 2024-01-14
title: "Implementing Make in Go"
slug: "implementing-make-in-go"
tags: ["make", "go"]
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

Before writing the code, I presumed that I'd need to [topologically sort](https://en.wikipedia.org/wiki/Topological_sorting) the targets in order to determine a correct execution order.
In practice, however, **a simple recursion that enforces that the current target's dependencies run before itself did the trick**.
Once all of a target's dependencies have been executed successfully, its own commands will run.

In simplified pseudocode, this is the execution logic:

```python
def execute(graph, name):
	# lookup the current target by name
	target = graph[name]

	# recursively execute the current target's dependencies
	for dependency in target.dependencies:
		execute(graph, dependency)

	# run the current target's commands
	for command in target.commands:
		run_in_cli(command)
```

In reality, the execution code is a bit more complex than that due to a couple nuances: error handling and concurrency.

### Error Handling

There are only a couple ways that the execution of a given target can go wrong but it is important to acknowledge them.
Failing to detect or handle the errors correctly could lead to an invalid execution of the graph and break the crucial invariant: a target's dependencies must **all execute successfully** before its own commands can run.

The first possible error arises if a dependency doesn't exist: its name isn't present in the graph.
This error can occur in malformed Makefiles.
Although the existence of each dependency _could_ be validated prior to execution, I chose perform this check during execution.

A second error scenarios is if a command runs but causes an error.
A command in a Makefile can be any arbitrary string provided by the user.
This opens the door to many possible issues: typos, missing programs, invalid args, etc.
I'm using Go's [os/exec](https://pkg.go.dev/os/exec) package to run the external commands, check for errors, and capture any output.
Since I execute each dependency concurrently, I need a way to track errors from asynchronous execution and prevent further execution if any arise.

### Concurrency

The execution of each target uses a [sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup) to ensure all dependencies have ran before moving onto its own commands.
Additionally, Go's [sync.Once](https://pkg.go.dev/sync#Once) utility is used to ensure that a target's commands only run once regardless of how many goroutines are competing.
A simple modification gives the `Target` struct this new superpower:

```go
type Target struct {
	// embed sync.Once to enable "execute once" behavior
	sync.Once

	Dependencies []string
	Commands     []string
}
```

Executing a target's commands now requires calling `.Do` but gives us safety against race conditions:

```go
func executeCommands(target *Target) error {
	var commandErr error
	target.Do(func() {
		// execute commands via os/exec
	})
	return commandErr
}
```

Any errors returned from the execution of a dependency are delivered to the current execute's context via a channel.
If any dependencies return an error, then execution stops and the error is propagated without running the current target's commands.

Once all deps are kicked off, a `select` is used to wait for either: the success of **all** deps or the failure of **any** deps.
This is the full, unsimplifed execution logic in Go:

```go
func execute(graph Graph, name string) error {
	// lookup current target by name
	target, ok := graph[name]
	if !ok {
		return fmt.Errorf("target does not exist: %s", name)
	}

	// create a channel for receiving errors from dependencies
	errors := make(chan error)

	// recursively execute all dependencies
	var wg sync.WaitGroup
	for _, dependency := range target.Dependencies {
		dependency := dependency

		wg.Add(1)
		go func() {
			defer wg.Done()

			// submit any dependency errors to the errors channel
			err := execute(graph, dependency)
			if err != nil {
				errors <- err
			}
		}()
	}

	// turn wg.Wait() into a select-able channel
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// wait for dependencies to finish / check for errors
	select {
	case <-done:
	case err := <-errors:
		return err
	}

	// execute the current target's commands
	return executeCommands(target)
}
```

# Lessons Learned

In many aspects, this project was a success: I was able implement a useful Make clone in a short cabin trip!
It was one of those rare occurrences where a project ended up being _less_ work than expected.
As mentioned earlier, I thought that I would have to fully build and topologically sort the graph before execution.
I learned that using a depth-first recursion and enforcing that each target's dependencies run first does the trick.

Despite implementing a personally-useful subset of Make behavior, there are many features that this project is missing.
The [Make specification](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html) outlines a number of useful things that my version lacks: variables, non-phony targets, default rules, etc.
I'd like to continue iterating on this project eventually.
Achieving full compatibility would be a solid milestone to work toward.

Another lesson: Go is awesome, especially when it comes to concurrency!
I guess this isn't a "lesson learned" as much as it is a "lesson reinforced".
The tools that Go provides in its standard library are comprehensive and well-designed.
A vast scope of projects can be built that depend only on Go's builtin packages.

Overall, this was a great mini-project and I enjoyed building it.
Thanks for reading!
