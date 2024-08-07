---
date: 2024-08-04
title: "Programming on Vacation"
slug: "programming-on-vacation"
---

I'm going to be on vacation this week!
My wife and I (and some friends) are headed down to the Lake of the Ozarks to swim, boat around, and chill.
That being said, I genuinely enjoy programming and vacations give me time to "switch gears" and adapt a different development schedule.
See, while we are here, people are always around: chatting and having a good time.
We also tend to change plans on the fly which makes committing to long blocks of uninterruped time difficult.
I love these aspects of vacation but they make deep work somewhat difficult.

For me, deep work requires isolation and quietness in order to be successful.
But again, that's fine!
This is vacation after all.
Relaxing and socializing are my priorities down here.
However, there are still _some_ tasks I can tackle.
In fact, I think that this environment might be the **best** fit for certain genres of tasks: those that seem too simple / tedious / mundane to be worth wasting valuable focus time at home.

For example, this is a great opportunity to explore new testing strategies (like Go's [benchmark](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go) feature).
I added a simple [BenchmarkParse](https://github.com/theandrew168/bloggulus/commit/6495db10f9f5fb98b1cf91683f0162aa5127e346) test to [Bloggulus](https://github.com/theandrew168/bloggulus) to get a sense of how long parsing a simple RSS feed takes.
Granted, I don't ever expect that to be one of the application's bottlenecks, but I like the experience of seeing what it takes to write that sort of performance-oriented test.

Another good vacation task is documentation!
One of my other side projects is a very early-days, bare-bones game idea called [bugforth](https://github.com/theandrew168/bugforth).
It is being written for the browser with [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API) graphics and a [Svelte](https://svelte.dev/) UI.
I've been hacking on the code long enough that I figured it was time to explain how things work (especially the WebGL / batch renderer aspects).
So, this trip is a great time to go through all of the code and [document](https://github.com/theandrew168/bugforth/commit/877f1dce897cd1b3c309bbfd75142ce448eb6bff) the project's classes, methods, and functions!
Adding docs requires little focus but adds a ton of value to any project.

Even on vacation, it is possible to find small but valuable tasks to get done.
I like to save the deep work for home and make the most of my time off.
Thanks for reading!
