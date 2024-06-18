---
date: 2024-06-16
title: "Reinforcing Indirect Joins"
slug: "reinforcing-indirect-joins"
tags: ["Databases"]
---

Recently, I ran into an issue where a very large PostgreSQL query (involving roughly 30 tables) was taking much longer to plan than to execute.
At its worst, it was taking multiple _seconds_ to plan and a few hundred milliseconds to execute.
My problem wasn't really one of query optimization: I was already past that.
This was a problem of query complexity and trying to get the PostgreSQL [planner/optimizer](https://www.postgresql.org/docs/current/planner-optimizer.html) to more quickly arrive at an ideal query.

# Research

Early research made me realize that there is much less content to be found on this specific topic.
There are endless articles explaining how to speed up a query's _execution_ with tools like indexes and clustering.
However, reading material surrounding speeding up the planner/optimizer itself is much more rare.

I was at least able to find some [slides from a presentation](https://www.postgresql.eu/events/pgconfeu2017/sessions/session/1617/slides/9/FromMinutesToMilliseconds.pdf) at the PostgreSQL Conference Europe 2017 that pointed me in the right direction.
I started by messing around with some planner tweaks (like [join_collapse_limit](https://www.postgresql.org/docs/current/runtime-config-query.html#GUC-JOIN-COLLAPSE-LIMIT)) which did help a bit but didn't completely solve my problem: planning was still the query's bottleneck.

I decided to try a new strategy.
I reassembled the query join by join and measured how long the planner/optimizer took at every step.
Once the planning time eventually jumped, I took a moment to think about the data model and how the tables were being joined.
I eventually noticed a pattern that was causing the planner to take extra time and even came up with a name for the fix: **reinforcing indirect joins**.

# An Example

I'll try to explain this concept with an example.
Consider a simple schema with three tables: `post`, `tag`, and `comment`.
One post can have many tags and many comments.
Therefore, both the `tag` and `comment` tables have a foreign key reference back to `post`.

![Simple schema for posts, tags, and comments](/images/20240616/schema.webp)

Now, perhaps you are in a position to ask questions about the relationship between tags and comments.
Do posts with certain tags attract more comments than others?
Are some tags associated with positive sentiment and others with negative?
To start gathering data and answering these questions, you dive into the SQL.
Since you only really care about tags and comments and they can be easily linked together (via `post_id`), let's get it done:

```sql
select
  "tag"."name",
  "comment"."message"
from "tag"
inner join "comment"
  on "comment"."post_id" = "tag"."post_id";
```

# The Problem

Now, based on my what I've seen and experienced, we've just thrown a curveball at the planner/optimizer.
The tables `tag` and `comment` are indirectly linked to each other: they link to each other "through" the `post` table.
Despite this indirect join returning the correct data, PostgreSQL may have trouble deciding how best to optimally execute the query.
With large enough queries, this type of short-cutting can lead to the planner/optimizer taking way longer than expected.

What can we do about it?
Well, we can take this indirect join and reinforce it!

# The Solution

The solution is actually quite simple: add a join to the "missing link" table even if you don't need any data from it.
Here's what the fix looks this looks for our current example:

```sql
select
  "tag"."name",
  "comment"."message"
from "tag"
-- join from "tag" to "comment" through "post"
inner join "post"
  on "post"."id" = "tag"."post_id";
inner join "comment"
  on "comment"."post_id" = "post"."id";
```

To be honest, I'm not exactly sure _why_ this helps, but I have a guess.
PostgreSQL only understands (from a statistics point of view) the relationships between two individual tables.
It knows hows tags relate to posts and how comments relate to posts, but not how tags relate to comments.
By adding posts into the query, the planner/optimizer is able to say "Ahh, I understand: they want the comments assocated with this tag's post".

Perhaps, without the extra join, the database has to search through more plans and therefore takes longer to find one that is optimal.

# Another Example

Another sneaky way that this problem can manifest is when using CTEs to build smaller subsets of data into a larger response.
Maybe you were first querying for tags that meet a certain condition before joining to comments:

```sql
-- build a CTE of relevant tags
with "relevant_tag" as (
  select *
  from "tag"
  where "tag"."name" ilike 'golang'
)
select
  "relevant_tag"."name",
  "comment"."message"
-- find comments related to the relevant tags
from "relevant_tag"
inner join "comment"
  on "comment"."post_id" = "tag"."post_id";
```

Despite having more SQL here to make the relationships slightly less clear, the same problem can occur: PostgreSQL doesn't understand how tags relate to comments.
Thankfully, the fix is the same: just reinforce the indirect join by adding the missing link:

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
-- join from "relevant_tag" to "comment" through "post"
inner join "post"
  on "post"."id" = "relevant_tag"."post_id";
inner join "comment"
  on "comment"."post_id" = "post"."id";
```

# Conclusion

There you have it!
If you ever run into large queries with slow plan times, consider checking for (and reinforcing) any indirect joins.
For smaller queries (like the examples in this post), the planner/optimizer is unlikely to be a bottleneck.
In my experience, this problem only occurs on extra large queries with many tables being joined together.
This could also be a problem that only arose due the specific data model and statistics of the database that I was working with.

At the end of the day, I wish I had been able to find more information about why this slowdown occurs.
I really want to understand why this happens and if anything else can be done to "tell" PostgreSQL about these indirect relationships between tables.
Perhaps this is simply common knowledge: only join on columns that properly reference each other.
Either way, this experience taught me a lesson that I won't soon forget.

Thanks for reading!
