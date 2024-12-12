---
date: 2024-12-22
title: "Prepared Statements: An Easy Win?"
slug: "prepared-statements-an-easy-win"
tags: ["Databases"]
draft: true
---

Some PostgreSQL database drivers transparently optimize your application by utilizing prepared statement to cache your queries (Go’s PGX and Java’s JDBC).
When implemented correctly, the user gets a faster application essentially for free.

However, some libraries don’t do this (like every npm PG driver).
Why not?
Is there a reason that this was intentionally left out?
Or maybe the package developers weren’t aware of the feature (or how transparently it could be implemented) and pursued an architecture that now doesn’t support it.
