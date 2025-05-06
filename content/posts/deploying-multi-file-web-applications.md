---
date: 2025-05-06
title: "Deploying Multi-File Web Applications"
slug: "deploying-multi-file-web-applications"
draft: true
---

Since this is a NodeJS app, you can't really bundle / build it into a clean, single binary for deployment (like I do with Go).
Given that, what needs to happen?
I still want to use systemd.
The server will need to have NodeJS installed to run the app.
Prior to starting, we need to run both `npm install` and `npm run build`.
Actually running the app will run `npm run start` or perhaps simply `node build/`.

Naively, this approach has issues.
Once the app is running, how can a new version of the code be deployed, built, and started **without impacting the currently running app**?
This problem also exists for other non-single-binary ecosystems such as Python.
There, tools like [shiv](https://shiv.readthedocs.io/en/latest/) solve the problem by transparently unzipping the app into a unique directory prior to running.
That way, the individual files of different versions don't conflict with each other and cutover race conditions are avoided.

Can we do something similar?
Systemd doesn't have a native way to say "use this arbitrary working directory (controlled via a variable, maybe?)".
One possible solution to this that I really like is simply making the `WorkingDirectory` a symlink ([reference](https://unix.stackexchange.com/questions/242019/set-workingdirectory-using-a-variable/629958#629958)).
Then, that symlink can point to any given revision of the project available on the server.

How would deployments work using this approach?
We'd first have a separate directory for each revision of the app that has been deployed (`/usr/local/bin/fussy/384a283` for example).
There will then exist a symlink (something like `/usr/local/bin/fussy/active`) that points to whichever version of the app is currently active.

When deploying a new version, clone and checkout its code into a new directory that correponds to its commit hash.
Then, from this new version directory, run both of the "build" steps: `npm install` and `npm run build`.
Next, update the "active" symlink to point to the new version directly.
Lastly, restart the systemd service which will pick up the code from the new version.

This will switch to the new code while minimizing the amount of “race condition” time spent running the old process with the new files.
Does NodeJS / NextJS even read anything from the FS once the service is running?
Surely something will be read from the FS during execution, like static resources or maybe compiled JS / CSS resources?
Either way, this approach minimizes the risk (though it doesn’t completely eliminate it).

Actually, this might work completely as intended without the “old process, new files” race condition.
Since systemd checks `WorkingDirectory` at startup, the old files should still be used even after the symlink is swapped.
The new files won’t be picked up until the service restarts (which also starts a new process). This might be perfect!