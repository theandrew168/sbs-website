---
date: 2025-03-23
title: "A ParallelMap Implementation in Go"
slug: "a-parallel-map-implementation-in-go"
draft: true
---

```go
func ParallelMap[T any](concurrency int, items []T, fn func(T)) {
	// Use a weighted semaphore to limit concurrency.
	sem := semaphore.NewWeighted(int64(concurrency))

	// Perform tasks in parallel (up to "concurrency" at once).
	for _, item := range items {
		// For each item, attempt to acquire a "task slot" from the semaphore.
		// If the maximum number of tasks are already running, this will block.
		sem.Acquire(context.Background(), 1)

		// Kick off a goroutine that runs our function on an item and then
		// releases a "task slot" back to the semaphore for future tasks.
		go func(item T) {
			defer sem.Release(1)
			fn(item)
		}(item)
	}

	// Wait for all tasks to finish.
	sem.Acquire(context.Background(), SyncConcurrency)
}
```
