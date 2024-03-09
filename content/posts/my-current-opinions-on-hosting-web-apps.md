---
date: 2024-02-18
title: "My Current Opinions on Hosting Web Apps"
slug: "my-current-opinions-on-hosting-web-apps"
tags: ["Go", "Hosting"]
---

A little while back, I was chatting with some tech friends about my experience using [Fly](https://fly.io/) and [Neon](https://neon.tech/) for hosting web apps that are under active development.
Between these two services and their low pricing for small projects, my total bill for the month was only $0.01 (Fly doesn't even collect invoices below $5.00).
I shared how Fly and Neon seemed like a great fit for projects that are a work in progress but that I'm not sure if I'd use them for production services.
I'd have to do more pricing estimations and stability testing before having the trust and confidence needed to rely on these tools for anything that had users or was generating income.

# The Question

Given that summary, my friend [Minh](https://github.com/minhio) asked a follow-up question:

> What would you recommend for production apps?
> Not for milliions of users, but say "a startup at the bootstrap stage".

I wrote a lengthy answer in Discord but wanted to share it here for visibility and preservation.
That being said, my opinions are still changing as I write more code, deploy more services, and experience the pros and cons of different strategies.

# My Answer

Honestly, I'm not sure since none of my personal projects have ever made it that far.
If I was using a stack that builds into a single binary (like a [Go](https://go.dev/)-based app), then I'd probably deploy it onto a standard VPS (virtual private server) via [Digital Ocean](https://www.digitalocean.com/) or [Linode](https://www.linode.com/) (I'm a bit biased against AWS but that's mostly subjective).
My current approach for deploying single-binary Go apps is to build them into a [deb package](<https://en.wikipedia.org/wiki/Deb_(file_format)>) via [GoReleaser](https://goreleaser.com/).
The deb handles user creation, systemd unit setup, etc.
Then I use [Ansible](https://www.ansible.com/) to perform basic server setup and install the deb.

For non single binary stacks (NodeJS, Django, and basically everything else), I prefer to deploy the app as a container in order to ensure all the deps / static files are located where the app expects them.
I've messed with deploying these stacks with Ansible but it never feels very clean.
You've gotta check out the code, ideally create a virtual environment, install the deps, and run the app from that virtual env.
Kinda fussy but I'm sure there is room for improvement in these types of apps on traditional servers.

That being said, I've only messed with Fly for deploying containerized hobby projects.
If I was going to deploy something more serious, I might consider using a managed [Kubernetes](https://kubernetes.io/) environment (Digital Ocean and Linode both offer managed Kubernetes).
However, there is a non-zero chance that that'd feel too fussy / complex / unreliable and I'd revert back to a traditional server deployment.
Or perhaps I'd be surprised and find that using Fly / Kubernetes in production is a dream and I prefer it to traditional servers.
Time will tell. 😂

I even wonder if there is a future where I transition out of the deb approach and align on containers for _all_ tech stacks.
Then, even for single-binary VPS deployments, I'd install docker and run the app as a single container without any resource restrictions.
If / when the time comes to "scale up", I've already got the container and I just need to migrate from a single server to a managed Kubernetes environment.
This idea is attractive to me because it'd work for any stack.
I'm know that many developers already do this.

If the app requires a database, I'm torn between self-hosting (on another VPS) or rolling with something managed (Digital Ocean and Linode both offer managed databases).
Self-hosting Postgres is cheaper than managed and backups are pretty easy, but the high-availability machinery isn't something I've ever configured and tested myself.
Given that, I'd probably use a managed database if I had any sort of income or actual users.
Maybe even one of the fancier new ones like Neon as I mentioned earlier.
I'd have to do some pricing estimation, of course.

# TL:DR

There are so many options and I'm still undecided on a "best stack".
I do know that I'll always avoid "one-way door" stacks such as fully-serverless AWS or CloudFlare.
I don't want to depend on anything that isn't ubiquitous across hosting providers.
If (when) my opinions change as time goes on and I get more experience, I'll be sure to post updates!
