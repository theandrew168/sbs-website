---
date: 2024-01-28
title: "Utilizing All Cores in a NodeJS Web Application"
slug: "utilizing-all-cores-in-a-nodejs-web-application"
tags: ["NodeJS", "Go"]
---

Coming from Go-based web development to NodeJS, one big thing surpised me: my web server only ever uses one core.
This is because NodeJS is a single-threaded runtime environment.
That being said, NodeJS is still highly concurrent: it uses [modern event loop technology](https://libuv.org/) to implement non-blocking, IO-based concurrency.
NodeJS might not be truly parallel, but it is certainly concurrent.
Knowing this, how is it that NodeJS performs reasonably well in server-side environments?
Is it only ever using a fraction of its available multi-core CPU power?
In some ways, yes, but in practice this isn't usually a bottleneck.

Even a single-threaded instance of NodeJS performs well in IO-bound environments.
In general, web servers spend more time waiting on the internet (connections from users, talking to databases, etc) than they do waiting for the CPU to burn through intense computations.
This means that a single NodeJS thread can keep up with typical web traffic without breaking a sweat.
If your web application was CPU intensive, however, then perhaps you'd start experiencing the downsides of this design.
In such cases, running your application to explicitly utilize all cores could come in handy.

# The Experiment

Consider a minimal, CPU-bound web server program.
It simluates a workload by calculating the square root of ten million numbers before returning "Hello World!".
How does NodeJS handle this scenario vs a more "multi-core aware" web server such as Go's builtin [net/http](https://pkg.go.dev/net/http#hdr-Servers) package?

I performed a load test to examine how these two servers behave with regards to CPU usage.
On a 2-core [Digital Ocean](https://www.digitalocean.com/) droplet, I ran each server and bombarded it with web traffic.
I'm using [hey](https://github.com/rakyll/hey) to send heavy traffic to the the server for 10 seconds: `hey -z 10s <server>`.
While the test is running, I watched the CPU usage characteristics (via [htop](https://htop.dev/)).

### NodeJS Server

```js
import http from "node:http";

console.log("Listening on port 3000...");

http
  .createServer((req, res) => {
    let n = 0;
    for (let i = 0; i < 10000000; i++) {
      n += Math.sqrt(i);
    }

    res.writeHead(200);
    res.end("Hello World!");
  })
  .listen(3000);
```

![NodeJS single-threaded performance](/images/node-single.png)

As you can see, the simple NodeJS server only utilizes a single CPU core.
The 2.7% on the other core comes from other, unrelated operating system tasks.

### Go Server

```go
package main

import (
  "math"
  "net/http"
)

func main() {
  println("Listening on port 3000...")
  http.ListenAndServe(":3000", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    n := 0.0
    for i := 0; i < 10000000; i++ {
      n += math.Sqrt(float64(i))
    }

    w.Write([]byte("Hello World!"))
  }))
}
```

![Go multi-threaded performance](/images/go-multi.png)

The Go version, on the other hand, uses both cores and squeezes the maximum performance out of the server.
This is because Go's runtime load balances work across all cores (via goroutines).
Since the builtin web server handles each incoming request in a separate goroutine, the web traffic is naturally spread across all cores.

# Optimizing NodeJS

What can we do about the NodeJS version?
Enter the [cluster](https://nodejs.org/api/cluster.html) module which is built into NodeJS.
This package allows you to start your web server multiple times: one for each core.

```js
import cluster from "node:cluster";
import http from "node:http";
import { availableParallelism } from "node:os";

if (cluster.isPrimary) {
  // if this is the primary process, start one worker process per core
  const numCPUs = availableParallelism();
  for (let i = 0; i < numCPUs; i++) {
    cluster.fork();
  }
} else {
  console.log("Listening on port 3000...");
  http
    .createServer((req, res) => {
      let n = 0;
      for (let i = 0; i < 10000000; i++) {
        n += Math.sqrt(i);
      }

      res.writeHead(200);
      res.end("Hello World!");
    })
    .listen(3000);
}
```

When I first read this example in the [cluster documentation](https://nodejs.org/api/cluster.html#cluster), I was confused.
How can multiple processes listen on the same port at the same time?
Based on past experience, attempting to do that results in an error: `EADDRINUSE`.
Well, it turns out that NodeJS's http server sets the `SO_REUSEPORT` flag by default on the socket(s) it opens.
_This_ is the secret sauce that allows multiple processes to listen on the same port at the same time.
Additionally, since [Linux 3.9](https://man7.org/linux/man-pages/man7/socket.7.html), this socket option causes incoming TCP connections to be evenly distributed across all listening processes.
This is effectively operating system based load balancing.
Fascinating!

Enough jabbering about sockets, let's see the results.

![NodeJS multi-threaded performance](/images/node-multi.png)

How about that: our CPU-bound NodeJS web server is now making use of all available cores!
With just a few lines of code, we can get around one of NodeJS's largest limitations.
That being said, there are a few things to consider when using this trick.

Since you are running multiple instances of your application, you need to think about how this affects connections to external systems (such as databases).
For example, consider how this would affect database connection pools.
If your application was configured to support up to 50 database connections and you employed this strategy to run multiple instances of it, then each instance will respect this limit independently.
This means that on a four core server, your application could open up to 200 database connections!
Your database could get overwhelmed by receiving more incoming connections that it expects or is configured to handle.
In this scenario, you should think about the multiplicative total across all program instances when setting a connection limit.

# Conclusion

For most NodeJS web applications, the single-threaded limitation will not be an issue.
Between waiting for packets from clients and waiting for rows from a database, your code will be seldom bound by the CPU.
If, by chance, your application _does_ require a heavy amount of CPU usage, the `cluster` module can help out.
By running a separate instance of your program per core and allowing Linux to load balance incoming connections, the full potential of your CPU can be realized.

Thanks for reading!
