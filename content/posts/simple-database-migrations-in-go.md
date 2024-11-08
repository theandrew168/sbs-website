---
date: 2024-11-08
title: "Simple Database Migrations in Go"
slug: "simple-database-migrations-in-go"
tags: ["Databases", "Go"]
draft: true
---

Talk about how simple it can be to handle database migrations when do you don’t need many features.
Just list files in a dir, check which ones are missing, and apply them in order (each within a transaction).
Roughly 50 LoC can save you from an entire dependency.
Be sure to name your migrations with an ascending order: either simple numbers (`0001`, `0002`, etc) or timestamps.

Example [migration code](https://github.com/theandrew168/bloggulus/blob/main/backend/postgres/migrate.go) from Bloggulus.

## Steps

1. Ensure "migration" table exists
2. List migrations that have already been applied (from the table)
3. List migrations that should be applied (from the filesystem / embed)
4. Determine missing migrations (simple set operation)
5. Sort missing migrations to preserve order
6. For each missing migration:
   1. Begin a transaction
   2. Apply the migration
   3. Update the "migration" table
   4. Commit the transaction

## Code

```go
func Migrate(conn Conn, files fs.FS) ([]string, error) {
	ctx := context.Background()

	// create migration table if it doesn't exist
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS migration (
			id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		)`)
	if err != nil {
		return nil, err
	}

	// get migrations that are already applied
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

	// get migrations that should be applied
	subdir, _ := fs.Sub(files, "migrations")
	migrations, err := fs.ReadDir(subdir, ".")
	if err != nil {
		return nil, err
	}

	// determine missing migrations
	var missing []string
	for _, migration := range migrations {
		name := migration.Name()
		if _, ok := existing[name]; !ok {
			missing = append(missing, name)
		}
	}

	// sort missing migrations to preserve order
	sort.Strings(missing)

	// apply each missing migration
	var applied []string
	for _, name := range missing {
		sql, err := fs.ReadFile(subdir, name)
		if err != nil {
			return nil, err
		}

		// apply each migration in a transaction
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(context.Background())

		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			return nil, err
		}

		// update migration table
		_, err = tx.Exec(ctx, "INSERT INTO migration (name) VALUES ($1)", name)
		if err != nil {
			return nil, err
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return nil, err
		}

		applied = append(applied, name)
	}

	return applied, nil
}
```
