---
date: 2024-06-02
title: "Bloggulus: A Responsible RSS Reader"
slug: "bloggulus-a-responsible-rss-reader"
draft: true
---

Recently, Rachel by the Bay wrote about an idea for "scoring" RSS feed readers.
Criteria would be things like: how often do they check, do they respect caching headers, etc.
Since I maintain my own RSS reader / aggregator called Bloggulus, I was curious how correct / responsible my own program was.
I messaged Rachel with my interest in participating and she promptly replied with instructions to get it setup.

Fast-forward a few days and Bloggulus is scoring quite well!
It has only sent one unconditional request (which is expected when a new feed is added).
All `If-None-Match` headers values have been valid.
No useless cookies, referrers, or query params have been sent, either.
Not bad!

Sometimes, however, the system ends up polling slightly sooner than hourly.
Occasionally it'll re-poll after 59 minutes which isn't huge deal.
Once, though, it re-requested content after only 44 minutes.
I'm not 100% sure _why_ this happened but it could be that the server restarted.
When that happens, the "wait an hour" timers get restarted.
I should probably update the code to account for that!
I'd have to track (in the database) _when_ I most recently synced each blog but that isn't a huge deal.

Thanks for reading and keep your RSS reading responsible!
