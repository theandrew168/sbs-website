---
date: 2025-11-16
title: "Introducing mfd"
slug: "introducing-mfd"
tags: ["Hosting"]
draft: true
---

In a [recent post](/posts/deploying-multi-file-web-applications/), I wrote about an idea for how to deploy multi-file web applications.
Given my background in Go, I was spoiled by the "trivially bundle everything into a single binary" deployment strategy.
Life was simple: build the binary, copy it out to the server, and restart the service.
Other language ecosystems don't have it quite as easy, however.

In Python or Node.js, for example, you can't really "build" the app into a single binary.
Tools like [shiv](https://shiv.readthedocs.io/en/latest/) _can_ emulate this by zipping your app's files together and then extracting them at first run.
This strategy glosses over another problem, though: what do you do about "native" dependencies?
Some Python / Node packages require native code to be built as part of the install.
Given that I develop on an ARM64 Mac and my servers are x86-64 Linux, I would need a way to both cross-compile native dependencies AND bundle everything together.

At some point I started to wonder if this Go-style "push" model of deployment wasn't the best fit for other tech stacks.
Instead, what about exploring the "pull" model?
Could that work?
If you read my [other post](/posts/deploying-multi-file-web-applications/) on this topic, you'd know that the answer is... yes!

# The Approach

To summarize the approach, each version (commit hash) of the app will exist in an isolated directory.
Then, a symlink (which I've named "active") will point to whichever version of the app is currently active.
Lastly, the systemd unit's `WorkingDirectory` will be the "active" directory.
This ensures that each version's files don't conflict with each other or cause any deployment race conditions.

When deploying a new version, the steps are as follows:
1. Clone and checkout the code into a new directory that correponds to the version's commit hash.
   1. You could also use a tag or something else here, if you prefer.
2. From this new version directory, run the install and build steps.
3. Update the "active" symlink to point to the new version directly: `ln -sfn <hash> current`.
4. Restart the systemd service to start using the code and resources from the new version.

Seems simple, right?
While it _is_ a simple approach, it'd nice to have a tool to help automate and orchestrate these steps.

# The Tool

This brings up to the topic of this post: [mfd](https://github.com/theandrew168/mfd).
At it's core, mfd is a tool that bundles these steps into a simple config file + CLI.
By providing the repo URL and build command(s), mfd will handle the rest: resolving commits, cloning and building new versions, and updating the "active" symlink.
The config file is written in TOML.

Here is an example that builds my [Bloggulus](https://bloggulus.com/) RSS feed aggregate web app ([source](https://github.com/theandrew168/bloggulus-svelte)).
```toml
[repo]
url = "https://github.com/theandrew168/bloggulus-svelte"

[build]
commands = [
	["npm", "install"],
	["npm", "run", "build"],
]
```

For private repos, you can add a `repo.token` field to the config file which contains a GitHub personal access token.
For basic auth, you can instead add `repo.username` and `repo.password`.

The mfd CLI provides a few useful options:
```console
usage: mfd <command> [<args>]
commands:
  list        List available deployments
  deploy      Resolve, fetch, build, and activate a revision
  resolve     Resolve a revision to a deployment
  activate    Activate a deployment
  rollback    Rollback to the previous deployment
  remove      Remove a deployment
  clean       Remove old deployments
  help        Show this help message
```

The `deploy` command is the most useful: it performs all of the necessary deployment steps.
Without any extra args, it'll resolve the latest commit on the repo's default branch.
You can also provide tags or specific hashes (full or partial).

You can use the `list` command to see which deployments are available and which one is active:
```
$ mfd ls
0cbbb84409331740a8a727ab09f1df163cfa4bc0
2c3e9dafda1b13f32604df11d6996fc384c9d681
92965fdb8a7c8da1fd1532ab9e3497990eee0519 (active)
```

Lastly, you can use the `clean` command to remove all but the most recent few deployments.

# The Future

Today, usage of the tool is still manual.
When I want to deploy, I have to SSH into the server and run `mfd deploy`.
I've been thinking about ways to automate this while maintaining visibility: did the deployment succeed or fail?
If it failed, what went wrong?

There are two parts to this problem: triggering the deployment and tracking / display its status.
Triggering the deployment is straightforward enough: use polling or webhooks to detect that a new version is available.
Tracking is more difficult because mfd lacks a persistent user interface.
When using the CLI, logs simply get printed to stdout and the user is expected to watch that output for any errors.

I think [commit statuses](https://docs.github.com/en/rest/commits/statuses?apiVersion=2022-11-28) could be a good solution here.
When a deployment starts, I can mark the commit as "pending" using the existing auth token.
Once finished, it can then be marked as "success".
If something goes wrong, it'll be marked as "error" with a useful message.

# The End

Let me know if mfd could be useful to you.
I'm definitely open to any feedback and / or suggestions.
Thanks for reading!