---
date: 2021-01-14
title: "Learning From geohot's minikeyvalue Project"
slug: "learning-from-geohots-minikeyvalue-project"
tags: ["Python", "Go"]
---

I first became aware of [minikeyvalue](https://github.com/geohot/minikeyvalue) while watching a [recording](https://www.youtube.com/watch?v=cAFjZ1gXBxc) of George's programming livestream.
This project was written in [Python](https://www.python.org/) and the design was simple: an HTTP-based interface for storing, getting, and deleting arbitrary content.
Many folks call this type of thing a "distributed key-value store".
The project was originally inspired by [SeaweedFS](https://github.com/chrislusf/seaweedfs) but had the goal of being much, much simpler (ideally less than 1000 lines of code).

The architecture of the program is straightforward: a single "index server" distributes and organizes data between any number of "volume servers".
The actual data could be anything: text, images, [SQLite](https://www.sqlite.org/index.html) databases.
It doesn't matter!
Everything is just bytes at the end of the day.

# Discoveries

The trick that first caught my eye was how the project implemented volume servers.
I cracked open the solitary [volume](https://github.com/geohot/minikeyvalue/blob/master/volume) executable and saw that it was nothing more than a Bash script that started an [NGINX](http://nginx.org/en/) in the foreground.
Even more interesting was how the NGINX config file was handled.
Instead of requiring any prerequisite setup, he just templated some env vars into a string and threw it out to a temporary file via [mktemp](https://man7.org/linux/man-pages/man3/mktemp.3.html).

This blew my mind!
I've spent many hours in the past trying to effectively coordinate Python-based web applications and NGINX reverse proxies via [Ansible](https://docs.ansible.com/ansible/latest/index.html).
I always believed that NGINX had to be configured separately due to having a mandatory [config file](https://www.nginx.com/resources/wiki/start/topics/examples/full/).
I would write a role to run the [WSGI server](https://gunicorn.org/) on some locally-visible port and then the proxy would listen globally on 80/443.
Each component would also have its own [systemd](https://www.freedesktop.org/wiki/Software/systemd/) unit file.

George's approach slashes this complexity by allowing NGINX to be something that the application manages itself.
As well as using NGINX for the volume servers, my Python version also starts its own reverse proxy (as a child process) for the index server.
I get all the benefits of NGINX (great performance, high extensibility, multiple workers) without the fuss of managing a persistent config file.

# The Volume Servers

A volume stores the data associated with a given key.
NGINX's [WebDAV module](http://nginx.org/en/docs/http/ngx_http_dav_module.html) is used to obtain the desired behavior out of the box: GET reads data, PUT stores data, and DELETE removes data.
George used Bash NGINX but I wanted to explore this idea using Python as the host language in order to be consistent with the index server.
Plus, I knew that Python ships with two modules that are perfect for the job: [tempfile](https://docs.python.org/3/library/tempfile.html) and [subprocess](https://docs.python.org/3/library/subprocess.html).

The first step was to generate the NGINX config as a string.
Python has two main methods of [string interpolation](https://www.python.org/dev/peps/pep-3101/).
The first uses percent placeholders while the second uses curly braces.
I decided to go with the first method in order to avoid having to escape all of the curlies that accompany an NGINX config.
Once I've got the config as string, I need a place to write it.
This is where the tempfile module comes in:

```python
import os
import tempfile

def nginx_temporary_config_file(conf):
    fd, path = tempfile.mkstemp()  # create the temp file
    os.write(fd, conf.encode())  # write the NGINX config to it
    os.close(fd)  # close it
    return path  # return its path
```

I could have also used one of the higher-level, context-managed classes like [NamedTemporaryFile](https://docs.python.org/3/library/tempfile.html#tempfile.NamedTemporaryFile) but that would've added an extra layer of `with` nesting and an explicit call to [flush](https://docs.python.org/3/library/io.html#io.IOBase.flush).
Without the flush, there is a race condition between writing the file's contents and NGINX reading it on startup.

With the config safely written to a temp file, NGINX can now be started.
The only extra flags needed are `-c` to specify where the config is and `-p` to control the runtime prefix dir (where the pidfile goes and things like that).
These args can be assembled into a list, passed to [subprocess.run](https://docs.python.org/3/library/subprocess.html#subprocess.run), and ran until interrupted:

```python
import subprocess

cmd = ['nginx', '-c', conf_path, '-p', prefix_path]
try:
    subprocess.run(cmd)
except KeyboardInterrupt:
    pass
```

The full code for starting these volume servers can be found [here](https://github.com/theandrew168/pymkv/blob/main/volume.py).
As far as performance goes, there is no meaningful difference between George's Bash script and my Python script.
They both run NGINX as a child process and wait for interruption.
The "kickoff" language doesn't really matter here.
Could be Bash, Go, Python, C, or anything else: the result is the same.

# The Index Server

The index server is the brain of the application.
It keeps track of what keys are currently in the system and where they are located (on which volume server and at what path).
George's current version of this component contains more features than I wanted to implement so I instead based my design on an [earlier version](https://github.com/geohot/minikeyvalue/blob/a60f742e59b8c11dccaf0f1bef97b605c027a220/server.go).

The core of the index server is a persistent lookup table for keys and the location of their values.
To accomplish this, George used [goleveldb](https://github.com/syndtr/goleveldb) which is a pure-Go implementation of the [LevelDB](https://github.com/google/leveldb) database.
While Python does have its own LevelDB wrapper called [plyvel](https://plyvel.readthedocs.io/en/latest/), I wanted to use something simpler: another standard library gem known as [dbm](https://docs.python.org/3/library/dbm.html).
On most Unix-like systems, the dbm module is implemented as a slim wrapper around the [GNU dbm](https://www.gnu.org.ua/software/gdbm/) key/value database.
The problem that dbm solves is similar to LevelDB but the fact that it ships with Python eases project setup and distribution.

Concurrency is a bit fussy with both options, to be honest.
Neither of them support simultaneous reads/writes from multiple processes.
However, since my [WSGI](https://www.python.org/dev/peps/pep-3333/) server of choice (more on that coming up) uses threads for workers instead of processes, everyone is happy.
As far as I know, CPython's [GIL](https://realpython.com/python-gil/) is "saving" me from multiple-writer corruption in this situation.

The source for the index server can be found [here](https://github.com/theandrew168/pymkv/blob/main/index.py).
One thing I found out from early performance tests was that while some WSGI servers could handle the load, others couldn't.
[Bjoern](https://github.com/jonashaag/bjoern), for example, held its own but [Waitress](https://docs.pylonsproject.org/projects/waitress/en/stable/) returned errors under the pressure.
I've been in this situation before and knew that the smart traffic handling and request buffering that NGINX offers as a [reverse proxy](https://en.wikipedia.org/wiki/Reverse_proxy) can make almost any WSGI server sing.
Plus, since I already had code to run arbitrary NGINX servers, this was a simple inclusion!

When ran, the index server starts up its own NGINX background process to act as an external reverse proxy to its internal Waitress WSGI server.
With a little help from [contextlib](https://docs.python.org/3/library/contextlib.html) I was able to construct a very simple context manager to start NGINX in the background while the index server continues forward.
Once the server is told to stop, the context manager exits and the reverse proxy shuts down with it.

```python
from contextlib import contextmanager
import subprocess
import waitress

@contextmanager
def run_in_background(cmd):
    proc = subprocess.Popen(cmd)  # start the child process
    yield  # wait til the context manager exits
    proc.terminate()  # send SIGTERM to the child
    proc.wait()  # wait for the child to exit

with run_in_background(['nginx', '-c', ...]):
    # ...
    try:
        waitress.serve(app, host='127.0.0.1', port=8080)
    except KeyboardInterrupt:
        pass
```

# Performance

I wanted to see how the two versions of this project performed in a small set of benchmarks.
I'm not a fan of testing locally so I spun up a small cluster of $5 [Digital Ocean](https://www.digitalocean.com/) droplets in order to introduce proper network latency.
I used one index server and three volume servers for all of the benchmarks and ensured that each version was tested under the same circumstances.

| **Benchmark**     | **Description**                            |
| ----------------- | ------------------------------------------ |
| **fetch missing** | fetch a key that doesn't exist             |
| **fetch present** | fetch a key that exists                    |
| **thrasher.go**   | simulate create, read, and delete behavior |

Benchmark results measured in "requests per second":

| **Benchmark**     | **My Python Version** | **George's Go Version** |
| ----------------- | --------------------- | ----------------------- |
| **fetch missing** | 1397                  | 2152                    |
| **fetch present** | 1012                  | 941                     |
| **thrasher.go**   | 99                    | 102                     |

The results were decent!
My server was able to hang with George's on the two more important benchmarks (in my opinion).
On the _most_ important benchmark which simulates realistic client behavior, our servers were matched.
Some other common factor must be the bottleneck here.
I'm assuming that the time spent waiting on network communication between components is the dominating factor.

# Takeaways

The full project with source code and documentation can be found [here](https://github.com/theandrew168/pymkv).
Overall, this was really fun!
After finishing things up, I decided to write my own [Go version](https://github.com/theandrew168/pymkv/blob/main/index.go) of the index server just to see how the language feels.
To be honest, I think that Go is the better language choice for this type of problem.
Sure, I was able to match the performance but it took an extra NGINX reverse proxy to save my straightforward Python.

I really do love this method of dealing with NGINX, though.
It greatly simplifies the setup of Python-based web apps and makes deployment easier, too.
No more separate config files, just a single "binary" that manages NGINX itself.
Amazing!
