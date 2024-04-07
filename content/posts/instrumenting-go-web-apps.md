---
date: 2024-04-07
title: "Instrumenting Go Web Apps"
slug: "instrumenting-go-web-apps"
tags: ["Go"]
draft: true
---

Prometheus is an incredible, open-source system for collecting, storing, and analyzing metrics.
Is it written in Go and has a very well-supported client library that can be used to add instrumentation to your own problem.
I always include it in my Go-based web servers (and collect the metrics for monitoring / visualization).
To do this, I use the awesome Prometheus [client_golang](https://github.com/prometheus/client_golang) package.
It is super easy to hook up:

```go
http.Handle("/metrics", promhttp.Handler())
```

Then you can visit your app's `/metrics` endpoint and view all the info (here is a small sample):

```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 6.1654e-05
go_gc_duration_seconds{quantile="0.25"} 9.9282e-05
go_gc_duration_seconds{quantile="0.5"} 0.00011521
go_gc_duration_seconds{quantile="0.75"} 0.000149918
go_gc_duration_seconds{quantile="1"} 0.006138503
go_gc_duration_seconds_sum 5.826040347
go_gc_duration_seconds_count 21144
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 10
```

For years, I've thought that this data included stats about the traffic across HTTP routes.
However, this turned out to _not_ be true.
It just included http stats for the metric endpoint itself (I saw data related to request counts and status codes and figured that the data was already there).
So, how do we fill in this missing info?
I want to know basic things like: how many requests does each endpoint get, what was their latency, and what about status codes.
I found a library to make instrumenting custom routes easy!
It is called [go-http-metrics](https://github.com/slok/go-http-metrics).
It supports all major routes and is simple to configure (check out the [examples](https://github.com/slok/go-http-metrics/tree/master/examples) for more details):

```go
package main

import (
	"net/http"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"
)

func main() {
	mmw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// ...
	mux.Handle("/foo", metricsWrapper.Handler("/foo", mmw, TODO))
	mux.Handle("/bar", metricsWrapper.Handler("/bar", mmw, TODO))
}
```

This example utilizes the std-lib middleware to wrap routes using the standard library HTTP router.
The first param is the route's label.
I've set them to `/foo` and `/bar` here to match the routes.
What about routes using clean URLs that might have params (like `/foo/{id}`)?
Well, you _could_ supply an empty string for a label which will cause the middleware to emit a new series _per unique URL_.
This could result in massive cardinality and you probably shouldn't do it!
I almost feel like you should _never_ using an empty string label: it is better to be explicit.
I'd probably label this route with its match pattern:

```go
mux.Handle("/foo/{id}", metricsWrapper.Handler("/foo/{id}", mmw, TODO))
```

This will show up as the literal string `/foo/{id}` without any values being replaced.
This way, the route will only ever emit a single series (per data point) and will allow you see aggregate counts across all `foo` IDs.
See the [custom example](https://github.com/slok/go-http-metrics/blob/master/examples/custom/main.go) for a better explanation of this issue + solution.
