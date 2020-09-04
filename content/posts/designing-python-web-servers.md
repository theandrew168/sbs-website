---
date: 2020-09-02
title: "Designing Python Web Servers"
slug: "designing-python-web-servers"
tags: ["python", "server", "design", "networking"]
---
I've spent the last few days designing and benchmarking extremely minimal pure-Python web servers.
There is a tech myth / superstition that exposing python web servers to the internet is a bad idea but I've seen little to no evidence as to why this is supposedly the case.
Are they too slow?
Are they too insecure for some reason?
One thing is for sure: performance must always be measured, not guessed.

# Designs
The designs I tested can be grouped into 4 broad categories: sequential, forking, threading, and asynchronous.

### Sequential
Sequential servers are the simplest and most straightfoward.
They process all clients in the order they arrive, one at a time.
This means that every client has to wait in a potentially very long line before they get service.
I expect this approach to be mediocre.

### Forking
Forking servers spawn a new operating system process for each client that arrives.
This means that clients don't have to wait on each other but more of the server's resources will be consumed.
Additionally, spawning new processes tends to have noticable overhead when many clients connect rapidly.
I expect this approach to be visibly affected by this overhead.

### Threading
Threading servers are similar to forking servers except that they spawn a thread to handle each client instead of a process.
Threads tend to have lower overhead than processes but aren't quite as independent when it comes to grabbing time on the CPU.
They share memory as well as some other server resources.
I expect this approach to perform well since web servers tend be bound by IO, not the CPU.

### Asynchronous
Asynchronous servers double down on the IO-heavy workload of web servers.
They keep all client connections together in such a way that no connections are ever waited on.
Only when data is ready to be read from a client does the server "wake up" and process it.
I expect this approach to perform very well since it is tailored to fit IO-bound workloads.

