---
date: 2024-08-25
title: "Has Science Gone Too Far?"
slug: "has-science-gone-too-far"
draft: true
---

Talk about how the recent frustration when trying to reconcile “traditional API design” with BFFs and N+1 API calls. At one point I even thought: maybe GraphQL is the answer? I quickly paused, though, and wondered if this was all getting to be too much.

Bloggulus doesn’t need to be an SPA: I just wanted to explore that architecture and see how it felt to developer the frontend and backend separately. While fun and quite informative, I’m starting to feel like the stack is overkill for this project (and probably many others, to be honest).

I’m considering going to back to server-side rendered HTML for the primary web UI. I’d probably keep the “traditional API” around but might drop the BFF endpoints I’ve created already (/api/v1/articles). What about making a browser extension? That’d mainly be for writing so the traditional API should cover that. What about a mobile app? I’m not sure. Maybe I just skip it since the site is already mobile-friendly from a CSS point of view. Or perhaps there is a quick way to bundle up an SPA as a web-view and still publish it as a standalone app.
