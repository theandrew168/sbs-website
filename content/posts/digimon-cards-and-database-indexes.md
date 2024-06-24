---
date: 2024-06-23
title: "Digimon Cards and Database Indexes"
slug: "digimon-cards-and-database-indexes"
tags: ["Databases"]
---

I'm a big fan of the [Digimon Card Game](https://world.digimoncard.com/) (the one released in 2020).
The games are quick, the mechanics are engaging, and the community is friendly.
Plus, I was a big fan of the Digimon anime as a kid so I have a lot of nostalgia for the characters.
When I first got into the game, I realized that only a few sets had been released so far.
I set my eyes on building and maintaining a complete collection: four copies of each card.
With enough spreadsheets and time, this would allow me build any deck that I wanted.
This was honestly a breath of fresh air coming from [Magic: The Gathering](https://magic.wizards.com/en).

# Digimon Cards

That being said, I've collected quite a few cards.
In addition to the minimum four copies of each card, I also have extra copies of most common and uncommon cards (and even some rares).
Historically, I stored these extra cards unsorted in a series of tall, unstable stacks.
"Managing" these stacks was simple: whenever a new set releases and I acquire more overflow, just throw 'em on the stack!
Recently, however, a couple things happened that motivated me to sort and organize my collection.

First, a friend of mine told me that they were building a binder with all Digimon cards illustrated by a particular artist: [koki](https://wikimon.net/Koki).
Since I have so many extra cards, I told them that I'd be willing to go through my collection and see if I could fill in any missing pieces.
However, since my cards were unsorted, this would mean looking through every single card in my collection.
That seemed pretty tedious and I figured that if I was going to touch every card, I might as well sort them.

Second, I was about to move to a new house and moving those loose stacks of cards didn't seem very fun.
I'd rather buy a large, multi-row trading card box to contain them.
If I was going to buy a nice box, I might as well use the individual rows to hold cards of a particular color.
This way, to find a card of a specific color, I'd only have to look through a subset (roughly 1/6th of all cards).
I ultimately decided to sort my collection by even more than just color: I further sorted them by card type and level.
Now, to find any specific card, I'd have to search through _even fewer_ cards (around 1/36th or so).

# Database Indexes

While sorting my cards and thinking about software (a classic combination), I realized a connection to database indexes.
By sorting my collection, I was effectively transforming it into an [index](https://www.postgresql.org/docs/current/indexes-intro.html).
If I had chosen to only sort by color, then it’d mirror a regular, single-column index.
Instead, sorting by color, type, and level maps more closely to a [multi-column index](https://www.postgresql.org/docs/current/indexes-multicolumn.html).

If this configuration was represented in PostgreSQL, the schema and index would look something like:

```
CREATE TABLE card (
  id TEXT PRIMARY KEY,
  name TEXT,
  color TEXT,
  type TEXT,
  level TEXT
);

CREATE INDEX card_idx ON card (color, type, level);
```

Instead of having to sequentially scan through all of my cards to find a specific one, I could instead jump to the specific color + type + level and only have to search through a handful.
By indexing my cards, lookups are elevated from a `Seq Scan` to an `Index Scan`.
However, my collection is now also susceptible to one of the downsides of indexes: increased insertion time.

Before, adding new cards was instant: just add them to the stack.
Now the process is more involved.
I need to first use the index to lookup where a specific card should go and _then_ insert it.
This adds a mild amount of overhead to every single card that gets added to the box.
It takes work to maintain an index but if you have a read-heavy workload, it is probably worth it.
