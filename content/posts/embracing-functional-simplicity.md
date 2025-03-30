---
date: 2025-03-30
title: "Embracing Functional Simplicity"
slug: "embracing-functional-simplicity"
draft: true
---

I’ve been thinking a lot recently about software design and how to write better code.
Here, "better" means more correct, more readable, easier to follow, and easier to test.
Two books have been of exceptional influene lately: [Grokking Simplicity](https://grokkingsimplicity.com/) by Eric Normand and [Architecture Patterns with Python](https://www.cosmicpython.com/) by Harry Percival and Bob Gregory.
The former focuses on functional design: emboldening pure, data-in-data-out functions.
The latter hoists the ideas up to a higher level and thinks about overall application architecture and how to write code that best represents the domain of your problem.

# Lessons

These ideas have already led me toward better code in my [Bloggulus](https://bloggulus.com/) project ([source code](https://github.com/theandrew168/bloggulus)): specifically with respect to functional design.
For some background context, the nastiest yet most critical part of Bloggulus is the “sync” process.
This is the code that enumerates all blogs in the system, checks their RSS feeds for updates, and then syncs any new / updated posts with the database.
It almost seems like an unfortunate pattern, sometimes: where the most important part of a system is the one that is most tightly-coupled to the outside world and therefore difficult to test, verify, and trust.

## Data, Calculations, and Actions

Early on, Eric Normand's [Grokking Simplicity](https://grokkingsimplicity.com/) book explores an idea of bucketing your code into three separate categories:

1. Data
2. Calculations
3. Actions

Data is what you'd expect: just data!
It could be a string, an array, a map, or anything else where its value comes from its identity.
Calculations and actions are where things get more interesting.
Calculations are "pure" functions that have zero interaction with the outside world.
They never have side effects and the same input yields the same output.
Lastly, actions are functions that _do_ interact with the world outside of your program.
For example, perhaps they talk to a database or call a remote API.

Eric argues (and I agree) that data is "better" than calculations which are "better" than actions.
Once again, "better" here means more correct, more readable, easier to follow, and easier to test.
Because calculations are pure by definition, they are much simpler to test and verify.
The lack of interaction with external systems means that tests can always be written as: "does this output match what I expect for a given input?".

See, one of the things that makes Bloggulus' sync process so nasty and tangled is the logic for fetching RSS feeds, querying the database, and deciding which posts to create / update are all heavily intertwined.
Prior to refactoring, I couldn't verify the small yet critical decisions that dictate what should happen.
Sure, the actual fetching of RSS feeds and interacting with the database are decoupled via interfaces (thankfully), but the tests still require a large amount of tiresome faking.
In many other systems I've worked on, these decouping abstractions are _not_ in place and therefore the test are even more difficult to setup, write, and teardown.

Refactoring your code to make more frequent usage of pure functions (calculations) will lead simpler, more reliable code and smaller, more powerful tests.

## Separating "What" from "How"

This next idea is so simple but it absolutely blew my mind when I first read it!
I immediately starting thinking about all the different places where I could apply it.
Taken from [Chapter 3](https://www.cosmicpython.com/book/chapter_03_abstractions.html) of [Architecture Patterns with Python](https://www.cosmicpython.com/), the idea is this: separate **what** we want to do from **how** to do it.
The example in the book involves syncing files between two directories.
The initial version of the code looks like any developer's first pass at the problem: iterate through each directory and, when necessary, immediately move / copy / delete the files that require action.
This effectively merges the "what" (which files require action) and "how" (actually performing the action) into a single, intertwined step.

We can already imagine the pain of testing this, especially without any fake-able filesystem abstractions in place.
Is there a better way?
Can we slice of a of pure function (calculation) that handles the most important question of the process: which directores should be moved / copied / deleted?
Their analysis of the problem and its required data yields an amazing result: represent the contents of each directory as a dictionary (data) and then write a pure function (calculation) to decide which files _should_ be moved / copied / deleted.

# Impact

With these two big ideas in mind, let's take a look at the sync process and see how we can make it better.

## Before

This version is kind of a mess: domain logic is mixed in with reckless abandon, actions and calculations are poorly defined yet heavily intertwined, and "what" vs "how" is confusingly mixed.
How can we isolate the most import question here: which posts should be created and which posts should be updated?

```python
# List all blogs in the database.
blogs = db.list_blogs()

# Process each blog individually.
for blog in blogs:
	# Fetch the blog's RSS / Atom feed.
	feed = fetch_feed(blog)

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

## After

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

# Process each syncable blog individually.
for blog in blogs:
	# Fetch the blog's RSS / Atom feed.
	feed = fetch_feed(blog)

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
