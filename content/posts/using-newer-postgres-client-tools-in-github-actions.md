---
date: 2024-05-05
title: "Using Newer PostgreSQL Client Tools in GitHub Actions"
slug: "using-newer-postgres-client-tools-in-github-actions"
---

Recently, while updating my [pg2s3 utility](https://github.com/theandrew168/pg2s3), I noticed that the project's `docker-compose.yml` file was pinning PostgreSQL to version 14.
I couldn't remember why I did that, so I went ahead and removed it (what could go wrong?).
Unfortunately, this led to some automated tests failing in GitHub Actions!

# The Problem

Thankfully, the error was very clear:

```
pg_dump: error: server version: 16.2 (Debian 16.2-1.pgdg120+2); pg_dump version: 14.11 (Ubuntu 14.11-1.pgdg22.04+1)
pg_dump: error: aborting because of server version mismatch
```

Classic version mismatch.
The PostgreSQL 16 server (running in a container) was not compatible with the PostgreSQL 14 client tools (installed on the GitHub Actions runner).
For some context, pg2s3 uses `pg_dump` and `pg_restore` to quickly export and import data.

I had two options: keep the server container pinned to version 14 or figure out how to install and use newer client tools on the Actions runner.
I opted for the latter since it is the more correct and future-proof solution.

# The Solution

As it turns out, the [PostgreSQL docs](https://www.postgresql.org/download/linux/ubuntu/) include a section about installing newer versions on stable releases of Ubuntu.
The docs explain how and why stable Linux releases can fall behind:

> PostgreSQL is available in all Ubuntu versions by default. However, Ubuntu "snapshots" a specific version of PostgreSQL that is then supported throughout the lifetime of that Ubuntu version.
> The PostgreSQL project maintains an Apt repository with all supported of PostgreSQL available.

So, we just need to configure the Actions runner to use the [PostgreSQL Apt Repository](https://wiki.postgresql.org/wiki/Apt) and then install whichever version we need.
The instructions were pretty straightforward: install `postgresql-common` and then run the automation script for setting up the apt repository:

```
sudo apt install -y postgresql-common
sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y
```

Once that configuration is complete, all modern versions of PostgreSQL will be available for install.
All that remains is removing any existing (older) client versions and then installing the most recent tools:

```
sudo apt purge -y postgresql-client-common
sudo apt install -y postgresql-client
```

Now both the server container and client tools will both be using the latest version which resolves the `version mismatch` error.
With the versions aligned, my tests were all passing once again!

# GitHub Actions

This process can be bundled into a single GitHub Actions step for usage in any `ubuntu-latest`-based workflow:

```yaml
- name: Install latest PostgreSQL client tools
  run: |
    sudo apt install -y postgresql-common
    sudo /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y
    sudo apt purge -y postgresql-client-common
    sudo apt install -y postgresql-client
```

If you ever find yourself building a tool that integrates with PostgreSQL's client tools, this snippet might come in handy.
Thanks for reading!
