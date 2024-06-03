---
date: 2024-06-02
title: "Bloggulus: A Responsible RSS Reader"
slug: "bloggulus-a-responsible-rss-reader"
---

Recently, Rachel by the Bay [wrote about an idea](https://rachelbythebay.com/w/2024/05/29/score/) for "scoring" RSS feed readers.
Criteria would be things like: how often do they check, do they respect [caching headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-None-Match), etc.
Since I maintain my own RSS reader / aggregator called [Bloggulus](https://bloggulus.com), I was curious to see how responsible my own program was.
Her previous posts have directly influenced my ability to make Bloggulus more efficient and more correct.
I messaged Rachel with interest in participating and she promptly replied with instructions on how to set things up.

Fast-forward a few days and Bloggulus is scoring quite well!

![RSS reader score](/images/20240602/rss-score.webp)

It has only sent one unconditional request (which is expected when a new feed is added).
All `If-None-Match` header values (this is a cache-related header) have been valid.
No useless cookies, referrers, or query params have been sent, either.
Not bad!

Sometimes, however, Bloggulus ends up polling slightly more oftan than the intended hourly schedule.
Occasionally, it'll re-poll after 59 minutes which isn't huge deal (just subtle timing variance).
Once, though, it re-requested content after only 44 minutes.
I'm not 100% sure _why_ this happened but it could be that the server restarted.
When that happens, the "wait an hour before polling again" timer gets restarted (likely for unattended upgrades).
I do have logic in-place to prevent "restart spam" by ensuring that a single feed doesn't polled more than once every 30 minutes.
I might consider increasing that duration value.

Anyhow, thanks for reading and for keeping your RSS readers responsible!
