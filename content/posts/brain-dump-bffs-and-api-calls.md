---
date: 2024-08-18
title: "Brain Dump: BFFs and N+1 API Calls"
slug: "brain-dump-bffs-and-api-calls"
draft: true
---

Options for handling followed blogs:

1. N+1 on the frontend: get 1 blog, call /following N times (current)
   This feels the most "pure" but has perf implications. It suffers from the
   N+1 API call pattern and causes 21 reqs every page load. Seems wasteful
   and the page has a noticeable pause. This is how GH structures their API
   (stars) but they also use a BFF to render (SSR) lists of repos.
2. Write a batch /blogs/following endpoint that supports ID filtering
   This means the FE needs to make 2 calls (one for the blogs, one for the
   following status) and merge the data before rendering.
3. Build a BFF-ish endpoint that includes isFollowing per blog
   Would this mean introducing a new blog type like BlogWithAccount?
4. Always return isFollowing on blogs
   4a. Account is required
   4b. Account is optional, default to false if nil
   I don't love this because blogs are their own thing. Whether or not they
   are being followed by the auth'd account is separate (from a data model POV).
   But from a user point of view, blogs _do_ always have this field. I feel like
   this sacrifices API / data model purity a bit.
5. Use RR's defer feature to render blogs after one round-trip.
   Then the subsequent "is followed" data will load in shortly after. This means
   that the user doesn't have to wait for 21 requests to finish before seeing
   anything... they only have to wait for one. This also means I get to keep the
   backend API "pure". User experience is greatly improved with this approach.

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

# References

https://docs.github.com/en/rest/activity/starring?apiVersion=2022-11-28#check-if-a-repository-is-starred-by-the-authenticated-user
https://softwareengineering.stackexchange.com/questions/448013/is-it-okay-to-combine-bff-and-rest-api
https://www.infoq.com/articles/N-Plus-1/
https://samnewman.io/patterns/architectural/bff/
https://reactrouter.com/en/main/guides/deferred
