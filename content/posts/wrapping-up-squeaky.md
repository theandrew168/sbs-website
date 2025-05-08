---
date: 2020-08-25
title: "Wrapping Up Squeaky"
slug: "wrapping-up-squeaky"
tags: ["C", "SDL2", "Scheme"]
---

My primary project for the last month or so has a been a proof of concept programming language called [Squeaky](https://github.com/theandrew168/squeaky).
I've had this idea in my head for a few years: build a functional language in C with a focus on game development.
I wanted window creation, OpenGL graphics, and input events to all be first class citizens of the language and not optional libraries to be pulled in later.
Cross-platform portability was important, too.
If you write a game that utilizes the language's builtin capabilities for graphics and keyboard input, it should work the same on all major platforms: Windows, macOS, and Linux.

The language currently shows that all of this is possible.
It may only be able to draw white lines on a black backgroud, but that proves that arbitrary graphics are possible.
It may only be able to handle left and right arrow keys, but that proves that arbitrary input is possible.
Given these facts, I am able to call the project a success even though it didn't meet my initial expectations.
There are some more advanced topics that I want to explore in smaller subprojects before returning to Squeaky.

## Bumps in the Road

One of my biggest gripes with the language is that you need the interpreter to be present to run any code.
This is also how [Python](https://www.python.org/) works but Python's popularity and ecosystem offset that cost.
What I'd really like is a language implementation that generates static binaries like [Go](https://golang.org/) does.
Thinking about this led me to a realization: if my language simply takes text as input and emits bytes as output, then the implementation language doesn't actually matter.
I could write the same language in Python just as much as I could write it in C.
But what about SDL2?!

The problem of linking in SDL2 functions is an interesting one that I plan to explore.
As long as I have the static libraries present during compilation, what prevents me from looking up the symbol and injecting its object code into the target binary myself?
This is super hand-wavey and complicated but I want to try it out.
I envision a Scheme compiler written in a friendly language (such as Python) that uses prebuilt static libraries as a sort of toolbox to be utilized during compilation.
A separate library for each platform can be held to enable cross-platform code generation.
If the program source wants to create a window, the compiler will reach into the target platform's `libSDL2.a` static library, find the `SDL_CreateWindow` symbol, and inline its code into the target binary.

## Future Plans

I don't actually know if this is feasible or not but I'm excited to find out!
If it does work, you can expect a more detailed write-up on the nuances and trade-offs for this approach to compiling foreign functions.
In the meantime, I have about 50 years of assembly basics to catch up on.
Wish me luck!
