---
date: 2024-09-01
title: "An Infinite io.Reader"
slug: "an-infinite-io-reader"
tags: ["Go"]
draft: true
---

I recently wrote some middleware to ensure that incoming requests to my server have an explicit size limit.
Until I determine this to be too small, I chose to limit request bodies to 4KB.
Go's [net/http](https://pkg.go.dev/net/http) package already includes a utility for this (called [MaxBytesReader](https://pkg.go.dev/net/http#MaxBytesReader)) which makes writing the middleware quite simple:

```go
// Represents a piece of HTTP middleware.
type Middleware func(http.Handler) http.Handler

// Limit the size of the request body to 4KB.
const MaxRequestBodySize = 4 * 1024

// Places an upper limit on the size of every request body.
func LimitRequestBodySize() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

			next.ServeHTTP(w, r)
		})
	}
}
```

Pretty simple!
The next question was obvious: how do I test this?
I need to create an `http.Request` with a body that is larger than 4KB.
I could probably do this by creating a `bytes.Buffer` that is bigger than the limit.
That being said, maybe there is a cleaner way to solve this problem using Go's [io.Reader](https://pkg.go.dev/io#Reader) interface?

## The Interface

What I created an implementation of `io.Reader` that always claimed to have more data available?
This way, I'd never have to increase any buffer sizes if and when I choose to increase the `MaxRequestBodySize`.
As a quick refresher, Go's `io.Reader` is a very simple interface:

```go
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

Implementations of this interface are expected to fill the slice `p` with available data (up to `len(p)`).
Then, it returns the number of bytes written and error (if one occurred).
Given this information, what would it take to simulate an `io.Reader` that never ends?

## The Implementation

The answer is quite simple: just always say that the entire buffer was filled.
If you don't care about what the data is, you can completely ignore it and just return the length of the given buffer.
Since slices are initialized to their zero value, this is essentially a Go implementation of [/dev/zero](https://en.wikipedia.org/wiki//dev/zero).

Let's take a look at the code:

```go
// Simulates an io.Reader of infinite size.
type infiniteReader struct{}

// Make a new infinite io.Reader.
func newInfiniteReader() *infiniteReader {
	r := infiniteReader{}
	return &r
}

// Always report that the entire buffer was filled.
func (r *infiniteReader) Read(p []byte) (int, error) {
	return len(p), nil
}
```

## The Test

With our infinite reader ready to rock, we can write the test.
Simply create a request with our custom `io.Reader` and try to read more than `MaxRequestBodySize` from it (just one extra byte should suffice).
If everything is working correctly, trying to read more than the limit should return a [MaxBytesError](https://pkg.go.dev/net/http#MaxBytesError):

```go
func TestLimitRequestBodySize(t *testing.T) {
	// Prepare the mock ResponseWriter and Request.
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", newInfiniteReader())

	// Prepare the stub handler.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to read more than the request body limit.
		buf := make([]byte, middleware.MaxRequestBodySize+1)
		_, err := r.Body.Read(buf)
		test.AssertErrorAs(t, err, new(*http.MaxBytesError))
	})

	// Wrap the stub handler in the middleware we want to test.
	limitRequestBodySize := middleware.LimitRequestBodySize()
	h := limitRequestBodySize(next)

	// Serve the HTTP request.
	h.ServeHTTP(w, r)

	// Verify that our stub handler was executed.
	rr := w.Result()
	test.AssertEqual(t, rr.StatusCode, http.StatusOK)
}
```

## Conclusion

I thought that this was a pretty clean solution to the problem of "how do create an infinite `io.Reader`".
Since Go's interfaces are often small and simple, it only takes a few lines of code to satisfy them.
While I didn't care about _what_ the data was in this case, it wouldn't take much work to extend this pattern to return something different (like a different value or a repeating sequence).
Gotta love Go!

Thanks for reading.
