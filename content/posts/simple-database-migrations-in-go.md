---
date: 2024-11-09
title: "Simple Database Migrations in Go"
slug: "simple-database-migrations-in-go"
tags: ["Databases", "Go"]
draft: true
---

Historically, database migrations were a problem I'd solve with a third-party library.
I eventually discovered, however, that handing migrations is actually quite simple if you don't need many features.
At a high level, the process is as follows: list out the migration files, check which ones are missing, and apply them in order (each within a transaction).
Implementing this logic can be done in a single function with roughly 50 lines of code.
Additionally, doing so can save your project from an entire (and sometimes complex) dependency.

## The Steps

The process of analyzing and applying migrations is actually quite straightforward.
First you must determine which migrations _should_ be applied... then apply them!
Not too fussy, right?

In reality, there is a bit more nuance to it (but not too much).
Here is a finer-grained list of the necessary steps:

1. Ensure "migration" table exists
2. List migrations that have already been applied
   - These come from the "migration" table
3. List migrations that should be applied
   - These come from a project directory (and can be [embedded](https://pkg.go.dev/embed))
   - Be sure to name your migrations with an ascending order: either simple numbers (`0001`, `0002`, etc) or timestamps.
4. Determine missing migrations
5. Sort missing migrations to preserve order
6. For each missing migration:
   1. Begin a transaction
   2. Apply the migration
   3. Update the "migration" table
   4. Commit the transaction

That's it!
While certainly not the shortest list of steps, I find it to be quite manageable.
Furthermore, writing this code saves me from having to introduce a new (and possibly complex) dependency into my application's workflow.

## The Code

Here is the code in its entirety (you can check out the [real thing](https://github.com/theandrew168/bloggulus/blob/main/backend/postgres/migrate.go) in Bloggulus).

```go
func Migrate(conn Conn, files fs.FS) ([]string, error) {
	ctx := context.Background()

	// 1. Ensure "migration" table exists
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return nil, err
	}

	// 2. List migrations that have already been applied
	rows, err := conn.Query(ctx, "SELECT name FROM migration")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	existing := make(map[string]bool)
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		existing[name] = true
	}

	// 3. List migrations that should be applied
	subdir, _ := fs.Sub(files, "migrations")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return nil, err
	}

	// 4. Determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := existing[name]; !ok {
			missing = append(missing, name)
		}
	}

	// 5. Sort missing migrations to preserve order
	sort.Strings(missing)

	// 6. For each missing migration:
	var applied []string
	for _, name := range missing {
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return nil, err
		}

		// 1. Begin a transaction
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(context.Background())

		// 2. Apply the migration
		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			return nil, err
		}

		// 3. Update the "migration" table
		_, err = tx.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return nil, err
		}

		// 4. Commit the transaction
		err = tx.Commit(context.Background())
		if err != nil {
			return nil, err
		}

		applied = append(applied, name)
	}

	return applied, nil
}
```

## Conclusion

So, should you study this code, understand how it works, and roll your own migrations?
Or should you go for one of the existing, populate libraries (like foo or bar)?
The decision is ultimately up to you.
I chose the former approach probably 5 or 6 years ago and it hasn't let me down (or ever failed).
Sure, I'm missing some features (like "down" migration and migrating to a certain point), but I haven't needed them yet.

Using more code than the current problem requires (often via a dependency) feels unnecessary to me.
I want to understand how my system works and not be beholden to another libraries documentation (if I can avoid it).
Perhaps this is too [Not Invented Here](https://en.wikipedia.org/wiki/Not_invented_here) of me.
But again: this approach has been successful and met all of my migration needs for years.
I think that means I made the right call for the scale of projects that I often work with.

Thanks for reading!
