---
date: 2024-08-04
title: "Programming on Vacation"
slug: "programming-on-vacation"
draft: true
---

I'm going to be on vacation this week!
My wife and I (and some friends) are headed down to the Lake of the Ozarks in Missouri to swim, boat around, and chill.
That being said, you all know me: I enjoy programming and vacations give me time to "switch gears" and adapt a different development schedule.
See, while we are here, people are always around: chatting, eating, or just having a good time.
I love that aspect of vacation but it makes deep work pretty difficult.
For me, deep work requires isolation and quietness in order to be successful.
We also change our plans on the fly which makes committing to long blocks of uninterruped time difficult.
But again, that's fine!
This is vacation after all.
Relaxing and socializing are my priorities down here.

However, there are still _some_ tasks I can tackle.
In fact, I think that this environment might be the best fit for certain genres of tasks: those that seem too simple / tedious / mundane to be worth wasting valuable focus time.
For example, this is a great opportunity to explore new testing strategies (like Go's benchmarks).
I added a simple [BenchmarkParse](https://github.com/theandrew168/bloggulus/commit/6495db10f9f5fb98b1cf91683f0162aa5127e346) test to [Bloggulus](https://github.com/theandrew168/bloggulus) to get a sense of how long parsing a simple RSS feed takes.
Granted, I don't ever expect that to be a bottleneck, but I like the experience of seeing what it takes to write that sort of performance-focused test.

Another good vacation task is documentation!
One of my other side projects is a very early-days, bare-bones game idea called [bugforth](https://github.com/theandrew168/bugforth).
It is written for the browser with WebGL graphics and a Svelte UI.
I'd been hacking on the code for so long that I figured it was time to explain how things worked.
I think that the code was kind of difficult to reason about: especially the WebGL / batch renderer aspects.
So, this trip is a great time to go through all of the code and [document](https://github.com/theandrew168/bugforth/commit/877f1dce897cd1b3c309bbfd75142ce448eb6bff) all classes, methods, and functions!
Adding docs requires little focus but adds a ton of value to any project.

Thanks for reading!
