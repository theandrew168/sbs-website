---
date: 2024-07-07
title: "A Better Pattern for Go HTTP Handlers"
slug: "a-better-pattern-for-go-http-handlers"
tags: ["Go"]
draft: true
---

In most web applications, handlers have dependencies.
These could be things like database connection pools, queues, or loggers.
For most of my Go-based web development history, I've embraced the "application struct" pattern as described by Alex Edwards in his "Let's Go" [book series](https://lets-go.alexedwards.net/) to manage these dependencies.
Despite serving my quite well over the past few years, I recently found myself wanting something a bit more flexible.
This post describes the original pattern, its limitations, and how I iterated on it to arrive at something even better (in my opinion, of course).

# Application Struct

This code snippet showcases the original pattern.
Essentially, we gather the union of all of our handler's dependencies together in a centralized, shared "Application" struct.
Then, we attach all handlers to this struct so that they can access the

```go
type Application struct {
    store *storage.Storage
}

func NewApplication(store *storage.Storage) *Application {
    app := Application{
        store: store,
    }
    return &app
}

// Construct the program's routing handler.
func (app *Application) Handler() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("GET /foo", app.handleFoo())
    return mux
}

// Handler for an individual route.
func (app *Application) handleFoo() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        foo, err := app.store.Foo().Read()
        if err != nil {
            util.ServerErrorResponse(w, r, err)
			return
        }

        err = util.WriteJSON(w, 200, foo, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
    }
}
```

# Dependency Closures

Instead of throwing all handler deps into an Application struct and attaching handlers to it, just closure the deps into each route.
This way, the routes only need what they NEED in order to be tested.
Also simplifies construction of individual routes and results in less test code (only need to setup the deps for the current route, not ALL routes).
Matches what I do for middleware, too.

```go
// Construct the program's routing handler.
func Handler(store *storage.Storage) http.Handler {
    mux := http.NewServeMux()
    mux.Handle("GET /foo", HandleFoo(store *storage.Storage))
    return mux
}

// Handler for an individual route.
func HandleFoo(store *storage.Storage) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        foo, err := store.Foo().Read()
        if err != nil {
            util.ServerErrorResponse(w, r, err)
			return
        }

        err = util.WriteJSON(w, 200, foo, nil)
		if err != nil {
			util.ServerErrorResponse(w, r, err)
			return
		}
    }
}
```

# Pros and Cons

## Pros

1. Locality / clarity of deps: everything a route needs is right next to it. It is very clear where routes need storage, for example.
2. No more need for the intermediate application struct, now all handlers are independent and de-coupled. Less state to worry about. Getting the handlers is one call now vs two. This could also enable handler to be split into separate packages.
3. Easier to test handler logic independently of router-level controls (auth, RBAC, etc) and things that are handled by middleware.
4. Simplified testing: each route only needs to setup the deps it needs now. Before, I had to setup everything the application struct needed just to a test a route that might not need ANY of the deps.

## Cons

1. Mild repetition. If many routes have the same deps, you have to pass em into every route. However, I consider this to be more of a facet of reality and is better left obvious. The truth is the same: these routes all need these deps. But shifting them into a central struct and coupling things together just to “save some lines” doesn’t seem worth it to me. I’d rather have a bit of repeated code than the wrong abstraction.
