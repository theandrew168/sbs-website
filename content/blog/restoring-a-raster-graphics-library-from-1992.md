---
title: "Restoring a Raster Graphics Library From 1992"
date: 2019-08-24
tags: []
draft: true
---
### Intro
I've recently been reading _Computer Graphics: Principles and Practice (2nd Edition)_.
It is a foundational graphics book first published in 1982.
Some of the content is dated but most of the ideas presented are fundamental to computer graphics theory.

The book uses two libraries to demonstrate principles of both 2D and 3D graphics: SRGP and SPHIGS, respectively.
From what I've seen, both of these libraries suffer from considerable amounts of [dormant code rot](https://en.wikipedia.org/wiki/Software_rot).
I really wanted to get them working again so that I could follow along with the book's examples.
In this post, I detail the process of evaluating the status of SRGP (Simple Raster Graphics Package), cleaning up the code rot, and getting it running again on a modern X11 desktop.

### Fixing the build
Testing on Ubuntu 18.04  
apt install libx11-dev  
mkdir objects  
works!  
building an example:  
`gcc test_keyboard.c -I ../src/srgp/ -L ../src/srgp/objects/ -lsrgp -lX11`  
run the example... segfault!  

### Following along with GDB
existing codebase makes it hard to follow the exact control flow  
use gdb to follow the critical path of X11 calls  
talk about xfollow.c, link to tronche tutorial  

### Enabling custom colors
TODO: Currently X11 is grabbing a TrueColor visual, I need DirectColor (I think)  

### Mission accomplished!
At this point the library "works".  
but it could def be cleaner and simpler  

### Refactoring the build process
Remove MAC files from project  
Make the X versions default  
Revamp Makefiles  

### Cleaning up the code
Simplify the code where possible  

### Cleaning up the docs
Rewrite manual in Markdown?  
