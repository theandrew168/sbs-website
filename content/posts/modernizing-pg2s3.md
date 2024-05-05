---
date: 2024-05-04
title: "Modernizing pg2s3"
slug: "modernizing-pg2s3"
draft: true
---

Hasn’t really been touched / updated in years. Update to the latest Go version (1.22), PGX v5, and all other deps. Maybe even clean up the code a bit if it needs it. Also talk about what it does, how it works, and why I built it.

# Overview

What is it?
Why did I build it?
Is it being used in production?

# Modernizations

### Logging

Current code just using `fmt.Print` for logging.
Use proper logging package.

### Running

Have the `func run() error` function return an error instead of an `int`.
Then `func main()` can call `run`, check for and log an errors, and exit with an appropriate code.

### Services

Instead of relying on gocron's `StartBlocking`, use my standard `Service.Run` interface w/ a `SignalContext` and `WaitGroup`.
