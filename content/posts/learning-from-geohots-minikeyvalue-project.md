---
date: 2021-01-14
title: "Learning From geohot's minikeyvalue Project"
slug: "learning-from-geohots-minikeyvalue-project"
tags: ["python", "go"]
draft: true
---
I first caught onto [minikeyvalue](https://github.com/geohot/minikeyvalue) while watching a [recording](https://www.youtube.com/watch?v=cAFjZ1gXBxc) of George's programming livestream.
This project was written in [Python](https://www.python.org/) and the design was simple: an HTTP-based interface for storing, getting, and deleting arbitrary content.
Many folks call this type of thing a "distributed key-value store".
The project was originally inspired by [SeaweedFS](https://github.com/chrislusf/seaweedfs) but had the goal of being much, much simpler (ideally less than 1000 lines of code).

The architecture of the program is straightforward: a single "index server" distributes and organizes data between any number of "volume servers".
The actual data could be anything: text, images, [SQLite](https://www.sqlite.org/index.html) databases.
It doesn't matter!
Everything is just bytes at the end of the day.

# Discoveries
The trick that first caught my eye was how the project implemented volume servers.
I cracked open the solitary [volume](https://github.com/geohot/minikeyvalue/blob/master/volume) executable and saw that it was nothing more than a Bash script that started an instance of [NGINX](http://nginx.org/en/).
Even more interesting was how the NGINX config file was handled.
Instead of requiring any prerequisite setup, he just templated some env vars into a string and threw it out to a temporary file via [mktemp](https://man7.org/linux/man-pages/man3/mktemp.3.html)!

This struck a chord of relief with me due to how much effort I've spent in the past trying to effectively coordinate Python-based web applications with NGINX reverse proxies via [Ansible](https://docs.ansible.com/ansible/latest/index.html).
I always believed that NGINX had to be configured separately due to having a mandatory [config file](https://www.nginx.com/resources/wiki/start/topics/examples/full/).
I would write a role to run the [WSGI server](https://gunicorn.org/) on some locally-visible port and then the proxy would listen globally on 80/443.
Each component would also have its own [systemd](https://www.freedesktop.org/wiki/Software/systemd/) unit file.

George's approach slashes this complexity by allowing NGINX to be something that the application manages itself.
As well as using NGINX for the volume servers, my Python version also starts its own reverse proxy (as a child process) for the index server.
I get all the benefits of NGINX (great performance, high extensibility, multiple workers) without the fuss of managing a persistent config file.

# The Volume Servers
A volume stores the data associated with a given key.
NGINX's [WebDAV module](http://nginx.org/en/docs/http/ngx_http_dav_module.html) is used to obtain the desired behavior out of the box: GET reads data, PUT stores data, and DELETE removes data.
George used a simple Bash script to start his NGINX volume servers.
I wanted to explore this idea using Python as the host language in order to be consistent with the index server.
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
    fd, path = tempfile.mkstemp()
    os.write(fd, conf.encode())
    os.close(fd)
    return path
```

I could have also used one of the higher-level, context-managed classes like [NamedTemporaryFile](https://docs.python.org/3/library/tempfile.html#tempfile.NamedTemporaryFile) but that would've added an extra layer of `with` nesting and an explicit call to [flush](https://docs.python.org/3/library/io.html#io.IOBase.flush).
Without the flush, there is a race condition between writing the file's contents and NGINX reading it on startup.

With the config safely written to a temp file, NGINX can now be started.
The only extra flags needed are `-c` to specify where the config is and `-p` to control the runtime prefix dir (where the pidfile goes and things like that).
Now these args can be assembled into a list, passed to [subprocess.run](), and ran until interrupted:
```python
import subprocess

cmd = ['nginx', '-c', conf_path, '-p', prefix_path]
try:
    subprocess.run(cmd)
except KeyboardInterrupt:
    pass
```

Altogether, the source for kicking of these volume servers can be found [here](https://github.com/theandrew168/pymkv/blob/main/volume.py).
As far as performance goes, there is no meaningful difference between George's Bash script and my Python script.
They both run NGINX as a child process and wait for interruption.
The "kickoff" language doesn't really matter here.
Could be Bash, Go, Python, C, or anything else: the result is still the same.

# The Index Server
The index server is the brain of the application.
It keeps track of what keys are currently in the system and where they are located (which volume server and what location).
George's current version of this component contains more features than I wanted to implement so I instead based my design of an [earlier version](https://github.com/geohot/minikeyvalue/blob/a60f742e59b8c11dccaf0f1bef97b605c027a220/server.go).
The core of the index server is a persistent lookup table for keys and the location of their values.

To accomplish this, George used [goleveldb](https://github.com/syndtr/goleveldb) which is a pure-Go implementation of [LevelDB](https://github.com/google/leveldb).
While Python does have its own LevelDB wrapper called [plyvel](https://plyvel.readthedocs.io/en/latest/), I wanted to use something simpler: another standard library gem known as [dbm](https://docs.python.org/3/library/dbm.html).
On most Unix-like systems, the dbm module is implemented as a slim wrapper around the [GNU dbm](https://www.gnu.org.ua/software/gdbm/) key/value database.
The problem that dbm solves is similar to LevelDB but the fact that it ships with Python eases project setup and distribution.

Concurrency is a bit fussy with both options, to be honest.
Neither of them support simultaneous reads/writes from multiple processes.
However, since my [WSGI server](https://www.python.org/dev/peps/pep-3333/) ([Waitress](https://docs.pylonsproject.org/projects/waitress/en/stable/)) uses threads for workers instead of processes, it all works.
As far as I know, CPython's [GIL](https://realpython.com/python-gil/) is "saving" me from multiple-writer corruption in this situation.

# Performance
I ran a bunch of benchmarks on a small digital ocean cluster of small droplets.  
I used one index server and three volume servers all of the same size.  
The results were decent!  
My server was able to hang with George's on the two more important benchmarks (in my opinion).  
On the _most_ important benchmark which simulates more realistic, full-circle client behavior, our servers were matched.
Some other common factor must be the bottleneck at this point (I'm assuming the client-index and index-volume IO).  

# Takeaways
This was really fun!  
I even wrote my own Go version of the index server just to see how it compares.  
I think that Go is the better language choice for this type of problem.  
Sure, I was able to match the performance but it took a whole extra NGINX reverse proxy to save my Python code.  
I do really love this method of dealing with NGINX.  
It is a lot simpler as far as server setup goes and makes deployment easier, too (just have to deploy a single "binary").  
