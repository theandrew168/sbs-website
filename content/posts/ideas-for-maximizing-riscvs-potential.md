---
date: 2020-11-11
title: "Ideas for Maximizing RISC-V's Potential"
slug: "ideas-for-maximizing-riscvs-potential"
tags: ["RISC-V", "Assembly", "Forth"]
---

There has been a lot of buzz surrounding the slowness, instability, and complexity of modern software systems.
It seems as though more and more people are feeling and observing bloat and bugginess in their day-to-day usage of computers.
In many situations, the mediocrity of today's technology has become so common that it isn't even seen as an issue.

During my few years spent in industry, these sort of quotes were heard almost daily:

- _Looks like Windows crashed. What a surprise!_
- _Shoot, my PC decided to do updates... see you in an hour._
- _This meeting room isn't working. Let's try to schedule a different one._
- _Just keeping refreshing the page until it renders correctly._

Businesses and individuals lose countless hours to these types of problems every day but they are usually brushed off with a laugh.
The stability of Windows is the source of many jokes.
Backup plans are necessary for inconsistent meeting room reservations.
Are these situations an inevitability of computer interactions or has something gone wrong?

## Existing Opinions

In his blog post [Software Disenchantment](https://tonsky.me/blog/disenchantment/), Nikita Prokopov (tonsky) presents evidence for "industry's lack of care for efficiency, simplicity, and excellence".
In his presentation [The 30 Million Line Problem](https://caseymuratori.com/blog_0031), Casey Muratori argues that the interaction between software and hardware has the potential to be simplified.
In his article [The Website Obesity Crisis](https://idlewords.com/talks/website_obesity.htm), Maciej Ceglowski describes how websites have gotten progressively larger and slower.
In his talk [Preventing the Collapse of Civilization](https://www.youtube.com/watch?v=pW-SOdj4Kkk), Jonathan Blow compares modern software development to historical technological advancements.
Lastly, in his 1995 essay [A Plea for Lean Software](https://cr.yp.to/bib/1995/wirth.pdf), Niklaus Wirth describes causes for what he calls "fat software".
1995!
This problem has a history.

## My Opinion

I agree with the sentiments expressed in each of these posts.
In some sense, it leaves me with a feeling of helplessness: what can we do about this?
Is the ship of modern software development too large to turn?
What is the real problem?
Even if we knew the answer to that, how could we fix our technological trajectory?
I don't want to accept that the tools we use everyday are incapable of being fast and reliable.

I've spent a lot of time thinking about this problem.
Not only what the root cause may be, but also what can be done to mitigate it.
The videos and articles mentioned prior have all helped guide me toward the following ideas.
In short:

1. Operating systems provide great value, but often at a cost
2. Modern software exhibits an abundance of code and complexity
3. Usefulness can present itself in different, non-traditional ways
4. There are other routes to actualizing useful behavior from computers
5. The Bronzebeard project is my attempt to enable one alternate route

## Idea 1: Where Things Went Wrong

Are modern operating systems the epitome of software efficiency?
By this, I mean to ask if they bring enough value to the table to justify their costs.
But Linux is free!
True, but I'm not referring to monetary costs: I'm referring to complexity and the sheer volume of code between you and chunk of silicon at the heart of your machine.

To be fair, I'm not really in a good position to answer this question.
I have zero kernel development experience at all.
However, there do exist veteran programmers who take issue with the state of modern kernel development.
[René Rebe](https://rene.rebe.de/) supports an ongoing conversation about this exact topic.
Through his [YouTube channel](https://www.youtube.com/user/renerebe), René presents frequent ideas and criticisms relating to operating system design and implementation.
He explains how better languages and better designs (such as [microkernel](https://en.wikipedia.org/wiki/Microkernel)) could be an important step in future OS innovation.

## Idea 2: Too Much Code, Too Much Complexity

Most of the time, code quantity is really a poor measure of anything.
Talking about a project's "lines of code" doesn't usually amount to much.
However, as a project's size increases, so does it's learning curve.
In my opinion, the easiest code to understand and maintain is code that doesn't exist.
A large amount of recent language debate has pertained to matters of safety.
Many agree that C is an "unsafe" language while Go and Rust are "safer" alternatives.
While I do concur that safety is important, I think that _quantity_ may be more significant.

I think at some scale, the safety of a language loses impact.
Thirty million line of code in _any_ language is going to result in subtle, nuanced, cross-cutting bugs.
As soon as a project grows to the point where no single developer can understand the whole picture, I think it's game over.
At such a scale, achieving optimal performance and reliability becomes nearly impossible because human intercommunication is lossy and inefficient.
I think that innovation will come not from writing _safer_ code, but from writing _less_ code.

## Idea 3: What Makes Something Useful?

Is a fully-featured, multi-user OS necessary for a computer to be useful?
Do specific, targeted applications such as a webserver need that?
I propose that, in fact, something like a webserver is made _less_ secure due to the multi-user model and features running behind it.
Perhaps even the performance of a dedicated server could be better without that extra overhead.

When it comes to "value", there are two equally-important angles to consider.
The first and most common is the value that comes from features: more code added in order to do more things.
The second and less common angle is the value added _implicity_ by exposing the same features with less resources.
If a program can do the same thing with half the code, that's value added.
If it can do the same task in half the time, that's also value added.
If you could accomplish 90% of the existing features with only 10% of the code, that may be a big win!

## Idea 4: Other Paths to Success

What about writing your program or application directly in [assembly](https://en.wikipedia.org/wiki/Assembly_language)?
Oof, that'd hurt portability... right?
I'd say "yes" if the goal was to support multiple architectures.
However, I'd say "no" if a calculated choice was made to target a _specific_ architecture.
Given the title of the post, you probably know where I'm going this: [RISC-V](https://en.wikipedia.org/wiki/RISC-V).
RISC-V is an open standard instruction set architecture released under open source licences.
What if you were to simply invest in the RISC-V ISA and embrace its assembly dialect as the foundation for modern programs?
Think about how much compiler and operating system cruft you could leave behind!

If you were to hit the reset button and completely start over, what would it take to achieve usefulness?
How many instructions?
How many lines of assembly?
Do we need C or Rust?
Assembly may not be the most productive language but there are simple, high-level alternatives.
I think that [Forth](<https://en.wikipedia.org/wiki/Forth_(programming_language)>) is definitely a contender and worth taking a look at.
A minimal Forth implementation can be built in a few hundred instructions and then a useful, interactive system can built in a few hundred Forth words.
That sounds pretty efficient!
Perhaps it would even be fairly secure.

**"There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies, and the other way is to make it so complicated that there are no obvious deficiencies. The first method is far more difficult."** - C. A. R. Hoare

## Idea 5: The Bronzebeard Project

[Bronzebeard](https://github.com/theandrew168/bronzebeard) is a minimal ecosystem for bare-metal RISC-V development.
I've written a basic, standalone [assembler](https://github.com/theandrew168/bronzebeard/blob/master/bronzebeard/asm.py) with no dependency on existing toolchains.
I also wrote a useful subset of [DFU](https://en.wikipedia.org/wiki/USB#Device_Firmware_Upgrade) in order to support running programs on a few existing RISC-V boards.
Can I make small-scale devices such as the [Longan Nano](https://www.seeedstudio.com/Sipeed-Longan-Nano-RISC-V-GD32VF103CBT6-Development-Board-p-4205.html), [Wio Lite](https://www.seeedstudio.com/Wio-Lite-RISC-V-GD32VF103-p-4293.html), and [HiFive1 Rev B](https://www.sifive.com/boards/hifive1-rev-b) useful without all of the default frameworks and SDKs?
If not using those, then what does the toolset look like?
Assembly is really the lowest level you can effectively go so that's where I've started.

In Abelson and Sussman's [Structure and Interpretation of Computer Programs](https://mitpress.mit.edu/sites/default/files/sicp/full-text/book/book-Z-H-10.html#%_sec_1.1), they describe the core elements of programming:

- **primitive expressions** - which represent the simplest entities the language is concerned with
- **means of combination** - by which compound elements are built from simpler ones
- **means of abstraction** - by which compound elements can be named and manipulated as units

How simply and effectively can these elements be achieved when starting from scratch?
Forth is the most minimal approach I've seen.
Furthermore, I think that Forth's design is rooted so literally upon these three elements that it almost _encourages_ productive abstrations.
It limits you to the "good" kind of abstraction: combining smaller units of low-level functionality into larger and larger high-level procedures.
I consider a "bad" abstraction to be when extra code is added in preemptive concern for future changes that rarely occur.

## Conclusion

In short, there is simply too much code and complexity underlying most modern systems.
The rise of RISC-V will enable programmers to embrace and control hardware at its lowest and most powerful level.
Programming languages such as C are valuable for their ability to paint over the machine that they target.
Perhaps that isn't always a good thing.
I think that there is sometimes value in doing just the opposite: making a conscious and well-intended choice to utilize a single platform to its fullest.

All of these ideas are heavily opinionated without a lot of science or evidence.
They are simply my opinion based on what I've experienced and the research I've done.
Modern operating systems and general computing are definitely a marvel but they are far from perfect.
There are many signs that point to things actually being quite poor.
If people ignore this problem then nothing will ever change.
If people just think the same way they always have then nothing will ever change.
If people assume that what we have now is "good enough" then nothing will ever change.

I want to enable and encourage other developers to explore bare-metal development.
I want to build an ecosystem that accessible to everyone with a computer, regardless of their operating system or prior programming experience.
By investing in the RISC-V ISA and its ecosystem, we no longer need an abstration over the assembly language by which it speaks.
RISC-V assembly can be a new foundation for future servers, programs, and operating systems.
