---
date: 2024-12-08
title: "Uploading 50GB of Music to S3"
slug: "uploading-50gb-of-music-to-s3"
draft: true
---

For DerzTunes, I wanna write an iTunes-style web application for listening to my old school collection of music files.
Part of the project involves uploading my 50GB of music from old, external hard drives up to an S3 bucket (DO Space).
What’s the best way to do that in an idempotent, restart-able way?
Can rsync do that?
The minio client?
Or the AWS CLI?

Options:

1. AWS CLI (the de-facto, first-party choice)
2. s3cmd (recommended by Digital Ocean)
3. mc (minio client, I'm a fan of minio)

My first choice, the AWS CLI, ended up working perfectly!
Initially, I just did a recursive copy.
To verify, I ran a sync.
If my collection was still growing, I'd continue to use sync to ensure everything lines up.
Using `cp --recursive` is simply an optimization for the "first time" upload since we know all the files are new and it isn't worth wasting time verifying which files are missing.

https://stackoverflow.com/questions/65321150/aws-speed-up-copy-of-large-number-of-very-small-files

```
aws s3 cp --recursive --endpoint=https://nyc3.digitaloceanspaces.com . s3://derztunes
aws s3 sync --endpoint=https://nyc3.digitaloceanspaces.com . s3://derztunes
```
