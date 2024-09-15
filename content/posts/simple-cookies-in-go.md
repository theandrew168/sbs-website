---
date: 2024-09-15
title: "Simple Cookies in Go"
slug: "simple-cookies-in-go"
tags: ["Go"]
---

Today's blog post is about cookies!
Not the chocolate chip kind, but the web application state managemant kind.
In short, [cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)! are bits of data (represented as name-value pairs) that web servers can request clients (mostly browsers) to store.
Since HTTP is a stateless protocol, it is useful to have a way to remember pieces of information between requests.
Often, cookies are used for identification and tracking purposes (like those obnoxious pop-ups you see on many sites).
For my projects, however, cookies are primarly used for session management and determining if a user is authenticated.

There are two varieties of cookies: **session** and **permanent**.
Session cookies disappear when the user closes their browser and ends the current session.
Permanent cookies, on the other hand, remain in the browser until they expire (with a lifetime recommended by the server).
The name "permanent" is somewhat confusing since they aren't _actually_ permanent: they just persist bewteen browser sessions until they eventually expire.

## Actions

In my Go-based web projects, I use cookies for session management.
When a user logs in, the server generates a random "session ID" value (using the [crypto/rand](https://pkg.go.dev/crypto/rand) package).
Then, the HTTP response includes a [Set-Cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie) header asking the browser to set a cookie with the provided name, value, max age, and other settings.
In a typical configuration, cookies will be included in subsequent requests to the same server.

In general, I seem to only ever perform one of three actions:

1. Creating **session** cookies
2. Creating **permanent** cookies
3. Creating **expired** cookies

Let's see how each of these actions are implemented as Go functions.

## NewSessionCookie

This is the "foundation" helper upon which the other two are built.
It creates a cookie with the given name and value while providing reasonable defaults for the seldom-changed details.

```go
// Create a session (not permanent) cookie that expires when the user's session ends.
func NewSessionCookie(name, value string) http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",  // This path makes the cookie apply to the whole site.
		Domain:   "",   // An empty domain will default to the server's base domain.
		Secure:   true, // Only send cookies on secure connections (includes localhost).
		HttpOnly: true, // Only send cookies via HTTP requests (not JS).
		// Don't send cookies with cross-site requests but include them when navigating
		// to the origin site from an external location (like when following a link).
		SameSite: http.SameSiteLaxMode,
	}
	return cookie
}
```

Since this cookie doesn't have an [Expires](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#expiresdate) or [MaxAge](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#max-agenumber) field, it will be removed whenever the user closes their browser session.

## NewPermanentCookie

This helper only differs from the previous in that it accepts a "time-to-live" duration which is converted to seconds and returned as the cookie's `MaxAge` value.
The browser will see this and know that the cookie is permanent (until its `MaxAge` is lifetime exceeded).

```go
// Create a permanent (not session) cookie with a given time-to-live.
func NewPermanentCookie(name, value string, ttl time.Duration) http.Cookie {
	cookie := NewSessionCookie(name, value)
	cookie.MaxAge = int(ttl.Seconds())
	return cookie
}
```

Note that cookies have two very similar fields for setting their expiration: `Expires` and `MaxAge`.
Which one should you use?
[TL;DR](https://stackoverflow.com/a/78348791): Use `MaxAge` because it is more reliable (since saying "this expires in N seconds" is less error prone than trying to provide a specific timestamp and hoping that the user's computer's time is correct).

## NewExpiredCookie

The last helper creates a cookie with a `MaxAge` of `-1`.
This tells the browser that the cookie with the given name is instantly expired (and should therefore be deleted).

```go
// Create a cookie that is instantly expired.
func NewExpiredCookie(name string) http.Cookie {
	cookie := NewSessionCookie(name, "")
	cookie.MaxAge = -1
	return cookie
}
```

One important thing to note is that your system shouldn't depend on browsers faithfully deleting expired cookies.
Especially when dealing with sessions, you should be sure to store the same expiration info in your database so that outdated sessions can be deleted server-side.
Otherwise, a malicious client could save their cookie value and simply re-add it themselves (or just make it expire much later than you intended).

## Conclusion

That's it!
Over time, I've converged on these three simple helpers for working with cookies.
I include these functions in almost all web-based Go projects that I work on.
While cookies have many options for tweaking their behavior, I only expose the ability to change `Name`, `Value`, and `MaxAge`.
Other details (such as `Path`, `Domain`, `Secure`, `HttpOnly`, and `SameSite`) are almost never change so I populate these values with reasonable defaults within the `NewSessionCookie` foundation helper.

If these helpers could be useful to you, feel free to try them out.
It not, I hope you at least learned something new!
Thanks for reading.
