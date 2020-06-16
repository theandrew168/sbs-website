---
date: 2020-06-16
title: "Loading OpenGL Functions for Fun and Profit"
slug: "loading-opengl-fuctions-for-fun-and-profit"
tags: ["c", "graphics", "opengl", "sdl2"]
draft: true
---
While writing my [previous blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/), I realized that there is a general information gap around dynamically loading OpenGL functions.
When developers encounter this task on new OpenGL-related projects, they tend to reach for a premade loader such as [glew](http://glew.sourceforge.net/), [gl3w](https://github.com/skaslev/gl3w), [glLoadGen](https://bitbucket.org/alfonse/glloadgen/wiki/Home), or [GLAD](https://github.com/Dav1dde/glad).
A more detailed list along with an overview of what all these libraries do can be found [here](https://www.khronos.org/opengl/wiki/OpenGL_Loading_Library).

These libraries all do a great job of solving this problem and I've used them many times!
If the goal is to simply get OpenGL working successfully and move on, then definitely go with a pre-built loader.
However, I believe that any dependency added to a project should be understood at a high level (at least).
That way you can have a better sense of what the library is doing and how it saves you time and effort.
Plus, if some aspect of the library starts working unexpectedly, having the ability to "open the black box" and troubleshoot the problem is incredibly valuable.

The purpose of this post is to present and explain what it takes to load OpenGL functions without using of the aforementioned libraries.
The process can be fussy because there are a number of different OpenGL-isms and C-isms that must come together in order to keep the compiler happy.
However, with some careful reading and planning, a minimal loader can be written in a clean, simple, and data-driven fashion.
When all is said and done, adding / removing specific OpenGL functions will only require changing a single line of code!

# Functions vs Function Pointers
An important distinction to understand when tackling this topic is the difference between C functions and C function pointers.
From a caller's perspective they behave exactly the same: pass the args in parentheses and get a return value back.
Both representations of a function delineate the same thing: an address in memory containing instructions that implement a procedure.
A function is typically initiated by a CALL instruction and finalized by a RET instruction.
This is a _very_ simlified explanation but the point remains: normal functions and function pointers are different representations of the same concept.

Consider these two declarations for a function that adds two integers together:
```
int add_normal(int a, int b);
int (*add_pointer)(int a, int b);
```

The first one declares `add_normal` as a function that accepts two integers and returns an integer.
The second declares `add_pointer` as a function pointer (a pointer to an address in memory that implements the procedure) that _also_ accepts two integers and returns an integer.

When calling these two functions, the syntax is exactly the same:
```
int foo = add_normal(5, 7);
int bar = add_pointer(5, 7);
```

So, if these two means for declaring and calling functions are so similar, why would we ever use one over the other?
Especially since the normal function declaration is so much simpler: why even bother with function pointers?
Aside from being able to be passed to other functions, function pointers provide extra flexibility in terms of how they get defined.
A normal function's definition (the part of the code that implements its behavior) must be present when the program is compiled (or linked, to be more specific).
Additionally, the memory address of a normal function cannot be altered once the program begins execution.

On the other hand, a function pointer must still be defined at compile time but the address that it points to _can_ be changed at runtime.
Consider how we could use this concept to dynamically set the body of `add_pointer` at runtime:
```
int add_normal(int a, int b) {
    return a + b;
}

int (*add_pointer)(int a, int b) = NULL;

void setup_add_pointer(void) {
    add_pointer = add_normal;
}
```

This snippet initially defines `add_pointer` as `NULL`.
If called in this state, the program would quickly hit a [segfault](https://en.wikipedia.org/wiki/Segmentation_fault).
However, the function `setup_add_pointer` can be called to point `add_pointer` to something valid.
In this case, we point it to the same function body as `add_normal`.
This mechanism of changing a function pointer's target at runtime is known as [dynamic loading](https://en.wikipedia.org/wiki/Dynamic_loading).

# Dynamic Loading
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
