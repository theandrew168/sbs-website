---
title: "Projects"
---
## Bloggulus ([website](https://bloggulus.com), [source](https://github.com/theandrew168/bloggulus))
[Bloggulus](https://bloggulus.com) is a "meta blog" that aggregates numerous other blogs into a single location.
Previously, I had been using a [Firefox extension](https://addons.mozilla.org/en-US/firefox/addon/feeder/) to track my blogs but I wanted something more personalized.
I also found myself spending too much time on [Reddit](https://bloggulus.com) and [Hacker News](https://bloggulus.com).
Like many others, I’d even sometimes read only the comments and not the actual content.
That’s definitely a waste of time!
So, I created a website that intentionally lacks a social commenting system and only pulls in posts from the blogs that I trust.

I also wanted an excuse to build something with [Flask](https://flask.palletsprojects.com/en/1.1.x/), [SQLite](https://www.sqlite.org/index.html), and [Peewee](phttp://docs.peewee-orm.com/en/latest/).
A special thanks goes out to [Charles Leifer's amazing blog](http://charlesleifer.com/) for inspiring me to use these tools together.
Together, they are definitely greater than the sum of their parts.

Bloggulus is hosted on a small [DigitalOcean](https://www.digitalocean.com/) droplet (the $5/month one).
[Gunicorn](https://gunicorn.org/) is used as the app's [WSGI server](https://www.python.org/dev/peps/pep-3333/) and [NGINX](http://nginx.org/) is deployed as a reverse proxy.
The domain is registered through [Google Domains](https://domains.google/) and the TLS cert comes from [Let's Encrypt](https://letsencrypt.org/).
Once an hour, an automated process checks each blog's [RSS feed](https://en.wikipedia.org/wiki/RSS) for new content.
All blog posts are stored and indexed in an [FTS (full-text search)](https://www.sqlite.org/fts3.html) table for quick and easy searching.

The frontend portion of Bloggulus is pretty simple.
It is all rendered server-side (via Flask’s [Jinja](https://jinja.palletsprojects.com/en/2.11.x/) templating) and styled using [Tailwind CSS](https://tailwindcss.com/).
Web design is definitely one of my weaker stills at the moment but practice is the only way to improve!
I intentionally brought together simple colors and hard lines to let the blog posts be the focus.
The color scheme did accidentally end up looking kind of like [Amazon’s](https://www.amazon.com/), though.

##### Relevant Skills
* Python programming
* Python web application development
* Flask development
* Full-stack development
* Deployment automation
* Linux server administration
* Data modeling with SQLite
* Content indexing and full-text search
* NGINX reverse proxy configuration
* TLS setup via Let's Encrypt

## JamQL ([website](https://jamql.com), [source](https://github.com/theandrew168/jamql))
This project started out as an idea for a way to create Spotify playlists based on a SQL-like query (hence the QL in JamQL).
My plan was to integrate Spotify's [Search API](https://developer.spotify.com/documentation/web-api/reference/search/search/) with SQLite's [virtual table](https://www.sqlite.org/vtab.html) functionality in order to bridge the gap.
I knew that an alternative to using SQLite would be to write my own simple query language and a corresponding lexer and parser for it.
However, I felt like reinventing that wheel wouldn't be worth the trouble and that SQLite's virtual tables would be a perfect fit regardless.

After getting into the weeds a bit, I realized that having users interact with the app purely through a query language wasn't great UX.
Especially on mobile, restricting the interface to one giant text box wasn't very fun _or_ usable.
Instead, I decided to take things in a different direction.
I remembered enjoying iTunes [Smart Playlist](https://support.apple.com/guide/itunes/use-smart-playlists-itns3001/mac) feature back in the day and wondered if that'd be a more viable approach.
So, the current version of the site functions similar to those smart playlists: you can search tracks and create playlists based on a set of filters.

I do plan to circle back to the query language idea at some point in the future.
I'd likely take the route of writing my own minimal language as well as a lexer and parser for it.
Since I'm still very much a novice when it comes to web design, I want to get the current interation working and looking great before tackling additional ways to jam.

##### Relevant Skills
* Go programming
* Go web application development
* Testing Go applications
* Introductory web design
* Web application security
* Full-stack development
* Fully automated deployments
* Zero downtime deployments
* Production-ready deployments
* Privileged port management via [systemd](https://www.freedesktop.org/software/systemd/man/systemd.socket.html)

## Squeaky ([example](https://github.com/theandrew168/squeaky/blob/master/examples/breakout.scm), [source](https://github.com/theandrew168/squeaky))
Squeaky is a proof of concept programming language for making games.
It is an implementation of [Scheme](https://en.wikipedia.org/wiki/Scheme_(programming_language)) (subset of R5RS) with a focus on writing simple and portable multimedia applications.
The initial goal was to support many features (2D graphics, 3D graphics, and networking to name a few).
However, after achieving the minimal viable product of cross-platform, line-based graphics, I felt as though I had learned enough and made enough progress to wrap this project up and move on.
I plan to study more complex compiler topics such as code generation and executable formats (ELF vs PE) and may return to this project in the future.

Back in college a good buddy of mine took a course on programming languages.
I knew nothing more about that topic than the fact that multiple different languages exist.
I wanted to know more.
How is a programming language made?
Can they even be made?
I thought that programming languages were somehow atomic: unable to be broken down into smaller pieces.
But then my experience with Python taught me something interesting: Python is _written_ in C ([CPython](https://github.com/python/cpython) is, at least).

That realization truly broke down a big mysterious wall for me.
Programming languages are just programs that read text and do stuff.
Some interpret the text literally while others "compile" the text's meaning into a leaner, compressed form.
There are many books on language design and implementation.
Some notable ones include [The Dragon Book](https://en.wikipedia.org/wiki/Compilers:_Principles,_Techniques,_and_Tools) and [SICP](https://en.wikipedia.org/wiki/Structure_and_Interpretation_of_Computer_Programs).
My fascination with the latter and with game development led to this proof of concept.

##### Relevant Skills
* C programming
* Scheme programming
* Cross-compiling for non [Unix-like](https://en.wikipedia.org/wiki/Unix-like) operating systems
* Build management with [POSIX-compatible](https://pubs.opengroup.org/onlinepubs/009695399/utilities/make.html) Makefiles
* Programming language design
* Functional API design
* Garbage collection

## Skylark ([source](https://github.com/theandrew168/skylark))
One of my earliest memories of being truly amazed by a programmer was first seeing Bisqwit's videos about writing an NES emulator ([Part 1](https://www.youtube.com/watch?v=y71lli8MS8s) and [Part 2](https://www.youtube.com/watch?v=XZWw745wPXY)).
His efficient and straight-forward approach to solving the problem was truly unlike anything I'd ever seen.
It was nothing short of incredible.
After seeing his videos I was inspired to get into emulator development myself.

I figured that starting with an NES emulator was bit too zealous and instead looked for something with a smaller scope.
A quick googling led me to [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8).
This simple, minimal instruction set is commonly referred to as the "Hello World" of emulator projects.
I was able to scrap together an implementation after a few weeks of iterating.
Despite being technically finished, it had numerous flaws that I didn't really realize until multiple years of C programming later.

For anyone curious, the name Skylark is a reference to the character [Chip Skylark](https://en.wikipedia.org/wiki/Chris_Kirkpatrick) from the old Nickelodeon show [The Fairly OddParents](https://en.wikipedia.org/wiki/The_Fairly_OddParents).

##### Relevant Skills
* C programming
* Cross-compiling for non [Unix-like](https://en.wikipedia.org/wiki/Unix-like) operating systems
* Build management with [POSIX-compatible](https://pubs.opengroup.org/onlinepubs/009695399/utilities/make.html) Makefiles
* Testing C applications
* Hardware [emulation](https://en.wikipedia.org/wiki/Emulator)
* Binary instruction decoding
* 2D graphics via [SDL2](https://www.libsdl.org/)
* Separation of IO-based and pure functional logic

## SimpleRISCV ([source](https://github.com/theandrew168/simpleriscv))
Assembly programming has always been one of those arcane, fundamental skills that I always attributed to older, wiser, expert programmers who have been writing "code" since its inception.
I remember thinking how crazy it was that early Pokemon games were written in assembly.
The idea of programming for a specific CPU and set of peripherals sounds pure and unemcumbered.
Sure, you may give up the potential for cross-CPU support but if that isn't a requirement then there is nothing lost.
You might even be able to make better use of your hardware because you don't have to rely on "common denominator" abstrations.

My desire to learn assembly led me to the [RISC-V](https://en.wikipedia.org/wiki/RISC-V) instruction set architecture (ISA).
The combination of simplicity, newness, and freedom seemed like a perfect fit for a beginner.
I ordered a couple of cheap RISC-V chips and started reading!
My search for a no-frills and portable RISC-V assembler came up short: I really didn't want to fuss with compiling a custom GCC toolchain.
Instead, I wrote a bit of Python to satisfy my very small RISC-V needs.

SimpleRISCV is that project and while correct (as far as I know), it is very bare-bones.
It only supports assembling single files and does absolutely no linking.
The program is simply a RISC-V instruction encoder with a thin layer of label and constant handling above it.
Nonetheless, it helped me get started with the platform and gave me an up-close experience with the ISA.

##### Relevant Skills
* Python programming
* RISC-V [instruction set architecture](https://riscv.org/specifications/)
* Testing Python applications
* Publishing Python packages
* Lexing and parsing of assembly programs
* Binary instruction encoding
* Programming bare-metal devices

## Derzy ([source](https://github.com/theandrew168/derzy))
One of the first topics to really get me interested in programming was real-time rendering.
After completing The Cherno's [Flappy Bird tutorial](https://www.youtube.com/watch?v=527bR2JHSR0) in single sitting, I knew that I wanted to learn more about this interesting [OpenGL](https://en.wikipedia.org/wiki/OpenGL) library.
Despite not having a specific use-case in mind, I created this project as more of a sandbox for learning more about graphics.
The Flappy Bird clone was a 2D game and I wanted to take it to the next level: 3D.

Along with being my first introduction to 3D graphics, this project was also my first time programming in C++.
I figured that this was a natural progression since the two languages I'd learned beforehand were C and Java.
Many books were read in order to familiarlize myself with all this new content: OpenGL, C++, and CMake (the build system).
After a few hundreds commits I had achieved most of the goals that I set out to accomplish.

I wanted to supporting loading 3D models and rendering in them a decent-looking scene with a skybox.
I wanted a camera that I could fly around the scene and see how the lighting looked from different angles.
Fancy shader effects were also on the menu: reflection, refraction, and bloom.
In hindsight, the design of this project was a bit too presumptuous with its abstractions.
After painting myself into a corner with every attempt at structuring the code, I decided that something was systematically flawed with my approach and transitioned to other programming languages and paradigms.

##### Relevant Skills
* C++ programming
* 3D graphics via OpenGL
* Concepts of real-time rendering
* Build management with [CMake](https://cmake.org/)
