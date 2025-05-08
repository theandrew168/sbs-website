---
date: 2024-12-08
title: "Uploading 50GB of Music to S3"
slug: "uploading-50gb-of-music-to-s3"
---

One of my upcoming projects for the [Year of Clojure](/posts/2025-the-year-of-clojure/) is to create a web-based, iTunes-style user interface for listening to my personal music collection.
I'm thinking something with an [old school aesthetic](https://www.versionmuseum.com/history-of/itunes-app) based on iTunes 6 or iTunes 7.
Part of this project involves hosting my ~50GB of music (~7500 files) somewhere that is easily accessed by the app.
To me, there was one obvious choice: S3-compatible object storage.
Since I use [Digital Ocean](https://www.digitalocean.com/) for infrastructure, this means utilizing [Spaces Object Storage](https://www.digitalocean.com/products/spaces).

Today, my music files live on an external USB hard drive (a Seagate FreeAgent from the mid-2010s).
The big question is this: what is the best way to get these files off of the hard drive and safely uploaded to S3?
My list of requirements is short:

1. The uploaded files should be private
2. The upload should be fast and concurrent
3. The upload should be idempotent and restartable

Let's start with creating the bucket.

## Infra

Creating Digital Ocean buckets [via Terraform](https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket) is super simple and something I've done before.
You really can't mess this up: give it a name, give it a region, and make it private!
Nice and easy.

```hcl
resource "digitalocean_spaces_bucket" "derztunes" {
  name   = "derztunes"
  region = "nyc3"
  acl    = "private"
}
```

This satisfies my first requirement by ensuring that files added to the bucket will be private.
Now for the actual upload!

## Upload

In my head, I told myself: "I want something like rsync but for S3".
To put this another way, I want something fast and idempotent.
If the upload were to get interrupted midway through, I don't want to start over from zero.
Instead, I want to use a program that can compare the local files to what exists in S3 and only upload _missing_ files.

Based on my personal experience and brief research, I was considering three different tools for the job.

1. [AWS CLI](https://aws.amazon.com/cli/) (the de-facto, first-party choice)
2. [S3cmd](https://s3tools.org/s3cmd) (recommended by Digital Ocean)
3. [MinIO Client](https://min.io/docs/minio/linux/reference/minio-mc.html) (I'm a fan of the [MinIO project](https://github.com/minio/minio) and its [client libraries](https://github.com/minio/minio-go))

I figured I'd start with the official AWS CLI since it is widely used and built to scale (I'm assuming).
The AWS CLI includes an rsync-style command called [sync](https://docs.aws.amazon.com/cli/latest/reference/s3/sync.html) which is exactly what I was looking for.
After installing the CLI and exporting the standard AWS credentials, performing the sync was straightforward:

```bash
export AWS_ACCESS_KEY_ID=$SPACES_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$SPACES_SECRET_ACCESS_KEY
aws s3 sync --endpoint=https://nyc3.digitaloceanspaces.com . s3://derztunes
```

Unsurprisingly, it ended up working perfectly!
The CLI first checks for any missing files and then uploads them concurrently (the max concurrency is [configurable](https://docs.aws.amazon.com/cli/latest/topic/s3-config.html#max-concurrent-requests)).
The process took roughly 90 minutes for the initial sync while subsequent "no-op" syncs took only a minute or two.

## Optimize

After the fact, I did learn of a [small optimization](https://stackoverflow.com/a/65321940).
When doing the initial sync (before any files have been uploaded), you can avoid the "comparing files" work by simply performing a recursive copy.
This unconditionally copies all files present in the source directory to the destination bucket.

```bash
export AWS_ACCESS_KEY_ID=$SPACES_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY=$SPACES_SECRET_ACCESS_KEY
aws s3 cp --recursive --endpoint=https://nyc3.digitaloceanspaces.com . s3://derztunes
```

In the future, I'll do this first instead of jumping right in to the sync.

Thanks for reading!
