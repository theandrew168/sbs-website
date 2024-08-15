---
date: 2024-08-11
title: "More Bloggulus RSS Improvements"
slug: "more-bloggulus-rss-improvements"
---

Thanks to [Rachel by the Bay](https://rachelbythebay.com/w/) and her posts about RSS correctness / etiquette, I’ve been able to more deeply understand the relationship between feeds and readers.
A few months ago, I decided to use [Bloggulus](https://bloggulus.com) to participate in her feed reader scoring program.
I [wrote about it](/posts/bloggulus-a-responsible-rss-reader/) at the time and noted how Bloggulus' performance as a responsible RSS reader was "good but not great".

Recently, while on vacation, I decided to fix the remaining three issues with Bloggulus:

1. Occasionally polling more than once per hour ([commit](https://github.com/theandrew168/bloggulus/commit/2e8003a32b78ea83c5dcb69ed17e4daa7c3e4ccb))
   1. This was simply caused by narrow time intervals
2. Not correctly detecting and storing the `Last-Modified` header ([commit](https://github.com/theandrew168/bloggulus/commit/78ae3b1dda24f720fec2d24905734298ef20275a))
   1. This was caused by a simple typo
3. Using Go's default user agent when making HTTP requests ([commit](https://github.com/theandrew168/bloggulus/commit/15042284ccffeb714655f40456c0ad8b499fc89b))
   1. Requests now set their user agent to: `Bloggulus/X.Y.Z (+https://bloggulus.com)`

# The Data

Looking at the data Rachel presents to those in the program, you can immediately see when the first change took effect (fixing the polling interval to always be at least one hour):

![Reading too quickly errors disappear](/images/20240811/interval.webp)

Further down, we can see the impact of the other two improvements (adding a user agent and fixing the `Last-Modified` detection):

![Better user agent and modified detection](/images/20240811/user-agent-last-modified.webp)

# Conclusion

At this point, I think Bloggulus is as well behaved as an RSS reader can be!
This is awesome and I'm grateful to Rachel for being so vocal about how these applications _should_ work.
It gives me confidence that Bloggulus is built on a solid foundation and ready to consume many more feeds in the future.

Thanks for reading!
