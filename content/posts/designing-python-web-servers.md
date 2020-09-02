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
One this is for sure: performance must always be measured, not guessed.

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
| [simple1.py](https://github.com/theandrew168/bloggulus/blob/master/design/simple1.py) | 54.23 | 0.84 | 0.92 | 8.10 | 8.07 |
| [simple2.py](https://github.com/theandrew168/bloggulus/blob/master/design/simple2.py) | 52.77 | 0.84 | 0.91 | 7.85 | 7.85 |
| [simple3.py](https://github.com/theandrew168/bloggulus/blob/master/design/simple3.py) | 45.71 | 0.78 | 1.64 | 15.88 | 6.80 |
| [forking1.py](https://github.com/theandrew168/bloggulus/blob/master/design/forking1.py) | 8.70 | 5.60 | 5.14 | 10.14 | 1.29 |
| [forking2.py](https://github.com/theandrew168/bloggulus/blob/master/design/forking2.py) | 81.83 | 0.54 | 0.87 | 7.60 | 12.17 |
| [threading1.py](https://github.com/theandrew168/bloggulus/blob/master/design/threading1.py) | 78.71 | 0.52 | 0.87 | 7.52 | 11.71 |
| [threading2.py](https://github.com/theandrew168/bloggulus/blob/master/design/threading2.py) | 79.06 | 0.51 | 0.88 | 7.73 | 11.76 |
| [processpool.py](https://github.com/theandrew168/bloggulus/blob/master/design/processpool.py) | 48.52 | 0.71 | 1.01 | 8.81 | 7.22 |
| [threadpool.py](https://github.com/theandrew168/bloggulus/blob/master/design/threadpool.py) | 70.04 | 0.62 | 0.87 | 7.49 | 10.42 |
| [nonblocking.py](https://github.com/theandrew168/bloggulus/blob/master/design/nonblocking.py) | 79.33 | 0.53 | 0.63 | 5.23 | 11.80 |
| [async.py](https://github.com/theandrew168/bloggulus/blob/master/design/async.py) | 82.97 | 0.53 | 0.84 | 7.35 | 12.34 |

# Conclusions
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

Aside from these numbers, another important design facet is whether or not a particular design allows me to "bring my own socket".
Some designs open the listening socket themselves which doesn't play nice with how I deploy my apps into production.
I like to get my privileged sockets (for ports 80 and 443) from [systemd](https://www.freedesktop.org/software/systemd/man/systemd.socket.html) which passes the already listening file descriptors to my app.
Due to this, I need to be able to bring my own socket to whichever design I choose.

Moving forward, I'm planning to pick one or two of these designs to take into profiling and optimization.
I really want to see how much closer I can get to NGINX's performance.
I also plan to integrate TLS via Let's Encrypt in a similar manner to Go's [autocert](https://pkg.go.dev/golang.org/x/crypto/acme/autocert?tab=doc) package.

Thanks for reading!
