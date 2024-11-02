---
date: 2024-11-02
title: "Migrating Numeric IDs to UUIDs"
slug: "migrating-numeric-ids-to-uuids"
tags: ["Databases"]
---

How do you migrate a database from numeric IDs to UUIDs?
If the data has no relationships, then the command is simple.
If there are relationships, however, then you need to utilize add intermediate values to the new table to track the old links and convert them to the new UUIDs.
This problem definitely gets more complex the larger your data model becomes.
At some point, it might not be worth it.
For Bloggulus, however, it was doable in a single, medium-complexity [migration](https://github.com/theandrew168/bloggulus/blob/main/migrations/0005_convert_ids_to_uuid.sql).

Why are UUIDs better (domain models can be created w/o asking the DB for new IDs)?
Enables a "functional core" design and easier DDD models.
Do they come with tradeoffs (space)?
UUIDs take 16 bytes while serial takes 4 bytes (or 8 for bigserial).
This means that each row will take up 12 more bytes.
IMO, this is a small price to pay for domain modeling clarity.

```sql
-- before
id SERIAL PRIMARY KEY

-- after
id UUID DEFAULT gen_random_uuid() PRIMARY KEY

-- easy case: without relationships
ALTER TABLE tag
	ALTER COLUMN id DROP DEFAULT,
	ALTER COLUMN id SET DATA TYPE UUID USING (gen_random_uuid()),
	ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE migration
	ALTER COLUMN id DROP DEFAULT,
	ALTER COLUMN id SET DATA TYPE UUID USING (gen_random_uuid()),
	ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- hard case: with relationships
-- create new parent table w/ UUID PK and old_id column
-- create new child table w/ UUID PK
-- populate new parent table (keep old ID around)
-- populate new child table (join to parent on old_id to get new UUID PK)
-- drop old_id column from parent
```
