---
date: 2024-03-24
title: "Testing with Transactions"
slug: "testing-with-transactions"
tags: ["Go", "TypeScript", "Databases"]
---

Most web applications eventually end up with tests that need to interact with a database.
Perhaps your business logic is tightly coupled to the database or maybe you are wanting to test a clearly-defined storage layer.
Either way, a common problem arises: how do you clean up the data used during testing? What should you do with all those scattered rows?

Most of the time, developers will resort to a few common strategies:

1. Cleanup the database after (or before) each integration test runs
   1. Can become very slow depending on the size of your database
   2. Can be optimized (somewhat) by only deleting tables that changed
   3. Could indicate that your tests must be ran serially (concurrency could cause conflicts)
2. Cleanup the database after (or before) _all_ integration tests have been ran
   1. Much quicker than wiping everything for each test
   2. This requires that all tests "play nice" and don't cause any data conflicts
3. Use an ephemeral database for each test
   1. Newer libraries such as [testcontainers](https://golang.testcontainers.org/) have made this much easier
   2. Not a bad idea, but I want a solution with more speed and less complexity

I want to present a better solution to _all_ of these strategies that depends only on a feature that is aleady at our fingertips: the humble **transaction**.

# Commits and Rollbacks

At a high level, database transactions are a feature that let you group a set of changes into an atomic operation that either _all_ succeed (and get committed to the database) or _all_ fail (and get rolled back as though they never happened).
The core idea is this: run your tests inside of a transaction and then roll it back (intentionally) when you are done.

That's it!
With this simple idea, you can verify that your storage layer is working as intended and not have to worry about the cleanup.
Additionally, you can run your database-related tests concurrently without having to worry about conflicts because the changes inside of each test (and therefore each transaction) will never "see" each other.

# In Practice

In order to support arbitrary transactions within the confines of the [repository pattern](https://medium.com/@pererikbergman/repository-design-pattern-e28c0f3e4a30), I like to attach an extra method to my storage class / struct / interface called something like `Transaction`, `WithTransaction` or `Atomically` (I'm still undecided on which name I like best).
This method accepts a function accepts a single argument: an instance of the storage "object".
The function is also expected to either _return_ an error (in languages like Go) or _throw_ an error (in languages like TypeScript) if something goes wrong.
A transaction is started before calling the given function and is automatically rolled back if an error occurs.
Otherwise, it commits the transaction like normal and permanently saves our changes to the database.

In my projects, I like to add a special sentinel error (usually called `ErrRollback` or `RollbackError`) that indicates to the reader that the error is being used _specifically_ to prevent the transaction from committing.
In TypeScript projects, I also check for this error and prevent it from propagating any higher.
As long as the [postgres library](https://github.com/porsager/postgres) sees the error first and rolls back the commit, it has served its purpose.

### TypeScript

This [example](https://github.com/theandrew168/bloggulus-svelte/blob/3c0572d736c2e3d97fa36a56bcbe6b9ec951f254/src/lib/server/storage/storage.ts#L27) is using the [postgres](https://github.com/porsager/postgres) library.
It exposes a [begin](https://github.com/porsager/postgres?tab=readme-ov-file#transactions) utility which automatically rolls back the transaction if any errors are thrown from the lambda it was given.

```ts
export class Storage = {
	// ...

	async transaction(operation: (storage: Storage) => Promise<void>) {
		try {
			// Calling begin() will start a new database transaction and automatically
			// roll it back if any errors are thrown from the given lambda.
			await this.sql.begin(async (tx) => {
				// Create a new instance of the Storage class using this transaction.
				const storage = new Storage(tx);

				// Use the transaction-based Storage class for this atomic operation.
				// If an error occurs within this operation, the transaction will
				// be rolled back.
				await operation(storage);
			});
		} catch (e) {
			// Check for our special sentinel error and prevent it from propagating
			// if found. All other errors should continue to bubble up.
			const isRollbackError = e instanceof RollbackError;
			if (!isRollbackError) {
				throw e;
			}
		}
	}
}
```

This allows me to write code like:

```ts
// create a Foo within a transaction and then roll it back
store.transaction(async (store: Storage) => {
	await store.foo.Create(...)
	throw new RollbackError();
})
```

### Go

This is [example](https://github.com/theandrew168/bloggulus/blob/eba4fee0f7083fe75b56e86bfb9033fe42c11e10/backend/storage/storage.go#L30) is using the [pgx](https://github.com/jackc/pgx) library which requires a slightly more "hands-on" approach to [managing the transaction](https://pkg.go.dev/github.com/jackc/pgx/v5#hdr-Transactions).

```go
type Storage struct = {
	// ...
}

func (s *Storage) WithTransaction(operation func(store *Storage) error) error {
	// Calling the Begin() method on the connection creates a new pgx.Tx
	// object, which represents the in-progress database transaction.
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	// Defer a call to tx.Rollback() to ensure it is always called before the
	// function returns. If the transaction succeeds it will be already be
	// committed by the time tx.Rollback() is called, making tx.Rollback() a
	// no-op. Otherwise, in the event of an error, tx.Rollback() will rollback
	// the changes before the function returns.
	defer tx.Rollback(context.Background())

	// Create a new Storage struct using the pgx.Tx as its Conn.
	store := New(tx)

	// Use the pgx.Tx-based Storage struct for this atomic operation.
	// If an error occurs within this operation, the transaction will
	// be rolled back.
	err = operation(store)
	if err != nil {
		return err
	}

	// If there are no errors, the operation can be committed
	// to the database with the tx.Commit() method.
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
```

This allows me write code like:

```go
// create a Foo within a transaction and then roll it back
store.WithTransaction(func(store *storage.Storage) error {
	store.Foo.Create(...)
	return ErrRollback;
})
```

# Conclusion

There you have it!
This is one of my favorite programming patterns that I've discovered so far in 2024.
If your code is structured such that all database interactions can be managed within explicit transactions, then this trick can be introduced and adapted to fit your needs.
Gone are the days of wiping the database before / after tests run: you can simply rollback and forget!
