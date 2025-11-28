---
date: 2025-11-16
title: "Automating Deployments with mfd"
slug: "automating-deployments-with-mfd"
tags: ["Hosting"]
---

In a [recent post](/posts/deploying-multi-file-web-applications/), I wrote about an idea for how to deploy multi-file web applications.
Given my background in Go, I was spoiled by the "trivially bundle everything into a single binary" deployment strategy.
Life was simple: build the binary, copy it out to the server, and restart the service.
Other language ecosystems don't have it quite as easy, however.

In [Python](https://www.python.org/) or [Node.js](https://nodejs.org), for example, you can't really "build" the app into a single binary.
Tools like [shiv](https://shiv.readthedocs.io/en/latest/) exist, however, than can emulate this deployment simplicity by zipping your app's files together and then extracting them at first run.
This strategy glosses over another problem, though: what do you do about "native" dependencies?
Some Python / Node.js packages require native code to be built as part of the install.
Given that I develop on an ARM64 Mac and my servers are x86-64 Linux, I would need a way to both cross-compile native dependencies AND bundle everything together.

At some point I started to wonder if this Go-style "push" model of deployment wasn't the best fit for other tech stacks.
Instead, what about exploring the "pull" model?
Could that work?
If you read my [other post](/posts/deploying-multi-file-web-applications/) on this topic, you'd know that the answer is... yes!

# The Approach

To summarize the approach, each version (commit hash) of the app will live in a separate directory.
Then, a symlink will point to whichever version of the app is currently active.
Lastly, the systemd unit's `WorkingDirectory` will be the "active" symlink.
This ensures that each version's files don't conflict with each other or cause any deployment race conditions.

When deploying a new version, the steps are as follows:
1. Clone and checkout the code into a new directory that correponds to the version's commit hash.
2. From this new version directory, run the install and build steps.
3. Update the "active" symlink to point to the new version directly: `ln -sfn <deployment> active`.
4. Restart the systemd service to start using the code and resources from the new version.

Seems straightfoward, right?
While it _is_ a simple approach, I figured it'd still be nice to have a tool to help automate and orchestrate these steps.

# The Interface

This brings up to the topic of this post: [mfd](https://github.com/theandrew168/mfd).
At it's core, mfd is a command-line tool (written in Go) that bundles these steps into a convenient interface.
By providing the repo URL, build command(s), and an optional systemd unit, mfd will handle the rest: resolving commits, cloning and building new versions, updating the "active" symlink, and restarting the service.

The config file is written in TOML.
Here is an example that builds my [Bloggulus](https://bloggulus.com/) RSS feed aggregator web app ([source](https://github.com/theandrew168/bloggulus-svelte)).
```toml
[repo]
url = "https://github.com/theandrew168/bloggulus-svelte"

[build]
commands = [
	["npm", "install"],
	["npm", "run", "build"],
]

[systemd]
unit = "bloggulus"
```

For private repos, you can add a `repo.token` field to the config file which contains a GitHub personal access token.
For basic auth, you can instead provide `repo.username` and `repo.password`.

The mfd CLI exposes a minimal set of necessary commands:
```console
usage: mfd <command> [<args>]
commands:
  list        List available deployments
  deploy      Resolve, fetch, build, and activate a revision
  rollback    Rollback to the previous deployment
  help        Show this help message
```

You can use the `list` command to see which deployments are available and which one is active:
```
$ mfd ls
0cbbb84409331740a8a727ab09f1df163cfa4bc0
2c3e9dafda1b13f32604df11d6996fc384c9d681
92965fdb8a7c8da1fd1532ab9e3497990eee0519 (active)
```

The `deploy` command is the most useful: it performs all of the necessary deployment steps.
Without any extra args, it'll resolve the latest commit on the repo's default branch.
You can also provide tags or specific hashes (full or partial).
After a successful deployment, mfd will automatically delete all but most recent few versions.
This way, you don't have to worry about old deployments slowly using up more and more drive space.

Lastly, `rollback` can be used to revert to the deployment that came before (chronologically) the active one.
This is a handy escape hatch that can be utilized if something is amiss with the latest deployment.
By default, mfd keeps the three most recent deployments around just in case a rollback is necessary.

# Continuous Deployment

Today, my usage of mfd is completely manual.
When I want to deploy, I have to SSH into the server, change to the `/usr/local/src/<project>` directory, and run `mfd deploy`.
On the upside, this process takes less than 30 seconds and, because it happens in the foreground, I'm immediately made aware if anything goes wrong.
[Some folks](https://nickherrig.com/) have asked me: what about automated continuous deployment?
It's a reasonable question to ask and there are at least a couple of approaches I could take.

## Polling

A polling-based solution is very straightforward: run `mfd deploy` every N minutes / hours to check for and deploy new versions.
This could be done as a [cron job](https://en.wikipedia.org/wiki/Cron) or I could build a basic timer into the tool itself.
This is somewhat wasteful, though, since **most of the time I'd be cloning a repo just to do nothing**.
Despite being trivial to implement and having the least moving parts, I actually think this is the worst option compared to deploying manually or the next strategy: notifying.

## Notifying
There are two parts to this problem: triggering the deployment and tracking its status.
For [GitHub](https://github.com/) specifically, webhooks could be used to notify mfd that new version of the code is available.
This _would_ require that an HTTP endpoint be exposed to the internet so that GitHub can successfully deliver events.
Since GitHub highly recommends SSL verification for webhooks, I'd either need to utilize the [autocert](https://pkg.go.dev/golang.org/x/crypto/acme/autocert) package or instead assume that mfd will be deployed behind a reverse proxy (like [Caddy](https://caddyserver.com/)).
Furthermore, mfd would itself would need a systemd unit with a dedicated user.
This is already a non-trival amount of moving parts and architectural overhead.

Tracking is more difficult with this approach, as well, because mfd lacks a persistent user interface or way to asynchrously communicate failures.
When using the CLI, logs get printed directly to stdout and the user is expected to watch that output for any errors (this is an argument in favor of manual deployments).
When _not_ using the CLI, some other method of communication will need to be used because silent failures aren't really something that I want to support.

I think [commit statuses](https://docs.github.com/en/rest/commits/statuses?apiVersion=2022-11-28) could be a good solution here.
When a deployment starts, I could mark the commit as "pending" using the existing auth token.
Once finished, it could then be marked as "success".
If something goes wrong, it'd be marked as "error" with a useful, human-friendly message.
However, applying these updates requires a GitHub API token with permission to read and write commit statuses on the target repository.
This is another moving part and another security vector to worry about.

# The End

At the end of the day, both of the aforementioned automation strategies have tradeoffs.
Most importantly, how do I maintain visibility into whether or not the deployment succeeded?
If it failed, what went wrong?

For this reason (and the fact that I don't run many services), I'm content with simply running `mfd deploy` manually.
I can always add a `Makefile` entry that handles the SSH + `mfd deploy` command.
That way I only have to type `make deploy` and then enter my sudo password.
This brings the entire deployment process from ~30s down to ~10s which isn't much of a blocker for me.

Let me know if [mfd](https://github.com/theandrew168/mfd) could be useful to you.
I'm definitely open to any feedback and / or suggestions.
Thanks for reading!