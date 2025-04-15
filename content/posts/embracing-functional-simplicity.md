---
date: 2025-04-14
title: "Embracing Functional Simplicity"
slug: "embracing-functional-simplicity"
draft: true
---

I’ve been thinking a lot recently about software design and how to write better code.
Here, "better" means more correct, more readable, easier to follow, and easier to test.
Two books have been of exceptional influence as of late: [Grokking Simplicity](https://grokkingsimplicity.com/) by Eric Normand and [Architecture Patterns with Python](https://www.cosmicpython.com/) by Harry Percival and Bob Gregory.

The former focuses on functional design: emboldening pure, "data in, data out" functions.
The latter zooms out a bit and considers overall application architecture: how to write code that best represents the domain of your problem.
More content in the same vein includes Brandon Rhodes' [Hoist Your I/O](https://www.youtube.com/watch?v=PBQN62oUnN8) talk and Gary Bernhardt's [Functional Core, Imperative Shell](https://www.destroyallsoftware.com/screencasts/catalog/functional-core-imperative-shell) screencast.

# Lessons

Let's discuss each lesson in a bit more detail.

## 1. Data, Calculations, and Actions

Early on, Eric Normand's [Grokking Simplicity](https://grokkingsimplicity.com/) book explores an idea of partitioning your code into three distinct categories:

1. Data
2. Calculations
3. Actions

Data is what you'd expect: just data!
It could be a string, an array, a map, or anything else where its value comes from its identity.
Calculations and actions are where things get more interesting.
Calculations are "pure" functions that have zero interaction with the outside world.
They never have side effects and the same input always yields the same output.
Lastly, actions are functions that _do_ interact with the world outside of your program.
For example, perhaps they talk to a database or call a remote API.

One of my takeaways is that data is "better" than calculations which are "better" than actions.
Once again, "better" here means more correct, more readable, easier to follow, and easier to test.
Because calculations are pure by definition, they are much simpler to test and verify.
The lack of interaction with external systems means that tests can always be written as: "does this output match what I expect for a given input?".

See, one of the things that makes Bloggulus' sync process so fragile and finnicky is that the logic for fetching RSS feeds, querying the database, and deciding which posts to create / update are all heavily intertwined.
**Prior to refactoring, I couldn't verify the small yet critical decisions that dictate what should happen.**
Sure, the actual fetching of RSS feeds and interacting with the database are decoupled via interfaces (thankfully), but the tests still require a large amount of tiresome faking.

In many other systems I've worked on, these decouping abstractions are _not_ in place and therefore the test are even more difficult to setup, write, and tear down.
Refactoring your code to make more frequent usage of pure functions (calculations) will lead to simpler, more reliable code and smaller, more powerful tests.

## 2. Separating "What" From "How"

This next idea is so simple but it absolutely blew my mind when I first read it!
I immediately starting thinking about all the different places where I could apply it.
Taken from [Chapter 3](https://www.cosmicpython.com/book/chapter_03_abstractions.html) of [Architecture Patterns with Python](https://www.cosmicpython.com/), the idea is this: separate **what** we want to do from **how** to do it.

The example in the book involves syncing files between two directories.
The initial version of the code looks like any developer's first imperative pass at the problem: iterate through each directory and, when necessary, immediately move / copy / delete the files that require action.
This effectively merges the "what" (which files require action) and "how" (actually performing the action) of this process into a single, intertwined step.
One can already imagine the pain of testing this, especially without an easily fake-able filesystem abstraction in place.

Is there a better way?
Can we slice off a pure function (calculation) that handles the most important question of the process: which directores should be moved / copied / deleted?
Their analysis of the problem and its required data yields an amazing result: represent the contents of each directory as a dictionary (data) and then write a pure function (calculation) to decide which files _should_ be moved / copied / deleted.

TODO: More here, maybe an example?

# Impact

Together, these lessons have already led me toward better code in my [Bloggulus](https://bloggulus.com/) project ([source code](https://github.com/theandrew168/bloggulus)): specifically with respect to functional design.
As mentioned earlier, the most critical (and unfortunately most tangled) part of Bloggulus is the “sync” process.
Sometimes, this almost feels like an unintentional anti-pattern: where the most important part of a system is the one that is most tightly-coupled to the outside world and therefore the most difficult to test, verify, and trust...

Anyways, back to the refactoring!
With these two big ideas in mind, let's take a look at the sync process and see how we can make it better (more correct, more readable, easier to follow, and easier to test).

## Before

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

This version is kind of a mess: domain logic is mixed in with reckless abandon, actions and calculations are poorly defined yet heavily intertwined, and "what" vs "how" is confusingly mixed.
How can we isolate the most import question here: which posts should be created and which posts should be updated?

## After

```python
# List all blogs in the database.
blogs = db.list_blogs()

# Process each blog individually.
for blog in blogs:
	# Fetch the blog's RSS / Atom feed.
	feed = fetch_feed(blog)

	# List all of this blog's posts in the database.
	posts = db.list_posts(blog)

	# Compare the known posts to the posts present in the feed.
	posts_to_create, posts_to_update = compare_posts(posts, feed.posts)

	# Create all posts that need to be created.
	db.create_posts(post)

	# Update all posts that need to be updated.
	db.update_posts(post)
```

This version feels much better!
Firstly, each blog's posts are fetched in bulk prior to comparison with the RSS feed data.
Not only is this faster than reading each post one by one, it enables the second improvment.
Instead of the deciding which posts need to be created / updated and then immediately performing that action, we instead utilize a pure function (`compare_posts`) to separate the "what" (create vs update) from the "how"  (interacting with the database).
**In short, the most critical decision that the system has to make is now a pure function.**

Given some known, existing posts and a list of current posts from the RSS feed, we can easily test and verify which ones should be created and which ones should be updated.
Then, the actual database interaction happens _after_ the decision has been mode.
The "what" is separated from the "how".
Another benefit of this approach is that since all new / changed posts are known together, we can batch their creation / updating instead of having to process them invidiually.

If you are curious, the real, Go-based implementation of the `compare_posts` helper can be read [here](https://github.com/theandrew168/bloggulus/blob/72eeeeb2ce949e59d9c5e59a08ff4fe204a1c8c7/backend/service/sync.go#L76).

# Conclusion

Through all of the concepts and ideas I've been reading about lately, two have made the most impact.
The first is to categorize code as either data, calculations, or actions.
The second is to separate the "what" from the "how" when performing state-changing or IO-based operations.

If you're already familiar with these concepts, then consider this a simple refresher.
If not, then I hope you can find as much value and inspiration in these patterns as I have.
There are surely more lessons to learn but I wanted to write about these two, specifically, because they have already proven themselves to be immensely useful in helping me write better, cleaner code.

Thanks for reading!