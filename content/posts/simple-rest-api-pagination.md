---
date: 2024-04-14
title: "Simple REST API Pagination"
slug: "simple-rest-api-pagination"
tags: ["Go"]
draft: true
---

Recently, I've been working on revamping my [Bloggulus](https://github.com/theandrew168/bloggulus) project to a split REST API + SPA architecture (for learning and experience).
As a part of this effort, I took a moment to research and revisit how the API handles pagination.
Since the system holds hundreds of blogs and thousands of posts, returning _all_ items from a request would be slow and unwieldy.
Instead, the user (or web frontend) should be allowed to specify which set of items they want to view.
This is where pagination comes into play.

Before I go any further, I want to give a special shoutout to Tania Rascia's [blog](https://www.taniarascia.com/) and her [post](https://www.taniarascia.com/rest-api-sorting-filtering-pagination/) about REST API design.
She has one of the best blogs on the internet and I always find myself coming back to it.
Thanks, Tania!

# Pagination

Pagination is the process of splitting up a large set of results into smaller, individual "pages".
Let's use Bloggulus' concept of "posts" as an example.
By default, the system will return the most recent 20 posts if you call `GET /api/v1/posts`.
What if, instead, you want the _second_ page of results?
Or maybe you only want to see 15 posts instead of 20.
How can we tell the server to only give us the posts that we want?
There are two common approaches but one thing is common to both: **results must have a consistent and deterministic ordering**.

# Limit and Offset

The first approach is very no-frills and "computer speak".
Instead of even bothering with the word "page", just think like a database and use the terms "limit" and "offset".
With this tactic, "limit" refers to how items you want to fetch and "offset" refers to how many items you want to skip before fetching the "limit" number of items.
Using the two examples from before, this is how you'd achieve them:

| Second Page of 20                      | First Page of 15                      |
| -------------------------------------- | ------------------------------------- |
| `GET /api/v2/posts?limit=20&offset=20` | `GET /api/v2/posts?limit=15&offset=0` |

Note that in this scenario, most APIs implement a default for both values.
In our current example, the defaults "limit" and "offset" are 20 and 0 (respectively) meaning that we can simplify our requests a bit:

| Second Page of 20             | First Page of 15             |
| ----------------------------- | ---------------------------- |
| `GET /api/v2/posts?offset=20` | `GET /api/v2/posts?limit=15` |

Overall, this approach works just fine and can handle most pagination needs.
The main downside is user experience.
For example, if you want to iterate through all pages of results, you'll have to compute each page's "offset" yourself: `offset = (page - 1) * limit`.
Additionally, "limit" and "offset" aren't the most natural terms to describe this concept and it might be more straightforward to think in terms of actual pages.

# Page and Size

Instead of thinking about this problem like a database, let's think about it like actual pages (it _is_ called pagination, after all).
We'll swap out "limilt" and "offset" for two new values: "page" and "size".
Here, "page" refers to the page of results you want to view and defaults to 1 in most REST APIs.
The other value, "size", refers to the number of results per page and defaults to 20 in this example.

Let's see how the requests look now:

| Second Page of 20          | First Page of 15            |
| -------------------------- | --------------------------- |
| `GET /api/v2/posts?page=2` | `GET /api/v2/posts?size=15` |

I find this to be bit more readable!
In the common case where users want the next page of results, they just increment "page".
If they _do_ want to utilize a different page size, they can simply change the "size".
Even with different page sizes, the "page" value will continue to work in the same way and enable easy page indexing.
The user won't have to do any special math in order to target the correct range of values.

That being said, _someone_ has to do this math because, on the backend, the database still works in terms of "limit" and "offset".
We've seen this equation before, but I'll write it again here as a refresher:

```
limit = size
offset = (page - 1) * limit
```

My backend is written in Go and the math transfer to code quite literally.
Here's the helper I wrote to handle the conversion:

```go
// Convert user-friendly "page/size" pagination to DB-friendly "limit/offset".
func PageSizeToLimitOffset(page, size int) (int, int) {
	limit := size
	offset := (page - 1) * limit
	return limit, offset
}
```

# Conclusion

This post outlined the need for pagination in REST APIs and explored two techniques for describing specific pages.
While "limit" and "offset" map more cleanly to underlying databases, "page" and "size" are a bit more user-friendly and defer a bit of the math to the server.
In the wild, you'll probably see more terms than just the four mentioned here.
Some APIs use more descriptive terms like "results_per_page" or "page_size" but the concepts are generally the same.
Thanks for reading!
