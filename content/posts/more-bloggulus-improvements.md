---
date: 2024-08-11
title: "More Bloggulus RSS Improvements"
slug: "more-bloggulus-rss-improvements"
draft: true
---

Thanks to Rachel by the Bay, I’ve been able to understand the relationship between RSS feeds and readers.
More recently, I made two fixes: one was to properly detect, store, and provide values related to the “Last-Modified” header and the other was to include a customized “User-Agent” alongside any requests made by the application.
This way, Bloggulus' actions on the internet are obvious and intentional: not hidden behind the generic go/http-client UA.
