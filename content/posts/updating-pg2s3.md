---
date: 2024-05-04
title: "Updating pg2s3"
slug: "updating-pg2s3"
draft: true
---

Hasn’t really been touched / updated in years. Update to the latest Go version (1.22), PGX v5, and all other deps. Maybe even clean up the code a bit if it needs it. Also talk about what it does, how it works, and why I built it.

# Overview

What is it?
Why did I build it?
Is it being used in production?

# Updates

### Go 1.22

No changes needed.

### go-systemd v22.5.0

No changes needed.

### gocron v2.4.0

Update `NewScheduler` to use [functional options](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis).
Use the updated `NewJob` API (more functional options).
Add my own "block til exit" since `StartBlocking` was removed.

### pgx v5.5.5

No changes needed.
The pgx library is only used to validate the PostgreSQL connection.
The actual exporting / importing is done by calling out to `pg_dump` and `pg_restore`, respectively.
