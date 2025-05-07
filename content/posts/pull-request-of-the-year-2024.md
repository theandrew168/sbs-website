---
date: 2025-05-06
title: "Pull Request of the Year 2024"
slug: "pull-request-of-the-year-2024"
draft: true
---

I know, I know.
It's already 2025 and I only came across this PR a few weeks ago.
That being said, it was merged in 2024 and I can't recall seeing any other PRs from last year that had me stomping and hollering so hard.
Finding this feature (and the PR) came at a time of frustration: realizing that I wouldn't be able to use systemd socket activate for NextJS (or RemixJS, for that matter) without having to jump through hoops.
By the end, I'd given up any hope for first-party support for this and figured I'd have to run the application with a "custom server": where I use Node's http server directly to host the pre-built application.

Here's the thing about writing "systemd socket aware applications": it requires zero user configuration (for the app itself, at least) and is a no-op for non-systemd cases.
If the `LISTEN_FDS` env var is set, then you can file descriptors directly (starting at 3) instead of opening a new socket on a given addr + port.
If that var is NOT set, then just open the socket like normal.

Why are you even writing web apps in JS?
Aren't you a hardcore server-side rendered, semantic HTML w/ forms gray beard?
While that is usually true, I do really like TypeScript as a language and I had the itch to build something "trendy" using "trendy" tools.
I wanted to see what all the vibe-coding hype was about (while applying my own supervision and scrutiny, of course).
I used [v0.dev](https://v0.dev/) to build a quick prototype and since _that_ used Next.js by default, that's where I decided to start.

For [Next.js](https://nextjs.org/), you can do this, but it has [limitations](https://nextjs.org/docs/pages/guides/custom-server):
> Before deciding to use a custom server, keep in mind that it should only be used when the integrated router of Next.js can't meet your app requirements. A custom server will remove important performance optimizations, like Automatic Static Optimization.

For [Remix](https://remix.run/), you use the [createRequestHandler](https://remix.run/docs/en/main/other-api/adapter#createrequesthandler) function from the `@remix-run/express` package to create the proper handler and then host it with your own Express server.
This would probably work for my use case but I wasn't using Remix at the time.

All that said, I wouldn't need any of this "custom server" stuff if these libraries had built-in support for checking for systemd sockets and using them (instead of opening their own).
I looked into the NextJS source code to consider if I could add this behavior myself but the lead time for developing, testing, and merging such a change would surpass my project timeline.
I may still do this in the future, though, since a reference implementation already exists... in [SvelteKit](https://svelte.dev/docs/kit/introduction)!

Honestly, my expectations were so low that I didn't even start my search by checking if SvelteKit had builtin support for socket activation: I just figured I'd have to take the build output and run it with my own custom NodeJS server like all the others.
That led me to the [adapter-node](https://svelte.dev/docs/kit/adapter-node) docs.
While this adapter _does_ easily support deploying a pre-built SvelteKit app with a custom server, I didn't even need it.
Looking at the sidebar, I couldn't believe my eyes: there was a section titled "Socket activation".
It was _exactly_ what I'd been looking for.
No joke, I was slack-jawed while reading through it.

I though to myself: someone in the SvelteKit community cares enough about traditional app deployment strategies to add this to the project?
I had to see [the PR](https://github.com/sveltejs/kit/pull/11653) for myself.
Did the whole feature come in one PR or was it multiple? (Just one).
Did one person write it? Essentially ([Karim Jordan](https://github.com/karimfromjordan)).
How did the implementation work? Check the env vars and pass the `fd` to the Node http server.
Was it complex or relatively simple? Pretty simple!

This this is actually already possible with Next.js or Remix and I simply missed it in the documentation / code, please let me know!
Until then, I'll be sticking with the framework that respects my server-ful, single node, systemd socket activated, "old school" style of [deploying applications](/posts/deploying-multi-file-web-applications/).