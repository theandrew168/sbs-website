---
title: "Projects"
---
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

## Skylark ([source](https://github.com/theandrew168/skylark))
Cross-platform [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) emulator.

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
Toy [RISC-V](https://en.wikipedia.org/wiki/RISC-V) assembler.

##### Relevant Skills
* Python programming
* RISC-V [instruction set architecture](https://riscv.org/specifications/)
* Testing Python applications
* Publishing Python packages
* Lexing and parsing of assembly programs
* Binary instruction encoding
* Programming bare-metal devices

## Derzy ([source](https://github.com/theandrew168/derzy))
Old college project exploring real-time graphics and OpenGL.

##### Relevant Skills
* C++ programming
* 3D graphics via OpenGL
* Concepts of real-time rendering
* Build management with [CMake](https://cmake.org/)
