---
date: 2025-05-10
title: "Deploying Multi-File Web Applications"
slug: "deploying-multi-file-web-applications"
tags: ["Hosting"]
---

For quite a long time, I've had a question: what is the best way to deploy multi-file web applications?
When I say "multi-file", I mean any application that doesn't easily build into a single binary.
In general, web apps have many files to worry about: HTML templates, static assets, migrations, etc.
For ecosystems where the build output contains multiple directories and files, how we can safely deploy a new version of the app without any disruptions?

The solution is actually quite simple: keep each version of the app in separate directories and mark one as "active" with a [symlink](https://en.wikipedia.org/wiki/Symbolic_link)!
Let's dig into the details, talk about some background context, and see how this works in practice.

## Background

In [Go](https://go.dev), you can use the [embed](https://pkg.go.dev/embed) package to bundle your code and all other app resources into single file.
Then, deployment becomes trivial: overwrite the old binary with the new one and restart the service.
This works (and is safe) because the binary itself gets loaded into memory when the process starts.
Therefore, when the file gets replaced by the new version, the current app process is unaffected.
Only when the service restarts will the new file be loaded and executed.

As it turns out, most other ecosystems are _not_ like Go and lack this functionality.
I'm talking about languages like Python (Flask, Django, etc) or Node.js (Next.js, Remix, SvelteKit, etc).
For these apps where multiple files _are_ required to be deployed, it isn't completely safe to just drop the new directory structure over top of the old one and restart.
While you _can_ do that and it'll technically work, there is a potential race condition issue.

Requests that reference old resources could break in the short gap between replacing the underlying app files and restarting the service.
For example, consider a request for the old version of the app that is bound for `/img/foo.png`.
If the new version removes it, the request could arrive _between_ deploying the new files and restarting the service which would result in an unexpected error (because the image can no longer be found).
How can this be solved?

## Deployment

For these multi-file applications, what actually needs to happen during a deployment?
I'm going to use the example of a Node.js app managed via systemd on an Ubuntu 24.04 server.
First things first: the server will need to have the app's runtime installed.
This is a one-time operation that can be performed when initializing the server: `apt install nodejs`.

Next, we need to get the source code to the server somehow.
We could either "push" the code to the server by copying and unzipping an archive or "pull" the code by cloning or downloading it directly.
Each method has tradeoffs and should be handled on a case-by-case basis.
In my experience, though, "pulling" (via `git clone` and `git checkout <hash>`) is pretty convenient, especially when dealing with open-source application that don't require any authentication to access the source code.

Then, for any given version of the app, we need to do some additional prep work before it will be ready to run: installing dependencies and building the code.
Installing dependencies is easy enough: `npm install`.
Building is usually quite simple as well (depending on the framework): `npm run build`.

The last question then becomes: how do we run this?
What command do we put in the `ExecStart` field of our [systemd service](https://www.freedesktop.org/software/systemd/man/latest/systemd.service.html#ExecStart=) file?
Well, it depends on the framework.
For something like Next.js, we'd run `npm run start`.
For a library that outputs a Node.js-ready "build" directory (like SvelteKit), the command would be `node build/`.

With all that out of the way, we can now restart the app: `systemctl restart <myapp>`.
And there we have it: a new version of the app has been deployed!

## Idea

Naively, this approach has issues.
As mentioned before, performing the deployment steps **in the same directory as the currently-running app** can lead to race conditions and ultimately user-facing errors.
Once the app is running, how can a new version of the code be deployed, built, and started **without impacting the app's availability**?

In other non-single-binary ecosystems such as Python, tools like [shiv](https://shiv.readthedocs.io/en/latest/) solve this problem by zipping your project into a single archive and then transparently unzipping it into a unique directory prior to execution.
This way, the individual files of different versions don't conflict with each other and cutover race conditions are avoided.

Can we do something similar to this but in a more ubiquitous way?
Unfortunately, systemd doesn't have a native way to say "use this parameterized working directory" at startup.
Instead, one possible solution to this that I really like is simply making the `WorkingDirectory` a symlink ([reference](https://unix.stackexchange.com/questions/242019/set-workingdirectory-using-a-variable/629958#629958)).
Then, that symlink can point to any given version of the project available on the server.

## Solution

How would deployments work using this approach?
We'd first have a separate directory for each version of the app that has been deployed (`/usr/local/src/<myapp>/<hash>` for example).
There would then exist a symlink (something like `/usr/local/src/<myapp>/active`) that points to whichever version of the app is currently active.
Lastly, in the systemd service file, we would set `WorkingDirectory` to our "active" symlink.

When deploying a new version, the steps are as follows:
1. Clone and checkout the code into a new directory that correponds to the version's commit hash.
   1. You could also use a tag or something else here, if you prefer.
2. From this new version directory, run the install and build steps.
3. Update the "active" symlink to point to the new version directly: `ln -sfn <hash> current`.
4. Restart the systemd service to start using the code and resources from the new version.

And... that's it!
That's the whole approach: build each version of the app in a separate directory and use a symlink to control which version is "active".
Also, because systemd checks `WorkingDirectory` at startup, the old files will still be used even after the symlink is swapped.
The new files won’t be picked up until the service restarts which avoids the race condition from earlier.

## Conlusion

In this post, we went over an approach for safely deploying multi-file web applications.
Since most web app frameworks and ecosystems _don't_ support bundling everything into a single, convenient binary, we have to get a bit creative.
Thankfully, with a bit of clever symlinking, we can ensure that new versions are deployed without impacting the service's availability.
Just be sure to clean up the old version directories from time to time and you'll be off to the races!

Thanks for reading!