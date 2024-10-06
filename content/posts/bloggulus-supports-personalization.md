---
date: 2024-10-06
title: "Bloggulus Supports Personalization!"
slug: "bloggulus-supports-personalization"
draft: true
---

[Bloggulus](https://bloggulus.com/) is my primary "for fun" passion project.
Any frequent readers of my blog (are you out there?) will have surely heard of it before.
In short, it is a custom-built RSS feed aggregator combined with a context index powered by PostgreSQL's [full text search](https://www.postgresql.org/docs/current/textsearch.html) feature.
I use it as a playground for exploring development strategies, software architecture, and even programming languages.

What originally started out as a simple server-side rendered Go application has since taken on many forms.
For a while it was a full-stack [SvelteKit](https://kit.svelte.dev/) app, then a REST API backend with a [Svelte](https://svelte.dev/) frontend, then with a [React](https://react.dev/) frontend, and eventually back to server-side rendered Go.
It has really come full circle!
I wouldn't trade the experience and learning for anything, though.
I feel much more well-rounded as a programmer having built the same application in so many different ways.

## Navigation

The first thing you'll notice are the new navigation links.
These change depending on whether or not you are logged in.
If you aren't, you'll see "Login" and "Register".
If you are, you'll see "Blogs" and "Logout".
I'll explain the blogs page more in a bit.
I also had to move the search bar down below the nav bar but I think it looks decent inline with the articles header.
The search only applies to articles, so it makes sense to have it live closer to the content it affects.

<div style="display:flex;justify-content:center">
	<img style="max-width:400px" src="/images/20241006/index.webp" alt="Screenshot of the Bloggulus index page">
</div>

## Authentication

Check it out, a legit login page!
I feel like a true [code monkey](https://www.youtube.com/watch?v=-CuOMpoY96Y) now.
For so long, Bloggulus has only been a collection of _my_ favorite blogs.
That changes today!
By registering and a new account and logging in, you'll be able to make some tweaks.
I still need to figure how to handle password resets, though, since the service only cares about usernames and not emails.
I suppose a reset link could still be sent to an email provided by the user if and when they forget their password.

<div style="display:flex;justify-content:center">
	<img style="max-width:400px" src="/images/20241006/login.webp" alt="Screenshot of the Bloggulus login page">
</div>

## Personalization

If you have an account and are logged in, you'll gain access to the blogs page.
From here, you'll be able to add, follow, and unfollow the blogs that Bloggulus knows about.
This includes the blogs you added as well as those added by other users (myself, probably).
One neat thing about this page is that it is powered by [HTMX](https://htmx.org/)!
All three blog operations (add, follow, and unfollow) get submitted without requiring a full page reload.
Instead, only the affected HTML sections get swapped out.

<div style="display:flex;justify-content:center">
	<img style="max-width:400px" src="/images/20241006/blogs.webp" alt="Screenshot of the Bloggulus blogs page">
</div>

## What's Next?

I've got a least one more major feature planned for Bloggulus before moving onto other projects.
I'm sure I'll revisit it again eventually, but I gotta wrap up my current pending changes in preparation for [The Year of Clojure](/posts/2025-the-year-of-clojure/).
Anyhow, thanks for reading this post and if you are in the market for a minimalist RSS feed aggregator, give Bloggulus a try.
If you do, please give me as much honest feedback as possible!
