---
date: 2020-06-07
title: "A Multi-Platform Modern OpenGL Demo with SDL2"
slug: "a-multi-platform-modern-opengl-demo-with-sdl2"
tags: ["c", "graphics", "opengl", "sdl2"]
draft: true
---
This post is largely inspired by [Chris Wellons'](https://nullprogram.com/) 2015 blog post about writing a [modern OpenGL demo](https://nullprogram.com/blog/2015/06/06/) that works on all three major desktop platforms (Windows, MacOS, and Linux).
I have come back to his post countless times over the years when looking for guidance on how to build any cross-platform C program.

In Chris' version of this demo he makes use of [GLFW](https://www.glfw.org/) for the window and input handling and uses [gl3w](https://github.com/skaslev/gl3w) for loading OpenGL functions.
I have used GLFW in the past but tend to prefer [SDL2](https://www.libsdl.org/) for its wealth of multimedia features and [polling-based event model](https://wiki.libsdl.org/SDL_PollEvent).
Therefore, in my version, I wanted to use SDL2 for the window and input handling and load the OpenGL functions myself.
Loading OpenGL functions is a fiddly, quirky topic that deserves its own blog post.
However, it doesn't require too much code once the nuances are understood.

# WIP Modern OpenGL
old stuff used to be immediate mode  
old stuff is usually present and linked simply  
new stuff exposes more custom flows  
new stuff has to be loaded dynamically  
newer GLSL lang  

# WIP The Demo
The demo itself is very minimal: a small window containing a rotating red square that exits upon hitting the `ESCAPE`, `Q`, or the close button in the corner.
{{<figure src="/images/sdl2-opengl-demo.png" alt="SDL2 OpenGL Demo">}}

More blerbs here about what it does: single VBO, single VAO, simple shader, etc.

# Structure
Since this is just a demo, the repo is quite sparse.
Aside from the usual `LICENSE` and `README.md`, we have a small `src/` directory containing the code itself and then three [Makefiles](https://pubs.opengroup.org/onlinepubs/009695399/utilities/make.html): one for each major platform.
The default Makefile (the one without a platform suffix) builds the Linux version of the demo.
```
.
├── LICENSE
├── Makefile
├── Makefile.macos
├── Makefile.mingw
├── README.md
└── src
    ├── main.c
    ├── opengl.c
    └── opengl.h
```

The `src/opengl.h` header contains three things: OpenGL function(-pointer) declarations, a function for loading the aforementioned functions, and then some miscellaneous helpers.
The `src/opengl.c` file contains the implementations of the three components listed in the header.
Lastly, `src/main.c` implements the demo: create a window, initialize an OpenGL context, load the modern OpenGL functions, and render a rotating square.

# Building
Build cross-platform C applications isn't that hard but it does require some intentional, upfront planning.
I have to thank Chris and his writings again for this facet of the demo, too.
He has written multiple posts on how to keep C projects simple and portable.
I highly recommend that anyone interested in this space read following posts:
* [Four Ways to Compile C for Windows](https://nullprogram.com/blog/2016/06/13/)
* [How to Write Portable C Without Complicating Your Build](https://nullprogram.com/blog/2017/03/30/)
* [A Tutorial for Portable Makefiles](https://nullprogram.com/blog/2017/08/20/)

### WIP Linux
Building the demo on Linux is a simple native build.
A small number of depenencies are required: git, make, and SDL2.
On a Debian-based Linux system, these can be installed via:
```
sudo apt install git make libsdl2-dev
```

### WIP MacOS
Building the demo on MacOS is a simple native build.
A small number of depenencies are required: git, make, and SDL2.
If using [brew](https://formulae.brew.sh/), these can be installed via:
```
brew install git make sdl2
```

### Windows
Building the demo for Windows requires cross-compiling from either Linux or MacOS.
This process relies on the amazing [mingw-w64](http://mingw-w64.org/doku.php) project which utilizes GCC tooling to build Windows executables and libraries.
Cross-compiling from Linux or MacOS will require the following depenencies: git, make, wget, and mingw-w64.

If using a Debian-based Linux system, these can be install via:
```
sudo apt install git make wget tar mingw-w64
```

If using MacOS with brew installed, the command is:
```
brew install git make wget gnu-tar mingw-w64
```

With all dependencies installed and ready to go, the demo source can be cloned and built.
```
git clone https://github.com/theandrew168/sdl2-opengl-demo.git
cd sdl2-opengl-demo/
make -f Makefile.mingw
```

This will build the demo into a standalone, portable Windows executable: `demo.exe`.
This file can be tested with [Wine](https://www.winehq.org/) or transferred to any modern Windows system for native execution.

One thing you might notice is that the Windows build didn't require installing any SDL2 dependencies.
This is because the build actually downloads (using wget), extracts (using tar), and statically links the pre-built SDL2 libraries directly into the executable.
This is my favorite trick from Chris Wellons' [original demo](https://nullprogram.com/blog/2015/06/06/)!
He used this approach for statically linking GLFW and I realized that it would also work just as well for SDL2.

# WIP Benefits of SDL2
TODO

# WIP Loading OpenGL Functions with SDL2
Make this its own post if it goes deep enough?

Info needed:
* Func typedefs (SDL2/SDL_opengl.h)
* Func ptr declarations (we do this)
* Func ptr impls (we do this)
* Func ptr loader (SDL_GL_GetProcAddress)
* Process to link things up (we do this)

All 3 of the "we do this" things can be simplified with simple macros.
Link to apoorva's post.
Talk about that the union thing?
Talk about why the func protos are so weird. Why not just declare it normal?
- Can't because it isn't known at link time.
- What about extern?
- Doesn't solve anything: still would have to be linkable at link time
