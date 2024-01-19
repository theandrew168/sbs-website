---
date: 2024-01-18
title: "Quick Tip: Using Make to Run Concurrent Services"
slug: "quick-tip-using-make-to-run-concurrent-services"
tags: ["make"]
draft: true
---

This short post explains a useful trick for using [Make](<https://en.wikipedia.org/wiki/Make_(software)>) to run multiple development services at the same time.
For example, one of my recent web projects required three active services in order to develop locally:

1. Running the [Go](https://go.dev/) backend
2. Running [esbuild](https://esbuild.github.io/) to (re)build the React+TypeScript frontend
3. Running the [Tailwind CSS CLI](https://tailwindcss.com/blog/standalone-cli) to re(build) the Tailwind CSS styles

Historically, I'd run each of these programs in a separate terminal window (or use something fancy like [tmux](https://github.com/tmux/tmux/wiki)).
This is a fairly straightforward solution but I always found it somewhat clunky: alt-tabbing between terminals to restart or stop the services.

# The Trick

If you are using Make to manage your project, you can utilize its builtin concurrency to run each service in parallel.
With only a single terminal and one command, you can start and stop all services effortlessly.
This works as long as you allow Make to run simulataneous "jobs" via the `-j` flag.
You also need to structure your Makefile such that each development service has its own target and that these targets are mutually independent (they don't depend on each other directly or indirectly).

Here is an example from the aforementioned web project:

```make
.PHONY: run-frontend-js
run-frontend-js:
	# build and bundle the frontend whenever files change
	npx esbuild --watch ...

.PHONY: run-frontend-css
run-frontend-css:
	# build Tailwind CSS styles whenever files change
	npx tailwindcss --watch ...

.PHONY: run-frontend
run-frontend: run-frontend-js run-frontend-css

.PHONY: run-backend
run-backend:
	# run the backend API web server
	go run main.go

.PHONY: run
run: run-frontend run-backend
```

When these targets are invoked via `make -j run`, the following chain of events occurs:

1. `run` invokes `run-frontend` and `run-backend`
2. `run-frontend` invokes `run-frontend-js` and `run-frontend-css`
3. `run-frontend-js` starts the esbuild service
4. `run-frontend-css` starts the tailwindcss service
5. `run-backend` starts the backend web service

Within a few seconds, all services will be running simultaneously!

Note that without concurrency, Make would execute until it encounters the first service and try to wait for it to finish.
However, since a service never truly finishes execution, none of the others would ever get a chance to start.

# Conclusion

Ever since discovering this trick, I use it everywhere.
For my personal develop workflow, it "just works".
I no longer have to fuss with a bunch of terminals and trying to remember which is which.
If you are a fellow Make user ([or would like to join the club](/posts/implementing-make-in-go/)), give this trick a shot!
