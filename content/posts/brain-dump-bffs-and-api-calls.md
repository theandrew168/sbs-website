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
I'm beginning to realize a truth: **they are different**.
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

# N+1 API Calls

If the goal is to maintain "REST API purity", then the frontend has no choice but to make multiple calls.
Since the underlying resources are separate and normalized, the frontend needs to fetch and aggregate the disparate pieces that comprise a single page (or element).
For a list of blogs and whether or not the user follows them, this means making 1 call for the blogs and then N more calls for each one: checking if it is already followed.
This pattern is known as the "N+1 Problem" and is [well documented](https://www.infoq.com/articles/N-Plus-1/).

### Optimization: Two Trips

Note that this doesn't necessarily imply N+1 round-trips to the backend.
With proper concurrency (like [Promise.all](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/all)), the problem reduces to only two round-trips: one for the blogs and another for their followed status (all at the same time).
In pseudocode, this looks something like:

```ts
// Make 1 round-trip fetch the blogs.
const blogs = await fetch("/api/v1/blogs");

// Make N round-trips simultaneously for their followed status.
const following = await Promise.all(
  blogs.map((blog) => fetch(`/api/v1/blogs/${blog.id}/following`))
);
```

So, the browser makes N+1 API calls to render this page but it only feels like two to the user.
Two is definitely better than twenty, but it still isn't great user experience.
Also, the backend API must be able to handle high volumes of requests at the same time without much throttling / queueing.
Otherwise, the user will end up waiting much longer than expected.
That beign said, can we do better?
Can we find a way to make user only have to wait for one round-trip instead of two (or more)?

### Optimization: One Trip

What if we render the list of blogs immediately after the first round-trip and then display the follow button as that second round of requests returns to the browser?
While waiting for that second batch, we could show some sort of loading indicator (like a spinner).
Thanks for the awesome features provided by [React Router](https://reactrouter.com/en/main), this is easily doable!
The project refers to this behavior as [deferred data](https://reactrouter.com/en/main/guides/deferred).
By taking the pseudocode above, wrapping it with [defer](https://reactrouter.com/en/main/utils/defer), and rendering the follow/unfollow button inside of an [Await](https://reactrouter.com/en/main/components/await) block, we can achieve this desired behavior.

Now, the user experience is even better.
Instead of having to wait for two round-trips before seeing any content, the user only has to wait for one.
Sure, they can't follow any new blogs until all requests have finished, but I think the quicker, intermediate rendering is worth it.
At least they can see _something_ and verify that the app is working in the meantime.

### Recap

To summarize, it _is_ possible to structure the frontend data loading such that the user-facing performance converges on a single round-trip.
It also means that the backend REST API can stay "pure", normalized, and designed around individual resources.
This approach does have risks, though.
Despite the appearance of quick rendering, we _are_ making a bunch of API requests to the backend.
Any one of those requests being slow could impact the responsiveness of the page.
Do we really to shift all of this data loading and aggregation responsibility to the client?
What other options are there?

# BFF Endpoints

What if we took the individualized needs of a specific frontend (web, mobile, etc) and built a backend that was tailored to those needs?
Instead of only exposing the normalized, underlying data models of the application, what we exposed the shapes of data required to power each page?
This way, the task of collecting and aggregating data gets shifted to the backend (instead of making the frontend worry about it).
In a single API calls, the frontend could get everything needed to render a page (or at least a section of a page).
This backend could return JSON or even pre-rendered HTML.

This concept already exists and is known as the "Backend for Frontend Pattern" (often shortened to "BFF").
Similar to the "N+1 Problem", this pattern is [well documented](https://samnewman.io/patterns/architectural/bff/).

Build a BFF-ish endpoint that includes isFollowing per blog.
Would this mean introducing a new blog type like BlogWithAccount?

1. Always return isFollowing on blogs
   1. Account is required
   2. Account is optional, default to false if nil

I don't love this because blogs are their own thing.
Whether or not they are being followed by the auth'd account is separate (from a data model POV).
But from a user point of view, blogs _do_ always have this field.
I feel like this sacrifices API / data model purity a bit.

### Articles

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
