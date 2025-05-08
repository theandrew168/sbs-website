---
date: 2024-01-07
title: "Why Write?"
slug: "why-write"
tags: ["WebGL"]
---

How long has it been since my last post?
Let's see... about a week shy of 3 years!
Blogging is one of those "things I'd like to do more of" but never make time for.
I've still been writing code since my last post but I haven't written about any of it.
I ran my own software company for a while, got a new job at a Silicon Valley startup, and even got married!
There are so many moments of progress and happiness that only exist in my memories (and photos, of course).
I'd like to fix this.

While not quite a proper "New Year's Resolution", my buddy [Nick](https://nickherrig.com) and I have challenged ourselves to write more.
Ideally, one post per week for all of 2024.
Have I learned at least 52 new things in the past 3 years? Definitely!
I'll start with something short and sweet.

## WebGL Rocks!

I have always had a lingering interest in graphics programming.
Back in the day, researching how to write graphics code led me to [OpenGL](https://www.opengl.org/).
This is a native graphics API that is written in C (with bindings written for many other languages).
Despite enjoying C programming and wanting to embrace OpenGL's native tooling, my small applications / games had a big program: distribution.

Trying to share a native program with friends is difficult.
You have two options: send them the compiled binary (for their specific operating system and architecture) or show them how to install a compiler and build it themselves.
The second is a poor choice due to complexity and the first is made difficult due to operating system security measures.
It turns out that Windows and macOS really dislike running executables that they don't trust... and rightfully so.
If _only_ there was a way to share my demos with folks without them needing to build, download, or install anything.

Enter [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API).
MDN describes it as "a JavaScript API for rendering high-performance interactive 3D and 2D graphics within any compatible web browser without the use of plug-ins".
This seems great!
I can write my code in [TypeScript](https://www.typescriptlang.org/), build it and throw the output onto a web server, and simply send the link to anyone.
Once deployed (with a valid domain and TLS cert), my app is instantly available to anyone that owns a computer and / or smart phone.

While the upper bounds of performance are certainly lower when using a browser-based, JavaScript-based graphics programming environment, I don't think that it impacts me very much.
I don't plan on creating anything incredibly large or high fidelity: just basic games and examples.
If completely solving the problem of distribution means sacrificing a bit of performance, then I'm satisfied.
