---
date: 2024-06-30
title: "Two Go + PostgreSQL Timestamp Gotchas"
slug: "two-go-plus-postgresql-timestamp-gotchas"
tags: ["Databases", "Go"]
draft: true
---

PG’s timestamptz in UTC doesn't really make sense.
The time will be stored without any TZ but converted to local when read.
Therefore, any time could go in as UTC and come out with a local TZ.
Instead, just use timestamp.

Go’s timestamps are in nanoseconds but PG’s are in microseconds.
Therefore, you lose precision when inserting and are better off rounding to micros beforehand.
