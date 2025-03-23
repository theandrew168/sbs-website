---
date: 2025-03-23
title: "Embracing Functional Simplicity"
slug: "embracing-functional-simplicity"
draft: true
---

I’ve been thinking a lot about software design and how to write better code.
Better here means more correct, more readable, easier to follow, and easier to test.
Two books have been influencing me lately: [Grokking Simplicity](https://grokkingsimplicity.com/) by Eric Normand and [Architecture Patterns with Python](https://www.cosmicpython.com/) by Harry Percival and Bob Gregory.
The former focuses on functional design: emboldening pure, data-in-data-out functions.
The latter hoists the ideas up to a higher level and thinks about project / application architecture and how to write code that complements the domain of your problem.
This idea was also expressed in Brandon Rhodes' [Hoist Your I/O](https://www.youtube.com/watch?v=PBQN62oUnN8) talk.

These ideas have already led me toward better code in my Bloggulus project.
More so on the functional side, so far.
The nastiest, most critical part of the app is the “sync” process.
Isn’t that usually how it goes: the most important part of the system is the one that is tightly-coupled to the outside world and therefore difficult to test, verify, and trust?
Here are some meaningful refactorings that I’ve done.

We’re going to separate what we want to do from how to do it.
Actions vs Calculations.
Note that the fetching of RSS feeds is implemented as an interface, but testing this can _still_ be painful.

# Before

This version is kind of a mess: domain logic is mixed in with reckless abandon, actions and calculations are poorly defined yet heavily intertwined, and "what" vs "how" is confusingly mixed.

```python
# List all blogs in the database.
blogs = db.list_blogs()

# Process each blog individually.
for blog in blogs:
	# Skip the blog if it has been recently synced.
	delta = now - blog.synced_at
	if delta < sync_interval:
		continue

	# Update the blog's synced_at time.
	db.update_blog_synced_at(blog, now)

	# Fetch the blog's RSS / Atom feed.
	feed = fetch_feed(blog)

	# If changed, update the blog's cache headers.
	if feed.etag:
		db.update_blog_cache_headers(blog, feed.etag)
	if feed.last_modified:
		db.update_blog_cache_headers(blog, feed.last_modified)

	# If the feed indicates no new content, skip it.
	if not feed.has_changed:
		continue

	# Process each feed post individually.
	for feed_post in feed.posts:
		# Check if the post already exists in the database
		post = db.read_post(feed_post.url)

		# If it doesn't, create it.
		if not post:
			post = db.create_post(feed_post)

		# If the post's content has changed, update it.
		if feed_post.content and not post.content:
			db.update_post_content(post, feed_post.content)
```

# Breakdown

Grokking Simplicity focuses on this idea of bucketing your code into two distinct categories: actions and calculations.
Actions are places where your code interacts with the outside world: a database, the internet, etc.
Calculations are pure functions where data goes in and data comes out: the outside world is completely irrelevant and unaffected.
In general, calculations are "better" than actions.
They are easier to write, test, and understand.

```
// Action: List all blogs in the database
// Calculation: Determine which blogs need to be synced (FilterSyncableBlogs)
// Calculation: For each sync-able blog, create it's FetchFeedRequest (CreateSyncRequest)
// Action: Update sync time for each sync-able blog
// Action: Exchange the request to fetch the blog for a response (with limited concurrency)
// Calculation: Update the blog's cache headers if changed
// Action: If changed, save the updated blog cache headers to the database
// Calculation: If the response includes data, parse the RSS / Atom feed for posts
// Action: List all posts for the current blog
// Calculation: Determine if any posts in the feed are new / updated
// Action: Create / update posts in the database
```

# Helpers

1. [CanBeSynced](https://github.com/theandrew168/bloggulus/blob/f4be7f7edaefdfa485d7d30bccfc3af17482d4e3/backend/model/blog.go#L104)
2. [FilterSyncableBlogs](https://github.com/theandrew168/bloggulus/blob/f4be7f7edaefdfa485d7d30bccfc3af17482d4e3/backend/service/sync.go#L43)
3. [UpdateCacheHeaders](https://github.com/theandrew168/bloggulus/blob/f4be7f7edaefdfa485d7d30bccfc3af17482d4e3/backend/service/sync.go#L54)
4. [ComparePosts](https://github.com/theandrew168/bloggulus/blob/f4be7f7edaefdfa485d7d30bccfc3af17482d4e3/backend/service/sync.go#L76)

# After

This version feels much better.
Domain concepts are bit more separated: `can_be_synced`, `filter_syncable_blogs`, `update_cache_headers`.
The big win, though, is the creation of the `compare_posts` helper.
This is better in multiple ways:

1. Some of the trickiest, nastiest logic is now a pure function.
2. The "what" (create vs update) is separated from the "how" (interacting with the database).
3. Posts are fetched and processed in a batch (per blog) which is faster than fetching and processing each post individually.

```python
# List all blogs in the database.
blogs = db.list_blogs()

# Filter blogs down to those that are able to be synced.
syncable_blogs = filter_syncable_blogs(blogs)

# Update the synced_at date for all syncable blogs.
for blog in syncable_blog:
	db.update_blog_synced_at(blog, now)

# Process each syncable blog individually.
for blog in syncable_blogs:
	# Fetch the blog's RSS / Atom feed.
	feed = fetch_feed(blog)

	# Check if the blog's cache headers have changed and update.
	has_changed = update_cache_headers(blog, feed)
	if has_changed:
		db.update_blog_cache_headers(blog)

	# If the feed indicates no new content, skip it.
	if not feed.has_changed:
		continue

	# List all of this blog's posts in the database.
	posts = db.list_posts(blog)

	# Compare the known posts to the posts present in the feed.
	posts_to_create, posts_to_update = compare_posts(posts, feed.posts)

	# Create each post that needs to be created.
	for post in posts_to_create:
		db.create_post(post)

	# Update each post that needs to be updated.
	for post in posts_to_update:
		db.update_post(post)
```

# Conclusion
