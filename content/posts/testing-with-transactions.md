---
date: 2024-03-22
# date: 2024-03-24
title: "Testing with Transactions"
slug: "testing-with-transactions"
draft: true
---

Most web applications eventually end up with tests that need to interact with a database.
Sometimes your logic is tightly integrated to your database or maybe you are wanting to test your repository / storage layer.
Either way, a common problem arises: how do you clean up the data used to test the database?

Most of the team, developers will resort to a few common strategies:

1. Same database
   1. No cleanup done
   2. Cleanup after all tests (ensure no conflicts)
   3. Cleanup after each test
2. Separate database (usually on a different port)
   1. No cleanup done
   2. Cleanup after all tests (ensure no conflicts)
   3. Cleanup after each test
3. Ephemeral database (container)
   1. No cleanup necessary

I want to present a better solution to _all_ of these strategies that depends only on the ability to run all storage operations within a controlled transactions.
Furthermore, I'll showcase how this affects architecture / tests in both TypeScript and Go projects.

# Rollback

The core idea is this: run your tests inside of a transaction and then simply roll it back when you are done.
That's it!
Such a simple idea that adds a ton of value.
You get to verify that your storage layer is working as intended (can CRUD stuff like you should) and not worry about the cleanup at all.
Additionally, you can run your database tests concurrently without having to worry about conflicts because the data in each test / transaction will never "see" each other.

# Examples

I usually structure my repo / storage layer as a struct with a public member for each table / model.
These members are abstract storage interfaces for each type.

In pseudocode:

```go
type Foo struct {
	id UUID
	// ...
}

type FooStorage interface {
	Create(foo Foo) error
	Read(id UUID) (Foo, error)
	List(limit, offset int) ([]Foo, error)
	Update(foo Foo) error
}

type Storage struct {
	conn database.Conn

	Foo FooStorage
}

func New(conn database.Conn) *Storage {
	s := Storage{
		conn: conn,

		// implements FooStorage by means of a Postgres database
		Foo: NewPostgresFooStorage(conn),
	}
	return &s
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
	// If an error occurs within this operation, the transction will
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
store := storage.New(conn)

// create a Foo for real
store.Foo.Create(...)

// create a Foo in a transaction and then roll it back
store.WithTransaction(func(store *storage.Storage) error {
	store.Foo.Create(...)
	return errors.New("skip commit")
})
```

### TypeScript

https://github.com/theandrew168/bloggulus-svelte/blob/main/src/lib/server/storage/storage.ts

### Go

https://github.com/theandrew168/bloggulus/blob/main/backend/storage/storage.go
