---
date: 2024-07-14
title: "Is SvelteKit a Poor Choice for SPAs?"
slug: "is-sveltekit-a-poor-choice-for-spas"
tags: ["Go", "SvelteKit", "React"]
draft: true
---

Recently, I've been working on the frontend for my [Bloggulus](https://bloggulus.com/) project.
Back in the day, the frontend was implemented via server-rendered HTML pages.
This worked reasonably well, but I always felt like Go's HTML templates were a bit painful to work with.
I also wanted to transition the application's backend to a REST API so that I could use it as a foundation for learning how to build native apps.

# Background

Ultimately, I chose SvelteKit because it has seen a lot of hype lately and I had just finished using it to completely rewrite bloggulus as a "learn SvelteKit" project.
SvelteKit is a full-stack framework built to handle both the frontend _and backend_ of your application.
To be honest, SvelteKit is awesome when using it this way.
I was able to re-implement the entirety of Bloggulus' feature set (and more) in just a couple weeks.
Let me stress: using SvelteKit for its intended use-case of being a fullstack framework is very productive.

Anyhow, once the backend was implemented as an API, I had to choose a frontend stack.
I was considering options (for a stack plus frontend routing capabilities):

1. React + React Router
2. Vue + Vue Router
3. SvelteKit
   1. SvelteKit has a router builtin

However, that project helped me realize that I prefer to write my backends in Go.
It makes better use of multi-core systems (out of the box) and is much easier to deploy.
I deploy my apps onto ubuntu servers so "ease of deployment" is important to me.
Go's ability to bundle the entire application (frontend and backend) into a single, static binary greatly simplifies the deployment process.

# Pain Points

Sometimes, ot feels like fitting a square peg into a round hole.
SvelteKit is built for fullstack and that is the project's focus (nothing wrong with that).
Using it for client-side SPAs is technically supported, but it feels secondary / like an afterthought.
Like "Yeah, we can technically do that if you disable some stuff".
But it makes me feel like my use case is orthogonal to SvelteKit's goals and ambitions.
I felt like my needs were going "against the grain" of the framework.
It wasn't impossible to do the things I was doing, but it also wasn't the streamlined, most common use case, either.
Honestly, maybe I'd have a better time if instead of using SvelteKit, I was just using vanilla Svelte with a third-party router.

Even common SPA tasks are not discussed in the official documentation so you end up trying to solve the puzzle by reading posts on [Reddit](https://www.reddit.com/r/sveltejs/comments/z6x5uj/sveltekit_spa_with_client_side_jwt_auth/) from others who found themselves struggling with the same thing.

1. Almost all of SvelteKit's first-party documentation applies to the fullstack approach.
2. There is [one page](https://kit.svelte.dev/docs/single-page-apps) on how to use it as an SPA... by disabling some features.
3. Even on that one page, the approach is not recommended.
4. Trying to figure how to do other SvelteKit things is SPA mode is difficult: errors, auth, etc.

# Moving Forward

I think I'm going to respect the lessons I've learning from trying to use SvelteKit in a SPA setting and transition to something that is a better fit.
At the moment, I'll probably swap to React + React Router since I'm quite familiar with both and that stack also moves me closer to being able to write a native app with React Native.
I'll say that it was immediately refreshing to see features like client-side "redirect if not authed" [documented](https://reactrouter.com/en/main/start/overview#redirects) (with [examples](https://github.com/remix-run/react-router/tree/dev/examples/auth)) within React Router's first-party docs.

We'll see how this goes and I'll report back.
Maybe the grass isn't always greener and I'll find that despite making _some_ things easier, switching to React + React Router introduces more, unexpected problems that make the developer experience even worse (I hope not, thought).
I also want to reiterate that I think SvelteKit is an awesome project and framework... if you are using it in a fullstack setting.
Otherwise, it might be worth choosing technology that more closely aligns with the use cases of your project.

Thanks for reading!
