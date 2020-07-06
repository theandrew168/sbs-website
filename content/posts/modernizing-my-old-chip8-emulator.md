---
date: 2020-07-05
title: "Modernizing My Old CHIP-8 Emulator"
slug: "modernizing-my-old-chip8-emulator"
tags: ["c", "sdl2", "emulator"]
draft: true
---
Back in 2017 I was really interested in emulator development.
I read that [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) was a great introductory system and decided to write my own emulator for it.
I was a C programming novice at the time but figured it'd still be a good choice for a project that dealt with a lot of low-level details and mechanics.
After a few weeks of work I was able to successfully emulate CHIP-8 games but knew that there were some lingering bugs.
Once I cleaned them up and got the finished the project, I moved on and never touched the codebase again.
You can find this old version at [this commit](https://github.com/theandrew168/skylark/tree/a24585b48de2923fd016f379c7b0ad8cbb0a9d75).

In subsequent years I continued to write more C and hone my overall programming skills.
Recently, I walked through this project again and my interim years of experience became clear.
Almost every aspect of the code and how it was structured was outdated relative to my modern habits.
I knew that a good amount of tender loving care was going to be needed in order to get the project up to my current standards.

# Lingering Problems
The issues I found upon reviewing this code can be grouped into the following categories:
* **Naive Makefile** - The Makefile I originally wrote was simple and effective.
However, it wasn't [POSIX-compliant](https://pubs.opengroup.org/onlinepubs/009695399/utilities/make.html) and relied on extensions that are specific to [GNU Make](https://www.gnu.org/software/make/).
This lack of compliance meant that the build was unlikely to work on non-GNU systems such as [macOS](https://en.wikipedia.org/wiki/MacOS) or [BSD variants](https://en.wikipedia.org/wiki/Comparison_of_BSD_operating_systems).
In addition to the lack of compatibility, the original Makefile only builds the target binary and nothing else.
I'd like to to also build a test binary, a static library, and a shared library.
* **Tight Coupling** - The original codebase was structured as three tightly-coupled sections: `input.c`, `graphics.c`, and `chip8.c`.
These sections corresponded to keyboard-based input handling, graphical output, and then everything else.
There are definitely better ways to break the pieces of an emulator apart than this.
Having such large, over-scoped chunks made it difficult to look at specific concepts in isolation.
I'd prefer to employ a project structure that more clearly defines the individual ideas of an emulator.
* **Global State** - Each of the three sections mentioned above had global state tied to their translation units.
This meant that almost none of the project's functions were [pure](https://en.wikipedia.org/wiki/Pure_function) or even [thread-safe](https://en.wikipedia.org/wiki/Thread_safety).
Even though this emulator doesn't use threads, I don't like having the possibility blocked from the start.
These quirks also made the sections very difficult to test because every function call would change some global state that affected subsequent calls.
* **Lack of Tests** - This problem is fairly self-explanatory.
The project had no tests at all!
This was due to a couple reasons.
First, I didn't know how to write tests into a C project.
I didn't understand Makefiles or libraries enough to add test execution into the build.
Second, the global state mentioned previously made it hard to decide _what_ to test.
The whole project was this big tangle of code which made any test look daunting.
* **Poor Platform Support** - Another facet of my inexperience with C and building C projects was the absense of portability.
The project was only ever built and ran on Linux.
I've since learned an effective approach for writing cross-platform C applications and would like to see this emulator working on Linux, macOS, and Windows.

# Mighty Makefile
The [original Makefile](https://github.com/theandrew168/skylark/blob/a24585b48de2923fd016f379c7b0ad8cbb0a9d75/Makefile) had a fairly linear process: use GNU Make extensions to list all of the source files, use a suffix rule to compile C source files into object files, and then link the target binary.
Simplified, it looks like this:
```
# glob all of the source files and name object files
SRCS = $(wildcard src/*.c)
HEAD = $(wildcard src/*.h)
OBJS = $(patsubst %.c, %.o, $(SRCS))

# the main executable depends on all object files
skylark: $(OBJS)
    $(CC) -o $@ $^ $(LDFLAGS)

# inference rule to compile C source to object files
%.o: %.c $(HEAD)
    $(CC) $(CFLAGS) -c $< -o $@
```

This small excerpt alone contains three non-POSIX, GNU Make extentions: `wildcard`, `patsubst`, and `$^`.
Additionally, this Makefile is hard-coded to only build a single target: the main binary.
There is little flexbility for customizing the build in its current form.
What I would prefer is a more free-form, ad hoc mapping between source files and build targets: executable binaries, static libraries, and shared libraries.

### Build Goals
More specifically, I'd like to be able to build the following targets:

| Target | Sources | Description |
| --- | --- | --- |
| libskylark.a | `src/chip8.c`, `src/isa.c` | a static library |
| libskylark.so | `src/chip8.c`, `src/isa.c` | a shared library |
| skylark | `src/main.c` | the CHIP-8 emulator |
| skylark_tests | `src/main_test.c` | an executable that runs the project's tests | 
| dis | `tools/dis.c` | a minimal CHIP-8 disassembler |
| rom2c | `tools/rom2c.c` | a tool for converting CHIP-8 ROMs to C source |

Fortunately, this situation is exactly what Make was built to solve.
This table of build targets can be easily expressed via Make's simple system of targets, rules, and dependencies.

### POSIX Preamble
Our Makefile starts out with a few lines of behavior specification:
```
.POSIX:
.SUFFIXES:
```

The first line `.POSIX:` tells Make to strictly adhere to POSIX functionality.
This is an important safeguard against unknowingly using GNU-specific Make extensions.

The second line `.SUFFIXES:` disables all default inference rules.
In order to keep every part of the build explicit, we don't want any implicit inference rules building things for us without our intention.
Individual suffixes will be added back in when and where they are necessary.

### Variables
I like to keep all of the configurable build variables right at the top of the Makefile.
These encompass details such as which C compiler to use, what extra libraries to link with, and various other settings and flags.
```
AR      = ar
CC      = cc
CFLAGS  = -std=c99
CFLAGS += -fPIC -g -Og
CFLAGS += -Wall -Wextra -Wpedantic
CFLAGS += -Wno-unused-parameter
CFLAGS += -Isrc/ -I/usr/include/SDL2
LDFLAGS =
LDLIBS  = -lSDL2
```

These variables can all be easily overriden on the command line.
For example, to build the project using `clang` instead of your platform's default compiler, simply override the `CC` variable: `make CC=clang`.

### Targets
The next sections defines all of the build artifacts and which ones should be built by default.
This can be adjusted based on personal preference.
In my case, I only build the main emulator and the tests executable by default.
To build everything, you can use `make all`.
```
default: skylark skylark_tests
all: libskylark.a libskylark.so skylark skylark_tests dis rom2c
```

With all that bookkeeping out of the way, we can start compiling things.

### Libraries
Let's start off by declaring the source files that should be built into the `libskylark` libraries and their corresponding object file names (just swapping the `.c` for `.o` in this case).
```
libskylark_sources = src/chip8.c src/isa.c
libskylark_objects = $(libskylark_sources:.c=.o)
```

Now we can express dependencies between object and source / header files.
These lists are what allow Make to rebuild _only_ the parts of the project that are affected by a given file change.
```
src/chip8.o: src/chip8.c src/chip8.h src/isa.h
src/isa.o: src/isa.c src/isa.h
```

We can then specify the targets for each library.
They both depend on the compiled object files.
```
libskylark.a: $(libskylark_objects)
    $(AR) rcs $@ $(libskylark_objects)

libskylark.so: $(libskylark_objects)
    $(CC) $(LDFLAGS) -shared -o $@ $(libskylark_objects) $(LDLIBS)
```

The last step is to add is an inference rule for compiling C source files into object files.
Inference rules apply their actions to files with suffixes that match the rule.
They are handy for one-to-one steps like we have here (one object file for one C source file).
They are _not_ as useful for many-to-one steps such as linking executables and libraries.
```
.SUFFIXES: .c .o
.c.o:
    $(CC) $(CFLAGS) -c -o $@ $<
```

### Main Executable
With both a static and shared library already built, executables became trivial to build.
All we need is it to specify a C source file that contains a `main` function, compile it, and then link it with one of the libraries.
I'm linking with the static library here for simplicity.
I would only ever really use the shared library for sharing skylark's functionality with other developers and not for building executables locally.
```
skylark: src/main.c libskylark.a
    $(CC) $(CFLAGS) $(LDFLAGS) -o $@ src/main.c libskylark.a $(LDLIBS)
```

### Tests Executable
Building the tests executable is similar to the main executable except that it has more than just a single `main` file.
To keep the project modular, the tests for each emulator section are kept separate in their own files.
```
skylark_tests_sources = src/chip8_test.c src/isa_test.c
skylark_tests: $(skylark_tests_sources) src/main_test.c libskylark.a
    $(CC) $(CFLAGS) -o $@ src/main_test.c libskylark.a
```

### Loose Ends
With all of the important artifacts out of the way, all that remain are the two "helper" programs: `dis` and `rom2c`.
These, similar to the main executable, are just single-file programs that can be linked with a skylark library if necessary.
```
dis: tools/dis.c libskylark.a
    $(CC) $(CFLAGS) -o $@ tools/dis.c libskylark.a

rom2c: tools/rom2c.c
    $(CC) $(CFLAGS) -o $@ tools/rom2c.c
```

### Summary
That wraps the Makefile revamp!
It is definitely a lot more flexible than the old version and does a lot more work.
I'm much happier with the way it is written now.
You can find the full, current version [here](https://github.com/theandrew168/skylark/blob/master/Makefile).

# Functional Foundations
talk about data flow  
link to hoist your IO talk  
talk about how that applies here  
graphics and keyboard is the IO, all else can be pure  

# Tactical Testing
uses minunit with a small edit  
usual testing structure  
test the pure functional stuff  

# Plentiful Platforms
native for unix-like via posix-compat make  
cross-compile for windows via minwgw-w64

# Above and Beyond
baking ROMs into the binary  
game selection screen  
