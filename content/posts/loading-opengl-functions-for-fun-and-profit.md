---
date: 2020-06-17
title: "Loading OpenGL Functions for Fun and Profit"
slug: "loading-opengl-fuctions-for-fun-and-profit"
tags: ["c", "graphics", "opengl", "sdl2"]
---
While writing my [previous blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/), I realized that there is a general information gap around dynamically loading OpenGL functions.
When developers encounter this task on new OpenGL-related projects, they tend to reach for a premade loader such as [glew](http://glew.sourceforge.net/), [gl3w](https://github.com/skaslev/gl3w), [glLoadGen](https://bitbucket.org/alfonse/glloadgen/wiki/Home), or [GLAD](https://github.com/Dav1dde/glad).
A more detailed list along with an overview of what all these libraries do can be found [here](https://www.khronos.org/opengl/wiki/OpenGL_Loading_Library).

**"An OpenGL Loading Library is a library that loads pointers to OpenGL functions at runtime, core as well as extensions. This is required to access functions from OpenGL versions above 1.1 on most platforms. Extension loading libraries also abstracts away the difference between the loading mechanisms on different platforms."** - [Khronos Docs](https://www.khronos.org/opengl/wiki/OpenGL_Loading_Library)

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
Dynamic loading is an important facet of using OpenGL because the locations of most of the library's functions aren't known at compile time.
They also aren't necessarily known even at [dynamic link time](https://en.wikipedia.org/wiki/Dynamic_linker)!
Due to how OpenGL has evolved over the years (and to support more flexibility in handling different platforms and implementations), its modern functions often require dynamical loading before they can be used.

Okay, so now we know that the OpenGL functions that we want to use are going to have to be loaded sometime after the program is running.
How can we build an application around these functions if they aren't located until _after_ compiling and linking?
If you refer to the previous code snippet, we can initially define a function pointer to be `NULL` and then simply change what it points to at a later time.
That way, the rest of our application can use these initial definitions when being built.
We just have to be absolutely sure to point these definitions at something valid _before_ calling them.

# The Gameplan
With all of that background context out of the way, we can start to bring everything together.
I'll be using [SDL2](https://www.libsdl.org/) for this demo in order to keep the code consistent and portable.
Additionally, most the details on how to write something like this comes from [Apoorva Joshi's](https://apoorvaj.io/) excellent blog post: [Loading OpenGL without GLEW](https://apoorvaj.io/loading-opengl-without-glew/).
The official [Khronos Documentation](https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions) goes into great detail on the mechanics, as well.

| What do we need? | Where will we get it? |
| --- | --- |
| OpenGL function type signatures | The SDL2 header: [SDL2/SDL_opengl.h](https://github.com/spurious/SDL-mirror/blob/master/include/SDL_opengl.h) |
| OpenGL function declarations | We will write these ourselves! |
| OpenGL function definitions | We will write these ourselves! |
| OpenGL function address getter | The SDL2 function: [SDL_GL_GetProcAddress](https://wiki.libsdl.org/SDL_GL_GetProcAddress) |
| OpenGL function loader | We will write this ourselves! |

### OpenGL function type signatures
The [SDL2/SDL_opengl.h](https://github.com/spurious/SDL-mirror/blob/master/include/SDL_opengl.h) header (and the other version-specific headers that it pulls in) includes type signatures for all modern OpenGL functions.
In this scenario, a function's type signature refers to its name, return type, and argument types.
This information is enough to uniquely identify a given OpenGL function (or any C function, for that matter).
Since we plan to load these function implementations dynamically, we aren't looking for prototypes: we are looking for function pointer typedefs.

For example, here is the abridged prototype and function pointer typedef for function [glCreateShader](http://docs.gl/gl3/glCreateShader):
```
GLuint glCreateShader(GLenum type);
typedef GLuint (*PFNGLCREATESHADERPROC)(GLenum type);
```
About the odd typedef name `PFNGLCREATESHADERPROC`: it represents a **P**ointer to the **F**unctio**N** **glCreateShader**, which is a **PROC**edure.
This explanation cames from the [Khronos Docs](https://www.khronos.org/opengl/wiki/Load_OpenGL_Functions#Function_Prototypes).

### OpenGL function declarations
While the prototype is valuable for more clearly showcasing how [glCreateShader](http://docs.gl/gl3/glCreateShader) should be called, we don't actually want that prototype present in our codebase.
This is because of a nuance we covered earlier: a normal function's definition _must_ be present at compile/link time.
Since we won't have "real" definitions at compile/link time, we need to use the more flexible option: a function pointer.

The typedef seen above serves as a convenience.
It gives us an easy way to declare a function pointer with the same signature as the prototype.
Therefore, in our own header, we can declare all of the OpenGL functions that our project requires.
For the sake of brevity, I'll just be using two modern OpenGL functions in my examples: [glCreateShader](http://docs.gl/gl3/glCreateShader) and [glDeleteShader](http://docs.gl/gl3/glDeleteShader).
```
// opengl.h
#include <SDL2/SDL_opengl.h>

PFNGLCREATESHADERPROC glCreateShader;
PFNGLDELETESHADERPROC glDeleteShader;
```

### OpenGL function definitions
Similar to original example, the default definition for our function pointers will simply be `NULL`.
Once again, these functions **MUST NOT** be called until after the real locations have been loaded and validated.
```
// opengl.c
#include <SDL2/SDL_opengl.h>
#include "opengl.h"

PFNGLCREATESHADERPROC glCreateShader = NULL;
PFNGLDELETESHADERPROC glDeleteShader = NULL;
```

### OpenGL function address getter
In order to get our hands on the actual address of these functions, a special function is needed to look the OpenGL functions by name.
Since our demo uses [SDL2](https://www.libsdl.org/), we will be using the platform-agnostic helper that it provides: [SDL_GL_GetProcAddress](https://wiki.libsdl.org/SDL_GL_GetProcAddress).
```
// "proc" is the name of an OpenGL function; returns NULL upon error
void* SDL_GL_GetProcAddress(const char* proc)
```

Each OpenGL implementation platform has its own, specific flavor of this function:
* **Windows** has [wglGetProcAddress](https://docs.microsoft.com/en-us/windows/win32/api/wingdi/nf-wingdi-wglgetprocaddress)
* **Linux/X11** has [glXGetProcAddressARB](https://www.khronos.org/registry/OpenGL/extensions/ARB/GLX_ARB_get_proc_address.txt)
* **macOS** has [NSLookupAndBindSymbol](https://developer.apple.com/library/archive/documentation/GraphicsImaging/Conceptual/OpenGL-MacProgGuide/opengl_entrypts/opengl_entrypts.html)
  * Note that OpenGL was deprecated in macOS 10.14 in favor of [Metal](https://developer.apple.com/metal/)

### OpenGL function loader
The last step now is to call our address getter for each OpenGL function and update its initial definition with the real location.
We can wrap this process in a single helper function along with some simple error checking.
```

// opengl.h
bool opengl_load_functions(void);

// opengl.c
bool opengl_load_functions(void) {
    glCreateShader = (PFNGLCREATESHADERPROC)SDL_GL_GetProcAddress("glCreateShader");
    glDeleteShader = (PFNGLDELETESHADERPROC)SDL_GL_GetProcAddress("glDeleteShader");
    if (glCreateShader == NULL) return false;
    if (glDeleteShader == NULL) return false;
    return true;
}
```

There we go!
As long as `opengl_load_functions` is called _after_ acquiring an OpenGL context, [glCreateShader](http://docs.gl/gl3/glCreateShader) and [glDeleteShader](http://docs.gl/gl3/glDeleteShader) can now be called and used normally!
Even though we've only loaded two functions, the process is exactly the same for the rest.

However, something isn't quite right: the compiler isn't happy about this code.
Even though this is exactly what most OpenGL function loader libraries do, there is still a warning message lingering around in the compilation output.
It even only shows up when building with `-Wpedantic` enabled...
```
ISO C forbids conversion of object pointer to function pointer type
```

# The Problem
What we have encountered is a historic nuance in the C standard: objects pointers and function pointers are not convertible.
Also, it's not that the standard mandates that the two types _not_ be compatible, it just doesn't explicitly say that they are.
On some older platforms (and maybe modern ones, too), function and object pointers were indeed sized differently.

For us, this means that it is not technically legal to cast the `void*` pointer we get from [SDL_GL_GetProcAddress](https://wiki.libsdl.org/SDL_GL_GetProcAddress) to a function pointer.
In practice, however, this is very unlikely to ever cause issues.
The platform-specific function address getters listed earlier all imply that this conversion is valid (and required).
Otherwise, how would we ever call the modern OpenGL functions?
Some OpenGL loader libraries ignore this warning altogether and that is perfectly reasonable.

There exists a decent amount of prior conversation of this issue around the web.
Here are a couple of StackOverflow questions discussing [casting void pointers to function pointers](https://stackoverflow.com/questions/13696918/c-cast-void-pointer-to-function-pointer) and [what is guaranteed about function pointer sizes](https://stackoverflow.com/questions/3941793/what-is-guaranteed-about-the-size-of-a-function-pointer).
There have even been [bugs reported to GCC](https://gcc.gnu.org/bugzilla/show_bug.cgi?id=83584) that seek clarification on this warning and how seriously it should be taken.

The POSIX [dynamic linking API](https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/dlfcn.h.html) includes a more generic function address getter.
The docs for [dlsym](https://pubs.opengroup.org/onlinepubs/9699919799/functions/dlsym.html) make special note of this function-casting behavior:

**"Note that conversion from a void * pointer to a function pointer [...] is not defined by the ISO C standard. This standard requires this conversion to work correctly on conforming implementations."**

Lastly, the [C11 standard](http://port70.net/~nsz/c/c11/n1570.html) actually contains a [common extension](http://port70.net/~nsz/c/c11/n1570.html#J.5.7) that permits this behavior and makes it valid:

**"A pointer to an object or to void may be cast to a pointer to a function, allowing data to be invoked as a function"**

# Appeasing the Compiler
Even though this warning isn't something that we need to be very worried about, it's still there: cluttering up our compilation output.
How can we make the compiler happy and still compile with pedantic warnings enabled?
The issue essentially boils down to a potential difference in sizes: an object pointer _might_ be smaller or larger than a function pointer.
What feature does C give us that can bridge a size gap between different types?

A [union](https://en.cppreference.com/w/c/language/union)!
As a reminder, members of a union overlap and occupy the same space in memory (unlike a [struct](https://en.cppreference.com/w/c/language/struct) which stores its members sequentially).
If we create a union that holds both an object pointer and a function pointer, then it is guaranteed to be large enough to hold both types.
```
union bridge {
    void* object_ptr;
    void (*function_ptr)(void);
};
```

With this, we can effectively "push" the `void*` from [SDL_GL_GetProcAddress](https://wiki.libsdl.org/SDL_GL_GetProcAddress) into the union and "pull" a function pointer out.
Using a [C99 designated initializer](https://en.cppreference.com/w/c/language/struct_initialization), we can initialize the union's `object_ptr` and grab its `function_ptr` all in a single expression:
```
glCreateShader = (PFNGLCREATESHADERPROC)(union bridge){
    .object_ptr = SDL_GL_GetProcAddress("glCreateShader")
}.function_ptr;
```

The cast to `PFNGLCREATESHADERPROC` is required here because the `function_ptr` in our union is of a different function pointer type.
It doesn't matter that `function_ptr` is a different function pointer type than any OpenGL function: it just needs to be a function pointer so that its size is accounted for.

With this code now in place, the compiler is quiet again!
We can compile with all warnings enabled (all, extra, pedantic) and no fuss is raised.
The only thing left to do is to find a way to wrap all of these fiddly concepts together in a concise collection of macros.

# Macros to the Rescue
In order to keep things simple and extensible, our custom OpenGL function loader two goals:
1. It should work on most modern platforms without fuss
2. It should be possible to add and remove functions from the project by changing a single line of code

The good news is that this is possible!
The bad news is that it involves a small amount of [C preprocessor](https://en.cppreference.com/w/c/preprocessor) complexity.
The core strategy is to develop a set of macros that all accept the same two parameters: the OpenGL function name and its corresponding function pointer name.
Then, we will be able to define a list of all the functions we need and swap out which macro expands them at different locations throughout the code.
```
#define OPENGL_FUNCTIONS                                    \
    OPENGL_FUNCTION(glCreateShader, PFNGLCREATESHADERPROC)  \
    OPENGL_FUNCTION(glDeleteShader, PFNGLDELETESHADERPROC)  \
    ...
```

Keep this list in mind as we explore each of the four helper macros.

### OPENGL_DECLARE
This macro declares an OpenGL function pointer.
```
#define OPENGL_DECLARE(func_name, func_type)  \
    func_type func_name;
```

It is used in our project's OpenGL-specific header file to expose the symbols that other parts of our application will link to.
```
// opengl.h
#define OPENGL_FUNCTION OPENGL_DECLARE
OPENGL_FUNCTIONS
#undef OPENGL_FUNCTION
```

Notice how this code says "make `OPENGL_FUNCTION` expand to `OPENGL_DECLARE`" before dropping in the whole list of functions.
This way, as the preprocessor handles all of the listed functions from earlier, they will each expand into a function pointer declaration.
The `OPENGL_FUNCTION` macro is then reset back to being undefined just for peace of mind.
This is pattern is how _all_ of these macros interact with that simple, data-driven list of functions!

### OPENGL_DEFINE
This macro defines the initial function pointer implementations (which is `NULL`).

```
#define OPENGL_DEFINE(func_name, func_type)  \
    func_type func_name = NULL;
```

It is used in our project's OpenGL-specific implementation file.
```
// opengl.c
#define OPENGL_FUNCTION OPENGL_DEFINE
OPENGL_FUNCTIONS
#undef OPENGL_FUNCTION
```

### OPENGL_LOAD and OPENGL_VALIDATE
This first macro encompasses the union pass-through and cast decribed earlier.
```
#define OPENGL_LOAD(func_name, func_type)                \
    func_name = (func_type)(union bridge){               \
        .object_ptr = SDL_GL_GetProcAddress(#func_name)  \
    }.function_ptr;
```

The second one is technically optional but still highly recommended in some form or another.
It validates that an OpenGL function was successfully loaded and prints an error if something isn't right.
```
#define OPENGL_VALIDATE(func_name, func_type)                      \
    if (func_name == NULL) {                                       \
        fprintf(stderr, "failed to load func: %s\n", #func_name);  \
        return false;                                              \
    }
```

These two macros together form the entirety of the `opengl_load_functions` function.
```
bool opengl_load_functions(void) {
    #define OPENGL_FUNCTION OPENGL_LOAD
    OPENGL_FUNCTIONS
    #undef OPENGL_FUNCTION

    #define OPENGL_FUNCTION OPENGL_VALIDATE
    OPENGL_FUNCTIONS
    #undef OPENGL_FUNCTION

    return true;
}
```

Not _too_ complicated once all is said and done!

# Conclusion
This post has discussed quite a few code snippets in isolation.
To bring it all together, [here is a gist](https://gist.github.com/theandrew168/2eec79a145396c5d08b774096f91c922) containing a complete header and implementation for the loader we’ve built.
You can also see a full version with extra comments in my [sdl2-opengl-demo](https://github.com/theandrew168/sdl2-opengl-demo) project.

Despite knowing more about how this process works, should you still do it yourself?
Is it worth reinventing the wheel or would an established function loader library be a better and safer option?
That's a call that you'll have to make yourself given the specific circumstances of your project and what you truly want to learn.
Using a library here has legitimate benefits: they cover more platform-specific edge cases, they have many more eyes on the code, and rarely require any manual changes or updates.

Reglardess of the choice you make, hopefully this post has shed a bit of light on what goes on behind the scenes of an OpenGL function loader.
