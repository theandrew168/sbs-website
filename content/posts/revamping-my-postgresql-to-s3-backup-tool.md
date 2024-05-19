---
date: 2024-05-12
title: "Revamping My PostgreSQL to S3 Backup Tool"
slug: "revamping-my-postgresql-to-s3-backup-tool"
tags: ["Databases", "Go"]
draft: true
---

Hasn’t really been touched / updated in years. Update to the latest Go version (1.22), PGX v5, and all other deps. Maybe even clean up the code a bit if it needs it. Also talk about what it does, how it works, and why I built it.

# Overview

What is it?
Why did I build it?
Is it being used in production?

# Updates

### go v1.22

No changes needed.

### go-systemd v22.5.0

No changes needed.

### gocron v2.4.0

Update `NewScheduler` to use [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
Use the updated `NewJob` API (more functional options).
Add my own "block til exit" since `StartBlocking` was removed.
Instead of relying on gocron's `StartBlocking`, wait for a `SignalContext`.

### pgx v5.5.5

No changes needed.
The pgx library is only used to validate the PostgreSQL connection.
The actual exporting / importing is done by calling out to `pg_dump` and `pg_restore`, respectively.

# Refactors

### Logging

Current code just using `fmt.Print` for logging.
Use proper logging package.

### Running

Have the `func run() error` function return an error instead of an `int`.
Then `func main()` can call `run`, check for and log an errors, and exit with an appropriate code.
