---
date: 2024-09-07
title: "Simple Middleware in Go"
slug: "simple-middleware-in-go"
tags: ["Go"]
draft: true
---

Since Go released its [HTTP routing enhancements](https://go.dev/blog/routing-enhancements) in version 1.22, I've been quick to migrate (if you want to learn more about this, Eli Bendersky wrote a [great article](https://eli.thegreenplace.net/2023/better-http-server-routing-in-go-122)).
While Alex Edwards' [Flow](https://github.com/alexedwards/flow) router has served me well (pun intended) for years, I tend to prefer using the standard library whenever possible.
However, one feature that Go's [http.ServeMux](https://pkg.go.dev/net/http#ServeMux) lacks is convenient support for middleware.

## Middleware

Middleware is code that runs between (in the middle of) incoming (or outgoing) HTTP requests and your handlers.
It can be used for all sorts of things: handling panics, adding headers, compressing files, and checking authenication.
Writing these chunks of logic as middleware allows for great flexibility and readability (as we'll soon see).
As far as the code goes, middleware is typically written as a function that both accepts and return Go's most important HTTP interface: the [http.Handler](https://pkg.go.dev/net/http#Handler).

As a type, one could represent middleware like this:

```go
// Represents a piece of HTTP middleware.
type Middleware func(http.Handler) http.Handler
```

Using this type as a focus, I've come to appreciate two specific helpers (`Use` and `Chain`) that are inspired by Mat Ryer's [middleware pattern](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81).
Let's dive into the first one.

## Use

This helper allows us to use multiple pieces of middleware on a single `http.Handler`.
Often, that handler will be an HTTP router of some sort (Go's `http.ServeMux` in my case).

```go
// Apply a sequence of middleware to a handler (in the provided order).
func Use(h http.Handler, mws ...Middleware) http.Handler {
	// Due to how these functions wrap the handler, we apply them
	// in reverse order so that the first one supplied is the first
	// one that runs.
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}
```

Here's what the `Use` helper looks like in action:

```go
// Apply global middleware to all routes.
handler := middleware.Use(mux,
	middleware.RecoverPanic(),
	middleware.CompressFiles(),
	middleware.SecureHeaders(),
	middleware.LimitRequestBodySize(),
	middleware.Authenticate(repo),
)
```

## Chain

The second helper is slightly different than the first.
While it still accepts multiple pieces of middleware and connects them together, it delays the actual application to an `http.Handler` until later.
This can be useful when you have groups of middleware that need to be applied more precisely to specific routes.

```go
// Chain multiple middleware together for delayed application to a handler.
func Chain(mws ...Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		return Use(h, mws...)
	}
}
```

In [Bloggulus](https://bloggulus.com/), for example, some routes require the user to have an account while other require the user to have an account AND be an admin.
Since checking if an account is an admin depends on having an account in the first place, it means that _both_ middleware must be applied.
The `Chain` helper let's us express that dependency once and then reuse the pre-connected pieces.
Here's what that looks like in code:

```go
accountRequired := middleware.AccountRequired()
adminRequired := middleware.Chain(accountRequired, middleware.AdminRequired())
// ...
mux.Handle("GET /blogs", accountRequired(HandleBlogList(find)))
mux.Handle("POST /blogs/follow", accountRequired(HandleBlogFollowForm(repo, find)))
mux.Handle("POST /blogs/unfollow", accountRequired(HandleBlogUnfollowForm(repo, find)))
mux.Handle("GET /blogs/{blogID}", adminRequired(HandleBlogRead(repo)))
mux.Handle("POST /blogs/{blogID}/delete", adminRequired(HandleBlogDeleteForm(repo)))
```

## Ordering

An important detail to note when applying middleware is how their "wrapping" nature affects the order in which they run.
To put things another way, the _first_ middleware that wraps your handler will be the _last_ one to execute (since all subsequently-applied middleware will execute before it).

For example, say we want to apply the following middleware to a handler (in this order).
When request arrive at our server, we want "foo" to run first, then "bar", and then "baz".

```python
mws = [foo, bar, baz]
```

Applying the middleware in the provided yield a call chain that is NOT what we expected.
In fact, it is the opposite of what we want: incoming requests will hit "baz" first.

```python
for mw in mws:
	h = mw(h)

h = baz(bar(foo(h)))
```

However, applying the middleware in the _reverse_ order will achieve what we want.
This new call chain will cause incoming requests to hit "foo" first.

```python
for mw in reversed(mws):
	h = mw(h)

h = foo(bar(baz(h)))
```

## Conclusion

It doesn't take much code to create an incredibly useful and versatile middleware system.
In fact, the [entire file](https://github.com/theandrew168/bloggulus/blob/4de1abef12ca9a9ef3651b8928d3aeac768b0359/backend/web/middleware/middleware.go) containing these helpers is only 30 lines (many of which are comments).
Feel free to try them out in your own Go web application.
If you aren't using Go 1.22's enhanced router yet, try that out as well!
After navigating a small learning curve, I've come to really enjoy using it.
Plus, having fewer dependencies is always nice.

Thanks for reading!
