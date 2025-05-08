---
date: 2024-05-26
title: "The Court of Public Opinion"
slug: "the-court-of-public-opinion"
---

Over the past few weeks, I posted a couple of my Go-related blog posts to the [r/golang](https://old.reddit.com/r/golang/) subreddit.
I don't usually do this but have been asked the question numerous times: why don't you share your writing anywhere?
Honestly, I don't know the true answer.
Anxiety?
Fear of criticism?
Worried about wasting people's time?
Regardless of the reason, I did end up posting a few and I think it went okay!

## Solid Start

I first shared my recent post about [Parsing Recursive Polymorhpic JSON](https://old.reddit.com/r/golang/comments/1cw3egv/parsing_recursive_polymorphic_json_in_go/).
This was one received super well!
It was fun post to write and I think it has a nice, practical flow.
Some folks found it interesting and one user even said:

> Thank you for writing! This helped my problem :)

Now _that_ feels great to hear.
If my writing can help even one other development make progress, then I'm happy.
I mainly write these posts to solidify my own learning, but adding value to others is a nice bonus.

## Mixed Middle

The second blog I posted was an older one called [Conditional Embedding in Go](https://old.reddit.com/r/golang/comments/1czkqa9/conditional_embedding_in_go/).
I feel like the reception here was a bit more mixed.
While some folks "nodded" and mentioned how they follow a similar approach, other readers were left somewhat confused.
To summarize their question: why not just use the built-in `vite dev` server to host the frontend and let Go handle only the backend (while developing locally)?
This is a fair question and it made me pause and ask myself: wait, why _didn't_ I just do that?

I think that, at the time, my goal was to unite the code paths between local and production (and avoid [proxying](https://vitejs.dev/config/server-options.html#server-proxy), probably).
Instead of using two different servers for local development and embedding the frontend's static files only in prod, why not just _always_ use the latter approach (with a small twist for picking up local changes).
Did I get bit tunnel-visioned on that specific facet of the problem?
Maybe, but I still prefer my "conditional embedding" strategy at the end of the day.

## Foggy Finish

When readers are critical, my gut reaction is to be defensive or delete the post altogether.
However, after some consideration, I realize that the reality isn't so bleak.
After all, 2024 is the year of balance.
What I realized what twofold: perhaps they don't fully understand the exact nature of the problem I was solving and perhaps I don't understand their experience and knowledge.
We should probably all listen, ask questions, and eventually meet in the middle.

## P.S.

Also, my wife and I just bought a house!
Since then, I've been very busy fixing things, cleaning, and painting.
It has been a great mix of effort and reward.
Who knew houses had so many walls, doors, and sections of trim??
