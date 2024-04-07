---
date: 2024-04-07T23:00:00
title: "Instrumenting Go Web Apps"
slug: "instrumenting-go-web-apps"
tags: ["Go"]
---

[Prometheus](https://prometheus.io/) is an incredible open-source system for collecting, storing, and analyzing system metrics.
In addition to the server program, numerous [client libraries](https://prometheus.io/docs/instrumenting/clientlibs/) have been written to simplify the process of exposing metrics in your own projects.
When writing Go-based web apps, I always include [Go's client library](https://github.com/prometheus/client_golang) so that information about the program can be collected, visualized, and monitored.
As outlined in the [official guide](https://prometheus.io/docs/guides/go-application/), it is very easy to get these basic metrics up and running:

```go
package main

import (
    "net/http"

    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe("127.0.0.1:2112", nil)
}
```

After running the program (`go run main.go`), you can visit its `/metrics` endpoint to view all metrics and values at the current point in time:

```
curl http://localhost:2112/metrics
```

Here is a small sample of the data you'll be able to view, collect, and store.
Each line (that isn't a comment) represents a single [time series](https://prometheus.io/docs/concepts/data_model/).
Notice how different tags (`code="200"`, `code="500"`, etc) yield independent, unique series (this will be important later):

```
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 6
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 121216
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```

When sampled periodically and stored for weeks (or even months), you can begin to assemble a sense of how your application is behaving.
For example, you can observe garbage collection stats, memory usage, and number of active goroutines.
You can even track HTTP request traffic but it requires a bit of extra setup.

# HTTP Metrics

For years, I thought that Go's client library included HTTP traffic stats for _all_ routes.
However, this is not the case: it only includes information about the `/metrics` endpoint.
I must've seen a few keywords in the data (like `request` and `code`) and jumped to the conclusion that all routes were tracked by default.

So, how do we fill in this missing info?
I want per-route information about request volume, latency, and status codes.
Instead of having to write fully-custom metrics from scratch, we can utilize another library: [go-http-metrics](https://github.com/slok/go-http-metrics).
It supports all major routers, is simple to configure, and includes a bunch of useful [examples](https://github.com/slok/go-http-metrics/tree/master/examples).
At a high level, you integrate this library by wrapping its middleware around the HTTP handlers that you want to track.

In pseudocode:

```go
package main

import (
	"net/http"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	metricsMiddleware "github.com/slok/go-http-metrics/middleware"
	metricsWrapper "github.com/slok/go-http-metrics/middleware/std"
)

func main() {
	// initialize the middleware with a metrics recorder
	mmw := metricsMiddleware.New(metricsMiddleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// wrap the HTTP handlers that you want to track
	http.Handle("/foo", metricsWrapper.Handler("/foo", mmw, fooHandler))
	http.Handle("/bar", metricsWrapper.Handler("/bar", mmw, barHandler))
}
```

### Handler Parameters

Tracking each route requires using the [specific wrapper](https://github.com/slok/go-http-metrics/tree/master?tab=readme-ov-file#framework-compatibility-middlewares) that is tailored to whichever HTTP router you are using (I'm using Go's standard library router).
Each wrapper exposes a `wrapper.Handler()` function that accepts three parameters:

1. A label for the route
2. The initialized middleware
3. Your downstream HTTP handler

While the latter two are fairly straightward, the first parameter is worthy of some extra explanation.
In the above example, I've labeled the metrics as `/foo` and `/bar` to match their routes.
This works perfectly for routes that don't ever change.
That being said, what about routes that have dynamic parameters (like `/foo/{id}`)?
There are a couple ways to handle this.

### Option 1: Unique Label per URL

First, you could supply an empty label which will cause the library to emit a new series _per unique URL_.

```go
mux.Handle("/foo/{id}", metricsWrapper.Handler("", mmw, specificFooHandler))
```

This might seem convenient but is actually **quite risky and could result in a cardinality explosion**.
This means that every unique URL visited will result in a separate series of metrics.
The [Prometheus docs](https://prometheus.io/docs/practices/instrumentation/#do-not-overuse-labels) warn against this pattern:

> Each labelset is an additional time series that has RAM, CPU, disk, and network costs.
> Usually the overhead is negligible, but in scenarios with lots of metrics and hundreds of labelsets across hundreds of servers, this can add up quickly.
> As a general guideline, try to keep the cardinality of your metrics below 10, and for metrics that exceed that, aim to limit them to a handful across your whole system.

In my experience, providing the wrapper with an empty label and letting it generate a new series for each unique URL is rarely what I'm after.

### Option 2: Specific Label per Route

Instead, I prefer labeling dynamic routes with their generic match pattern:

```go
mux.Handle("/foo/{id}", metricsWrapper.Handler("/foo/{id}", mmw, specificFooHandler))
```

With this approach, all requests will be aggregated under the literal label `/foo/{id}` regardless of what specific ID was provided.
To put it another way: the route will only ever emit a single series (per metric) and will group the data related to each distinct `foo` together.
See the [custom example](https://github.com/slok/go-http-metrics/blob/c472df028d97fa53f3e99c760831d55908541bba/examples/custom/main.go#L51-L57) for another explanation of this issue how specific labels can be used to solve it.

# Conclusion

This post offers a brief introduction to [Prometheus](https://prometheus.io/) metrics and how you can easily expose them from any Go-based web application.
Despite not being tracked by default, I showed how the [go-http-metrics](https://github.com/slok/go-http-metrics) library can be used to track HTTP request metrics for _all_ routes in your application (not just the `/metrics` endpoint).
Lastly, you should be careful about how you label your per-route metrics in order to avoid excessive cardinality.
Thanks for reading!
