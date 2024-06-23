---
date: 2024-06-23
title: "Digimon Cards and Database Indexes"
slug: "digimon-cards-and-database-indexes"
tags: ["Databases"]
draft: true
---

Talk about how my collection sorting related to DB indexes.
I could leave the collection unordered which would make inserts fast but lookups slow.
Instead, I index by (color, type, level) which makes inserts a bit slower but lookups very fast!
If I just did color, it’d be a regular, single-column index.
But instead I chose a [multi-column index](https://www.postgresql.org/docs/current/indexes-multicolumn.html).
Similar to DBs, I must lookup cards “left to right” in order for the index to add value.
