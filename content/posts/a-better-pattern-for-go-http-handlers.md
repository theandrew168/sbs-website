---
date: 2024-07-07
title: "A Better Pattern for Go HTTP Handlers"
slug: "a-better-pattern-for-go-http-handlers"
tags: ["Go"]
---

In most web applications, handlers have dependencies.
These could be things like database connection pools, queue clients, or loggers.
For most of my Go-based web development projects, I've embraced the "application struct" pattern as described by Alex Edwards in his "Let's Go" [book series](https://lets-go.alexedwards.net/) (and a [blog post](https://www.alexedwards.net/blog/organising-database-access)) to manage these dependencies.
Despite serving my quite well over the past few years, I recently found myself wanting something a bit more flexible.
This post describes the original pattern, its limitations, and how I iterated on it to arrive at something even better (in my opinion, of course).

# Application Struct

This code snippet showcases the original pattern.
Essentially, we gather the union of all of our handlers' dependencies together in a centralized "application struct".
Then, we attach all handlers to this struct so that they can access the shared resources.

```go
// Keep all handler dependencies in a centralized struct.
type Application struct {
    store *storage.Storage
}

// Constructor for the application struct.
func NewApplication(store *storage.Storage) *Application {
    app := Application{
        store: store,
    }
    return &app
}

// The program's routing handler.
func (app *Application) Handler() http.Handler {
    mux := http.NewServeMux()
    mux.Handle("GET /foo", app.HandleFoo())
    mux.Handle("GET /bar", app.HandleBar())
    return mux
}


// Handler for a simple route without any dependencies.
func (app *Application) HandleFoo() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
}

// Handler for an average route with dependencies.
func (app *Application) HandleBar() http.Handler {
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
    })
}
```

This strategy works quite well.
However, it has a few downsides:

1. All handlers must live in the same package as the application struct.
2. All dependencies must be setup to test each handler: even those that don't need them.

For me, personally, that second detail started to become a nuisance.
Even for a handler that returns static HTML, I have to instantiate a database connection?
I wanted to find a better way.

# Dependency Closures

The solution is quite simple: just pass each handler's dependencies into its creation function.
This creates a closure around the dependency such that the handler can still see and use it even after being returned.
Let's take a look.

```go
// The program's routing handler.
func Handler(store *storage.Storage) http.Handler {
    mux := http.NewServeMux()
    mux.Handle("GET /foo", HandleFoo())
    mux.Handle("GET /bar", HandleBar(store))
    return mux
}

// Handler for a simple route without any dependencies.
func HandleFoo() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
}

// Handler for an average route with dependencies.
func HandleBar(store *storage.Storage) http.Handler {
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
    })
}
```

This is pretty nice!
Gone is that fussy, shared "application struct".
Additionally, each handler's dependencies are now explicitly clear.
If a handler doesn't need any extra resources, then it simply doesn't get any.
This makes testing much simpler and removes all of that "gotta spin up every dependency for every handler" boilerplate.
Plus, if I wanted to, I could now split my handlers up between different packages because they are no longer coupled to a shared struct.

There is at least one downside, though: repetition.
If many handlers have the same dependencies, then you have to pass them into every route.
However, I consider this to be more indicative of reality and is better off being written transparently.
The truth is the same either way: these handlers all need these dependencies.
Shifting them into a central struct and coupling things together just to save a few lines of code doesn’t seem worth it to me.

You know me, I’d rather have a bit of duplication than the wrong abstraction.
