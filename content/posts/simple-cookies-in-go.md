---
date: 2024-09-14
title: "Simple Cookies in Go"
slug: "simple-cookies-in-go"
tags: ["Go"]
draft: true
---

Today's blog post is about [cookies](https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies)!
Not the baked treat kind, but the web application state managemant kind.
In short, cookies are bits of data (in key-value pairs) that web servers can request clients (mostly browsers) to keep track of.
Since HTTP is a stateless protocol, it is useful to have a way to track persistent state in web clients.

There are two varieties of cookies: session and permanent.
Session cookies disappear when the user closes their browser.
Permanent cookies remain until they expire (which is a property controlled by the server).
The name "permanent" is somewhat confusing since they aren't actually permanent, they simpler persist bewteen browser sessions until they eventually expire.

In my Go-based web projects, I use cookies for session management.
When a user logs in, the server generates a random "session ID" value.
Then, the HTTP response includes a [Set-Cookie](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie) header to tell the browser to store the value.
Depending on their configuration, cookies will be included in subsequent request to the same server.

I seem to only ever find myself doing one of three actions:

1. Creating session cookies
2. Creating permanent cookies
3. Creating expired cookies

Let's see how each of these actions are implemented as Go functions.

## NewSessionCookie

This is the "foundation" helper upon which the other two are built.
It hard-codes the common defaults for seldom-changed cookie details.
It creates a cookie with the name and value.

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

## NewPermanentookie

This helper only differs from the previous in that it access a "time-to-live" duration which is converted to seconds and returned as the cookie's `MaxAge` value.
The browser will see this and know that the cookie is permanent (until its `MaxAge` is exceeded).

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
TL;DR: Use `MaxAge` because it avoids unexpected time zone surprises.

## NewExpiredCookie

The last helper creates a cookie with a `MaxAge` of `-1`.
This tells the browser that the cookie with the given name is instantly expired (and will therefore be deleted).

```go
// Create a cookie that is instantly expired.
func NewExpiredCookie(name string) http.Cookie {
	cookie := NewSessionCookie(name, "")
	cookie.MaxAge = -1
	return cookie
}
```

One important thing to note that you shouldn't "trust" that a browser will delete expired cookies.
Especially when dealing with sessions, you should be sure to store the same expiration info in your database so that outdated sessions can be deleted server-side.
Otherwise, a malicious client could save their cookie value and simply re-add it themselves (or just make it expire much later than you intended).

## Conclusion

That's it!
Over time, I've converged on these three simple helpers for working with cookies.
I include these functions in almost all web-based Go projects that I work on.
While cookies have many options for tweaking their behavior, I only usually change `Name`, `Value`, and `MaxAge` (time until a permanent cookie expires).
Other details (such as `Path`, `Domain`, `Secure`, `HttpOnly`, and `SameSite`) are almost always the same so I hard-code their values within these helpers.

If these could useful to you, feel free to try them out.
It not, I hope you at least learned something new!
Thanks for reading.