# Benchmark
The servers were tested whilst running on a minimal [DigitalOcean](https://www.digitalocean.com/) Droplet (1 CPU, 1GB RAM).
All of them are using a TCP listen backlog of 128.
Printing each request to stdout is disabled during the benchmark.
For the actual load testing, I used a tool called [hey](https://github.com/rakyll/hey).

It performs 2000 requests spread across 50 concurrent connections:
```
hey -n 2000 -c 50 http://bloggulus.com
```

# Results
I started with a load test against [NGINX](https://nginx.org/en/) as a baseline.
I then measured all of the other servers in the same way.

RPS stands for "Requests per second" and latencies are measured in seconds.
The column "% NGINX (RPS)" is extra cool in my opinion.
It clearly shows how much room for improvement exists within each design.

| Server | RPS | Mean Latency | Mode Latency | Worst Latency | % NGINX (RPS) |
| --- | --- | --- | --- | --- | --- |
| [NGINX](https://github.com/nginx/nginx) | 672.36 | 0.07 | 0.11 | 0.64 | 100.00 |
| [sequential1.py](https://github.com/theandrew168/bloggulus/blob/master/design/sequential1.py) | 54.23 | 0.84 | 0.92 | 8.10 | 8.07 |
| [sequential2.py](https://github.com/theandrew168/bloggulus/blob/master/design/sequential2.py) | 52.77 | 0.84 | 0.91 | 7.85 | 7.85 |
| [sequential3.py](https://github.com/theandrew168/bloggulus/blob/master/design/sequential3.py) | 45.71 | 0.78 | 1.64 | 15.88 | 6.80 |
| [forking1.py](https://github.com/theandrew168/bloggulus/blob/master/design/forking1.py) | 8.70 | 5.60 | 5.14 | 10.14 | 1.29 |
| [forking2.py](https://github.com/theandrew168/bloggulus/blob/master/design/forking2.py) | 81.83 | 0.54 | 0.87 | 7.60 | 12.17 |
| [threading1.py](https://github.com/theandrew168/bloggulus/blob/master/design/threading1.py) | 78.71 | 0.52 | 0.87 | 7.52 | 11.71 |
| [threading2.py](https://github.com/theandrew168/bloggulus/blob/master/design/threading2.py) | 79.06 | 0.51 | 0.88 | 7.73 | 11.76 |
| [processpool.py](https://github.com/theandrew168/bloggulus/blob/master/design/processpool.py) | 48.52 | 0.71 | 1.01 | 8.81 | 7.22 |
| [threadpool.py](https://github.com/theandrew168/bloggulus/blob/master/design/threadpool.py) | 70.04 | 0.62 | 0.87 | 7.49 | 10.42 |
| [nonblocking.py](https://github.com/theandrew168/bloggulus/blob/master/design/nonblocking.py) | 79.33 | 0.53 | 0.63 | 5.23 | 11.80 |
| [async.py](https://github.com/theandrew168/bloggulus/blob/master/design/async.py) | 82.97 | 0.53 | 0.84 | 7.35 | 12.34 |

# Analysis
Pretty interesting stuff!
Some designs were awful, some were middle-of-the-road, and some were decent.
As expected, the simple sequential servers were quite average and all floated around that 50 RPS range.
The `forking1.py` design was the worst by far.
Spawning a new process per client connection really does add a large amount of overhead.
I'm really not sure why it was so much worse than `forking2.py` which uses Python's builtin [ForkingTCPServer](https://docs.python.org/3/library/socketserver.html#socketserver.ForkingTCPServer) (this supposedly does the same thing).
By watching the actually system processes, the latter didn't seem to be creating a new process for each client as it describes.
I wonder if it has some sort of internal "max processes" limit to prevent things from getting out of hand.

All of the threading servers did well.
I'm not surprised by this given that a web server's workload tends to be very IO-bound.
If the workload was more CPU-bound then I'd expect to the threading designs rank a bit lower.
The `nonblocking.py` and `async.py` servers also performed well for a similar reason: an IO-bound workload.
These two designs both multiplex the socket IO handling around which clients are ready for action _right now_.
The two "pool" designs did decent but not as well as their unpooled counterparts.
The overhead of managing the pool must be coming into play here.

# Existing Servers
Just for completeness, I wanted to see how a few existing Python-based web servers stacked up against my quick and dirty designs.
Unsurprisingly, they beat me by a long shot!
Most of them did, at least.

| Server | RPS | Mean Latency | Mode Latency | Worst Latency | % NGINX (RPS) |
| --- | --- | --- | --- | --- | --- |
| [NGINX](https://github.com/nginx/nginx) | 672.36 | 0.07 | 0.11 | 0.64 | 100.00 |
| [Gunicorn[sync]](https://docs.gunicorn.org/en/latest/design.html#sync-workers) | 46.31 | 0.99 | 0.95 | 8.23 | 6.89 |
| [Gunicorn[gevent]](https://docs.gunicorn.org/en/latest/design.html#async-workers) | 338.36 | 0.14 | 0.09 | 0.43 | 50.32 |
| [Waitress](https://docs.pylonsproject.org/projects/waitress/en/stable/index.html) | 449.07 | 0.10 | 0.09 | 0.42 | 66.79 |

The synchronous Gunicorn did about as well as my own sequential servers which is to be expected.
Gunicorn with [gevent](http://www.gevent.org/) based workers did extremely well.
Gevent is a well-optimized implementation of asynchronous IO with its core event loop written in C (as either [libev](http://software.schmorp.de/pkg/libev.html) or [libuv](http://libuv.org/)).
I figured that it would be a high performer and I wasn't wrong.

Waitress was the underdog for me but... Wow!
I was really blown away by its performance.
It even gives NGINX a run for its money by smashing through 67% of its requests per second.
Even more amazing is that Waitress is a pure-Python implementation with zero external dependencies.
Talk about good engineering.

I'm definitely going to spend some time reading the source of Waitress to understand more about how they achieve such impressive performance.
Thankfully, they have a well-written [design overview](https://docs.pylonsproject.org/projects/waitress/en/stable/design.html) which will be a good place to start.

# Conclusion
Moving forward, I plan to use to Waitress as my web server of choice for upcoming Python web projects.
The performance is great (as we've seen) and it supports a "bring your own socket" model of initialization.
This is important to me because the apps I deploy to production get their privileged listen sockets (on ports 80 and 443) from [systemd](https://www.freedesktop.org/software/systemd/man/systemd.socket.html) as raw file descriptors.

Thanks for reading!
