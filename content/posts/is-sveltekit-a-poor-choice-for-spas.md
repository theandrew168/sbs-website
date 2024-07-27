---
date: 2024-07-14
title: "Is SvelteKit a Poor Choice for SPAs?"
slug: "is-sveltekit-a-poor-choice-for-spas"
tags: ["Go", "SvelteKit", "React"]
---

Recently, I've been working on the frontend for my [Bloggulus](https://bloggulus.com/) project.
Back in the day, the frontend was implemented via server-rendered HTML pages.
This worked reasonably well, but I always felt like Go's HTML templates were a bit painful to work with.
I also wanted to transition the application's backend to a REST API so that I could use it as a foundation for learning how to build native apps (someday...).

# Background

Ultimately, I chose [SvelteKit](https://kit.svelte.dev/) because it has seen a lot of hype lately and I had just finished using it to completely rewrite Bloggulus as a "learn a new skill" project.
SvelteKit is a full-stack web development framework built to handle both the frontend _and backend_ of your application.
To be honest, SvelteKit is awesome when using it this way.
I was able to re-implement the entirety of Bloggulus' feature set (and more) in just a couple weeks.
I want to be clear: using SvelteKit for its intended use-case of being a full-stack framework is very productive.

However, that Bloggulus rewrite project also helped me realize that I prefer to write my backends in Go.
It makes better use of multi-core systems (out of the box) and is much easier to deploy.
I deploy my apps onto bare-ish metal servers so "ease of deployment" is important to me.
Go's ability to bundle the entire application (frontend and backend) into a single, static binary greatly simplifies the deployment process.
Despite being marketed as a full-stack framework, SvelteKit _can_ be used to build single-page applications (SPAs).

# Pain Points

It's sort of hard to explain, but using SvelteKit this way feels like trying to fit a square peg into a round hole.
Trying to figure out how to utilize basic features while in SPA mode is difficult: **handling errors and implementing client-side auth** are just a few examples.
It makes me feel like my use case is orthogonal to the project's goals and ambitions: like my needs were going "against the grain" of the framework.
It wasn't _impossible_ to do the things I needed to do, but it also wasn't the use case that the library was optimized / built for.

> Using SvelteKit for building SPAs is technically supported but feels like an afterthought.

Being more specific, though, I found that common SPA tasks are not discussed in the official documentation so you end up having to look elsewhere for solutions.
I often had to rely on [Reddit posts](https://www.reddit.com/r/sveltejs/comments/z6x5uj/sveltekit_spa_with_client_side_jwt_auth/) from others who found themselves struggling with the same things.
If fact, there is only [one page](https://kit.svelte.dev/docs/single-page-apps) in the docs on how to use SvelteKit as an SPA.
Perhaps I'd have a better time if, instead of using SvelteKit, I just used vanilla [Svelte](https://svelte.dev/) with a third-party router.

# Moving Forward

I think I'm going to respect the lessons I've learned here (and the frustration I've felt) and transition to something that is a better fit.
At the moment, I'm planning to swap to React + React Router since I'm quite familiar with both.
That stack also moves me closer to being able to write a native app (with [React Native](https://reactnative.dev/)).
I'll say that it was immediately refreshing to see features like client-side "redirect if not authed" [documented](https://reactrouter.com/en/main/start/overview#redirects) (with [examples](https://github.com/remix-run/react-router/tree/dev/examples/auth)) within React Router's official docs.

We'll see how this goes and I'll report back.
Maybe the grass isn't always greener and I'll find that despite making _some_ things easier, switching to React + React Router introduces more, unexpected problems that make the developer experience even worse (I hope not, though).
I also want to reiterate that I think SvelteKit is an awesome project and framework... if you are using it in a full-stack setting.
Otherwise, it might be worth choosing technology that more closely aligns with the use cases of your project.

Thanks for reading!
