---
date: 2020-06-08
title: "Loading OpenGL Functions for Fun and Profit"
slug: "loading-opengl-fuctions-for-fun-and-profit"
tags: ["c", "graphics", "opengl", "sdl2"]
draft: true
---
Introduction and why I'm writing this  
Came up during the SDL2 demo and figured it'd be worth explaining  
You can always use a library but you can do it yourself just for fun  
It's always nice to understand what problems a library is solving and how they do it  

# Functions vs Function Pointers
[cppreference details](https://en.cppreference.com/w/c/language/function_declaration)  
They are very similar (and even the same from a caller's perspective)  
A function is just an address in memory: CALL + RET is how they work  
A normal function is declared and defined at compile time  
A function pointer is more flexible: put in a struct, load dynamically, etc  
Do an example here with add?  

# Dynamic Loading
[Dynamic loading](https://en.wikipedia.org/wiki/Dynamic_loading#Special_library)   
Functions can be located and called at runtime!  
Modern opengl functions work like this: they aren't present at compile time  
But we at least need to have _some_ sort of compile-time handle to them for the rest of our code to call  
How do we do that? What do they pointers initially point to? What happens if one is called before being dynamically loaded?  

# The Gameplan
[Loading OpenGL without GLEW](https://apoorvaj.io/loading-opengl-without-glew/)  

| What do we need? | Where do we get it? |
| --- | --- |
| Func typedefs | SDL2/SDL_opengl.h, GL/glcorearb.h, etc |
| Func ptr declarations | We do this |
| Func ptr defaults | We do this |
| Func ptr loader | SDL_GL_GetProcAddress, OS-specific funcs |

# Making the compiler happy
Why do all of these libraries (GLAD, glLoadGen, gl3w, glew) get a warning with -Wpedantic?  
Gap in the C standard that results in obj ptrs and func ptrs not necessarily being the same size  
Therefore (void ptr -> func ptr) is _not_ a valid type conversion  
Most of the libs just ignore the warning and that is likely never going to cause any issues  
Most of the time, obj ptrs and func ptrs ARE the same size on a given platform  

[SO](https://stackoverflow.com/questions/13696918/c-cast-void-pointer-to-function-pointer)  
POSIX sort of mandates this behavior as being valid due to funcs such as dlsym()  
Also C11 has a common extension that legalizes this behavior  
Plus this is the only way to access OpenGL funcs! Can't do much about that  

How can we fix it, though? I like compiling clean with -Wpedantic  
Pull the pointer through a UNION to eliminate the potential ptr size difference  
Does it actually _solve_ the problem? Not really.  
If this ran on a platform where the sizes were different AND the func addr sat between the two sizes , you'd end up jumping to a nonsense address and likely segfaulting  

How can we write this code is the most precise way?  
What needs to be included in the union? void ptr and ALL GL funcs  
TODO Is there a void ptr equiv for functions? Not that I know of  
Talk about the three methods I came up with and the pros / cons  

# Cleaning things up with macros
Data-driven system of macros that simply the loading process: just add the func to the list  
What if you want _all_ the funcs? Yea, that'd be more work. Maybe use a library at that point  
Use fancy macro nesting to reuse the same list of required funcs in different ways:  
OPENGL_DECLARE, OPENGL_DEFINE, OPENGL_LOAD, and OPENGL_VALIDATE 
