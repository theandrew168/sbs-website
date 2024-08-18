---
date: 2024-08-18
title: "Brain Dump: BFFs and N+1 API Calls"
slug: "brain-dump-bffs-and-api-calls"
draft: true
---

Writing software is hard.
Sometimes, I'll find myself stuck on a problem for multiple days or even weeks.
When that happens, I find it useful to simply write out everything I know about the problem: the nuances, possible solutions, external references, etc.
My most recent head-scratcher has been about balancing "purist" REST API design with the needs of a web frontend.
I'm beginning to realize a truth: they are different.
As always, I'm talking about [Bloggulus](https://bloggulus.com).

# The Problem

I've recently been working on adding support for individualized feeds.
This means that users can create an account and follow their own favorite blogs.
Under the hood, this involves adding a new many-to-many relationship between accounts and blogs (easy enough).
The problem comes when creating a UI for users to follow and unfollow blogs.
See, from a data model point of view, following blogs is a separate concept (separate table).
The `GET /api/v1/blogs` endpoint does NOT include an `isFollowing` field.
Instead, the `GET /api/v1/blogs/{blogID}/following` endpoint can be used to check if a specific blog is being followed by the authenticated user.

Given this limitation, how do we build this page?
Pardon the ugliness...

![Bloggulus blogs page with follow and unfollow buttons](/images/20240818/blogs.webp)

# Possible Solutions

### N+1 API Calls

[InfoQ - N+1 Problem](https://www.infoq.com/articles/N-Plus-1/)

N+1 on the frontend: get 1 blog, call /following N times (current).
This feels the most "pure" but has perf implications.
It suffers from the N+1 API call pattern and causes 21 reqs every page load.
Seems wasteful and the page has a noticeable pause.
This is how GH structures their API (stars) but they also use a BFF to render (SSR) lists of repos.

[React Router - Deferred](https://reactrouter.com/en/main/guides/deferred)

Use RR's defer feature to render blogs after one round-trip.
Then the subsequent "is followed" data will load in shortly after.
This means that the user doesn't have to wait for 21 requests to finish before seeing anything... they only have to wait for one.
This also means I get to keep the backend API "pure". User experience is greatly improved with this approach.

Write a batch /blogs/following endpoint that supports ID filtering.
This means the FE needs to make 2 calls (one for the blogs, one for the following status) and merge the data before rendering.

### BFF Endpoints

[Sam Newman - BFF Pattern](https://samnewman.io/patterns/architectural/bff/)

Build a BFF-ish endpoint that includes isFollowing per blog.
Would this mean introducing a new blog type like BlogWithAccount?

1. Always return isFollowing on blogs
   1. Account is required
   2. Account is optional, default to false if nil

I don't love this because blogs are their own thing.
Whether or not they are being followed by the auth'd account is separate (from a data model POV).
But from a user point of view, blogs _do_ always have this field.
I feel like this sacrifices API / data model purity a bit.

# Other Thoughts

The whole concept of "articles" is basically BFF.
The frontend _could_ just query for recent posts and then make subsequent reqs for blog metadata
(title, url) and tags.
But that'd be kinda complex on the FE and probably slow.
Lots of waterfall there.
The BFF shifts that complexity to SQL and the rest of the app (BE and FE) are faster and simpler (one API call, one
query).
Do I really wanna get into the habit of finding synonyms for stuff just to support the frontend?
Articles and Posts are the same thing but in different contexts.

Should every page load only require one API call?
That points me toward option 3 or 4.
How many is too many?
Maybe the bigger issue is round-trips.
Like two round-trips before showing anything is poor UX.
But I could probably get away with one (like doing multiple initial loads at the same time).
The waterfall is what hurts.
Maybe the FE / react-router can help here by rendering something as soon as the blogs come in.
Then, the following status will pop in when ready.
Can RR do that sort of thing?
Show individual loaders within a list?

# An Existing BFF Example

As it turns out, I've already solved a similar problem via the "BFF Endpoint" approach (I just didn't know it had a name).
Blogs can have multiple posts (one-to-many), and posts can have multiple tags (many-to-many).
In the API (and database), these are represented as separate resources that can be CRUD'd individually.
However, on the frontend, we frequently need to bundle these disparate models together into something useful for the user.
Here is a screenshot from the application:

![Bloggulus home page with three posts](/images/20240818/bloggulus.webp)

Let's take the first "article" as an example.
You can see how data from all three underlying models are aggregated into a single unit.
The published date, post title, and post URL all come from the underlying **post** model.
The blog title and URL come from the **blog** model.
Lastly, the tag names come from the **tag** model.
Knowing all of this, the question becomes: how do we efficently fetch and render this data?
