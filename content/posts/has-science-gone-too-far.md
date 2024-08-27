---
date: 2024-08-25
title: "Has Science Gone Too Far?"
slug: "has-science-gone-too-far"
---

I recently [wrote an article](/posts/brain-dump-bffs-and-api-calls/) explaining my difficulty with trying to reconcile “traditional API design” with BFFs (backend for frontend) and N+1 API calls.
It seemed tough to arrive at a "best" solution (and maybe one doesn't even really exist).
Do I build my API for general consupmtion (granular and resource-based) or specifically for my web frontend (coarse with multiple resources joined together)?
I settled for a middleground: use intentional frontend techniques to load secondary data _after_ rendering the initial, primary data.
Then, if that was too cumbersome, create a BFF endpoint to power heavier pages.

Even after documenting everything I knew about the designs and their tradeoffs, I didn't feel like the problem was truly resolved.
What if I create a mobile app that requires _different_ BFF endpoints?
Do I just make them as needed?
I really wished that I could just create a "write one, use everywhere" REST API.
At one point I wondered: is [GraphQL](https://graphql.org/) the answer?
After the upfront cost of specifying the data model and implemeting resolvers, clients would be able to request whatever data they need: granular OR coarse.
Sounds perfect, right?

But then, I took a beat.
I paused and remembered an old proverb:

> You can rarely make a system simpler by adding a bunch of stuff.

Was this was all getting to be too much?
Had science gone too far?
Bloggulus doesn’t need to be an SPA: I just wanted to explore that architecture and see how it felt to develope the frontend and backend separately.
While fun and quite informative, I’m starting to feel like the stack is overkill for this project (and probably many others, to be honest).
Perhaps it's time to get off of this wild ride through SPA-town.
To that end, I’m considering going to back to server-side rendering for the primary web UI.

I’d probably keep the “traditional API” around but might drop the BFF endpoints I’ve created already (such as `GET /api/v1/articles`).
What about making a browser extension?
That’d mainly be for modifying app state (following a blog / adding a page) so the traditional API should cover that.
What about a mobile app?
I’m not sure.
Maybe I just skip it since the site is already mobile-friendly from a CSS point of view.
Or perhaps there is a quick way to bundle up an SPA as a web-view and still publish it as a standalone app.
It seems like there has been [some work](https://developers.google.com/codelabs/pwa-in-play) done on this concept.
