---
date: 2020-06-07
title: "Loading OpenGL Functions for Fun and Profit"
slug: "loading-opengl-fuctions-for-fun-and-profit"
tags: ["c", "graphics", "opengl", "sdl2"]
draft: true
---

# WIP Loading OpenGL Functions with SDL2
[Dynamic loading](https://en.wikipedia.org/wiki/Dynamic_loading#Special_library)   
[cppreference details](https://en.cppreference.com/w/c/language/function_declaration)  
Make this its own post if it goes deep enough?

Info needed:
* Func typedefs (SDL2/SDL_opengl.h, or others like glcorearb.h)
* Func ptr declarations (we do this)
* Func ptr impls (we do this)
* Func ptr loader (SDL_GL_GetProcAddress, or others)
* Process to link things up (we do this)
* Macros to clean things up (we do this)

All 3 of the "we do this" things can be simplified with simple macros.
Link to apoorva's post.
Talk about that the union thing?
Talk about why the func protos are so weird. Why not just declare it normal?
- Can't because it isn't known at link time.
- What about extern?
- Doesn't solve anything: still would have to be linkable at link time
