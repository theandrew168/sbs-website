---
date: 2020-05-17
title: "Refactoring My Old CHIP-8 Emulator"
slug: "refactoring-my-old-chip8-emulator"
tags: ["c", "sdl2", "emulator"]
draft: true
---
Back in 2017, I was really interested in emulator development.
I read that [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) was a great introductory system and decided to write my own emulator for it.
I was a C programming novice at the time but figured it'd still be a good choice for a project that dealt with a lot of low-level details and mechanics.
After a few weeks of work I was able to successfully emulate CHIP-8 games but knew that there were some lingering bugs.
Once I got the project to some semblance of "done", I moved on and never touched the codebase again.
You can find the final old version at [this commit](https://github.com/theandrew168/skylark/tree/a24585b48de2923fd016f379c7b0ad8cbb0a9d75).

In subsequent years I continued to write more C and hone my overall programming skills.
Recently, I walked through this project again and my interim years of progress became clear.
Almost every aspect of the code and how it was structured was outdated in comparison to my modern habits.
I knew that a good amount of tender loving care was going to be needed in order to get the project up to my current standards.

# Lingering Problems
The issues I found upon reviewing this code can be grouped into the following categories:
* **Naive Makefile** - The Makefile I originally wrote was simple and effective.
However, it wasn't [POSIX-compliant](https://pubs.opengroup.org/onlinepubs/009695399/utilities/make.html) and relied on extensions that are specific to [GNU Make](https://www.gnu.org/software/make/).
This lack of compliance meant that the build was unlikely to work on non-GNU systems such as [macOS](https://en.wikipedia.org/wiki/MacOS) or [BSD variants](https://en.wikipedia.org/wiki/Comparison_of_BSD_operating_systems).
In addition to the lack of compatibility, the original Makefile only builds the target binary and nothing else.
My modern Makefile template builds not only the target binary but also a test binary and both static and shared libraries.
* **Tight Coupling** - The original codebase was structured as three tightly-coupled sections: `chip8.c`, `graphics.c` and `input.c`.

* **Global State** - The original sections all made heavy use of global state.
* **Lack of Tests** - The project had no tests at all!
* **Poor Platform Support** - The project only worked (and was ever tested) on Linux.

# Mighty Makefile
Standard template  
plus extra tool for baking ROMs  
means that the dist is just a single file  
no program files  

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
