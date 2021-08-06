---
date: 2020-06-07
title: "A Multi-Platform Modern OpenGL Demo with SDL2"
slug: "a-multi-platform-modern-opengl-demo-with-sdl2"
tags: ["c", "graphics", "opengl", "sdl2"]
---
This post is largely inspired by [Chris Wellons'](https://nullprogram.com/) 2015 blog post about writing a [modern OpenGL demo](https://nullprogram.com/blog/2015/06/06/) that works on all three major desktop platforms (Windows, macOS, and Linux).
I have come back to his post countless times over the years when looking for guidance on how to build any cross-platform C program.

In Chris' version of this demo he makes use of [GLFW3](https://www.glfw.org/) for the window and input handling and uses [gl3w](https://github.com/skaslev/gl3w) for loading [OpenGL](https://www.opengl.org/) functions.
I have used GLFW3 in the past, but tend to prefer [SDL2](https://www.libsdl.org/) for its [wealth of features](https://wiki.libsdl.org/Introduction) and [polling-based event model](https://wiki.libsdl.org/SDL_PollEvent).
Therefore, in my version, I wanted to use SDL2 for the window and input handling and load the OpenGL functions myself.
Loading OpenGL functions is a fiddly, quirky topic that deserves its own [blog post](/posts/loading-opengl-fuctions-for-fun-and-profit/).
However, it doesn't require too much code once the nuances are understood.

# Modern OpenGL
The differences between "legacy" and "modern" OpenGL are well-documented across the web.
Instead of rewriting the wheel, I'll refer you to this [excellent write-up](https://glumpy.github.io/modern-gl.html) about it over on the docs for the GLUMPY project.
In short, old OpenGL used a "fixed pipeline" that was generally simpler to navigate but restricted programmer flexibility.
Modern OpenGL uses a "programmable pipeline" that gives the programmer much more control over what their GPU is doing but can slightly increase the complexity of their code.
The newer versions even allow developers to write the actual code that runs on the GPU (these are known as [shaders](https://www.khronos.org/opengl/wiki/Shader) and are written in a C-like language called [GLSL](https://www.khronos.org/opengl/wiki/OpenGL_Shading_Language)).

As a quick aside, the [Vulkan graphics library](https://www.khronos.org/vulkan/) is positioned even further in this same direction: more graphical control at the cost of added complexity.
This tradeoff is definitely worth it, though, given that the graphical needs of modern video games and other technologies are ever-increasing.

# The Demo
The demo itself is very minimal: a small window containing a rotating red square that exits upon hitting the `ESCAPE`, `Q`, or the close button in the corner.

![SDL2 OpenGL Demo](/static/img/sdl2-opengl-demo.png)

The source can be found here:
* https://github.com/theandrew168/sdl2-opengl-demo

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

The `src/opengl.h` and `src/opengl.c` files contain two things: OpenGL function declarations and the code to dynamically load them at runtime.
The function `opengl_load_functions()` should be called once after obtaining a valid OpenGL context.
Lastly, `src/main.c` implements the demo: create a window, initialize an OpenGL context, load the modern OpenGL functions, and render a rotating square.

# Building
Building cross-platform C applications isn't that hard but it does require some intentional, upfront planning.
I have to thank Chris and his writings again for this facet of the demo, too.
He has written multiple posts on how to keep C projects simple and portable.
I highly recommend the following posts to anyone interested in writing clean, simple, and portable C projects:
* [Four Ways to Compile C for Windows](https://nullprogram.com/blog/2016/06/13/)
* [How to Write Portable C Without Complicating Your Build](https://nullprogram.com/blog/2017/03/30/)
* [A Tutorial for Portable Makefiles](https://nullprogram.com/blog/2017/08/20/)

### Linux
Building the demo on Linux is a simple native build.
Only two dependencies are required: make and SDL2.
On a Debian-based Linux system, they can be installed via:
```
sudo apt install make libsdl2-dev
```

Once the dependencies are installed, the build itself is as simple as possible:
```
make
```

### macOS
Building the demo on macOS is a simple native build.
Only two dependencies are required: make and SDL2.
If using [brew](https://formulae.brew.sh/), they can be installed via:
```
brew install make sdl2
```

The build process for macOS is exactly the same as Linux aside from using a different Makefile.
```
make -f Makefile.macos
```

### Windows
Building the demo for Windows requires cross-compiling from either Linux or macOS.
This process relies on the amazing [mingw-w64](http://mingw-w64.org/doku.php) project which utilizes GCC tooling to build Windows executables and libraries.
Cross-compiling from Linux or macOS will require the following dependencies: make, wget, tar, and mingw-w64.

If using a Debian-based Linux system, these can be install via:
```
sudo apt install make wget tar mingw-w64
```

If using macOS with brew installed, the command is:
```
brew install make wget gnu-tar mingw-w64
```

With all dependencies installed and ready to go, the demo can be cross-compiled.
```
make -f Makefile.mingw
```

This will build the demo into a standalone, portable Windows executable: `demo.exe`.
This file can be tested with [Wine](https://www.winehq.org/) or transferred to any modern Windows system for native execution.

One thing you might notice is that the Windows build didn't require installing any SDL2 dependencies.
This is because the build actually downloads (using wget), extracts (using tar), and statically links the pre-built SDL2 libraries directly into the executable.
This is my favorite trick from Chris' original demo!
He used this approach for statically linking GLFW3 and I realized that it would work just as well for SDL2.

# Why SDL2?
I have used both GLFW3 and SDL2 in the past and struggled initially to decide which library I preferred.
Since then, though, I have settled on SDL2 as my de-facto choice for programs that require graphics and user input.
My reasons for preferring SDL2 are three-fold:
1. SDL2 comes with more batteries included than GLFW3
2. SDL2 supports more platforms than GLFW3
3. SDL2's polling-based input is cleaner to deal with than GLFW3's callbacks (GLFW3 _does_ support querying for current input states but that means potentially "missing" short-lived events)

Here is a personal list of "pros" for choosing SDL2:
* SDL2 is a one-stop-shop for cross-platform game development
* SDL2 supports Windows, Linux, macOS, iOS, Android, Nintendo Switch, and others
* SDL2 supports mouse and keyboard, joysticks, game controllers, and multi-touch gestures
* SDL2 supports low-level audio playback
* SDL2 includes extra helpers for dealing with threads, timers, CPU detection, and power management
* SDL2 includes a simple, accelerated 2D renderer for smaller projects
* SDL2 is released under the permissible [zlib license](https://opensource.org/licenses/Zlib)
* SDL2's reference docs are very detailed
* The SDL2 ecosystem includes "extension projects" that add specific functionality:
* [SDL\_image](https://www.libsdl.org/projects/SDL_image/) adds image loading
* [SDL\_mixer](https://www.libsdl.org/projects/SDL_mixer/) adds sound mixing
* [SDL_net](https://www.libsdl.org/projects/SDL_net/) adds networking
* [SDL_ttf](https://www.libsdl.org/projects/SDL_ttf/) adds TrueType font rendering

On the other hand, here are some "cons" related to SDL2:
* SDL2's officials docs are void of full examples and tutorials
* Third-party tutorials and examples are often outdated
* Third-party tutorials and examples often use legacy OpenGL

# Conclusion
There we have it!
Cross-platform graphical applications don't have to be fussy: they just require a bit of planning.
Even though I prefer SDL2, GLFW3 is a solid alternative when it comes to platform-independence library.
By loading our own OpenGL functions and keeping dependencies to a minimum, we can achieve maximum portability.
Even Windows can join in on the fun!
