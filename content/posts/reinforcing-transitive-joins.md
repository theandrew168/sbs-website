---
date: 2024-06-16
title: "Reinforcing Transitive Joins"
slug: "reinforcing-transitive-joins"
draft: true
---

Recently, I ran into an issue where a very large query (joining roughly 30 tables together) was taking much longer to plan than to execute (multiple seconds to plan and a few hundres milliseconds to execute).
My problem wasn't really one of query optimization: I was already past that.
This was a problem of query complexity and trying to get the PostgreSQL [planner/optimizer](https://www.postgresql.org/docs/current/planner-optimizer.html) to more quickly arrive at an ideal query.
I started by messing around with some planner tweaks like [join_collapse_limit](https://www.postgresql.org/docs/current/runtime-config-query.html#GUC-JOIN-COLLAPSE-LIMIT) which did help (knocked a few seconds off of the plan) but didn't completely solve my problem: planning was still the query's bottleneck.

So, I reassembled the query join by join and waiting for the planning time to jump.
Once it did, I took some time think about the data and how the tables were being joined.
I eventually noticed a pattern that was causing the planner to take extra time and even came up with a name for the fix: reinforcing transitive joins.

# An Example

I'll explain this concept with an example.
Consider a simple schema with three tables: `post`, `tag`, and `comment`.
Once post can have many tags and many comments.
Therefore, both the `tag` and `comment` tables have a foreign key reference back to `post`.

![Simple schema for posts, tags, and comments](/images/20240616/schema.webp)

Now, perhaps you are in a position to ask questions about the relationship between tags and comments.
Do posts with certain tags attract more comments than others?
Are some tags associated with positive sentiment and others with negative?
To start gathering data and answering these questions, you dive into some SQL.
Since you only really care about tags and comments and they can easily linked together (via `post_id`), let's get it done:

```sql
select
  "tag"."name",
  "comment"."message"
from "tag"
inner join "comment"
  on "comment"."post_id" = "tag"."post_id";
```

# The Problem

Now, based on my own anecdotal experience, we've just thrown a curveball at the PostgreSQL query optimizer.
The tables `tag` and `comment` are "transitively" linked to each other: they link to each other via the `post` table.
While we know that quite intuitively, postgres can't quite come to the same decision as quickly.
Instead, with large enough queries, this type of short-cutting can lead to the planner taking way more time than necessary.
So what can we do about it?
Well, we can take this transitive join and reinforce it!

# The Solution

The fix is actually quite simple: add a join to the "missing link" table even if you don't need any data from it.
Here' how this looks in our current example:

```sql
select
  "tag"."name",
  "comment"."message"
from "tag"
inner join "post"
  on "post"."id" = "tag"."post_id";
inner join "comment"
  on "comment"."post_id" = "post"."id";
```

Now instead of comparing the tables based on `tag.post_id == comment.post_id`, we route the join through that shared post: `tag.post_id == post.id && post.id == comment.post_id`.
I'm not exactly sure _why_ this helps, but I have a guess.
PostgreSQL only understands (from a statistics point of view) the relationships between two individual tables.
It knows hows tags relate to posts and comments relate to posts, but not how tags related to comments.
By adding posts into the query, the planner is able to say "Ahh, I understand: they want the comments assocated with this tag's post".
Without the extra join, perhaps PG has to guess at how the tables are related and it therefore takes longer to converge.

# Another Example

Another sneaky way that this problem can manifest is when using CTEs to build smaller subsets of data into a larger response.
Maybe you were first queries for tags that meet a certain condition and _then_ joining to comments:

```sql
with "relevant_tag" as (
  select *
  from "tag"
  where "tag"."name" ilike 'golang'
)
select
  "relevant_tag"."name",
  "comment"."message"
from "relevant_tag"
inner join "comment"
  on "comment"."post_id" = "tag"."post_id";
```

Despite having some more SQL to muddle the waters and spread things out, the same problem can occur: PG doesn't understand how tags related to comments.
Thankfully, the fix is the same here: just reinforce the transitive join by adding the missing link:

```sql
with "relevant_tag" as (
  select *
  from "tag"
  where "tag"."name" ilike 'golang'
)
select
  "relevant_tag"."name",
  "comment"."message"
from "relevant_tag"
inner join "post"
  on "post"."id" = "relevant_tag"."post_id";
inner join "comment"
  on "comment"."post_id" = "post"."id";
```

# Conclusion

Talk about how this probably only applies at scale with larger queries.
Small examples like this are so minimal that you it won't likely make a difference.
This is mainly something to keep in mind if you ever find yourself up against slow query plans but fast executions.
Why does this happen?
Does PostgreSQL not understand these transitive relationships?
Or are there fancy things you can do with table stats to help it understand these scenarios?
If you know, please reach out and let me know.

Thanks for reading!
