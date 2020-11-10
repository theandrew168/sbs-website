---
date: 2020-11-10
title: "Ideas for Maximizing RISC-V's Potential"
slug: "ideas-for-maxinimizing-riscvs-potential"
tags: ["risc-v", "assembly", "forth"]
draft: true
---
There has been a lot of buzz surrounding the slowness, instability, and complexity of modern software systems.
It seems as though more and more people are feeling and observing bloat and bugginess in their day-to-day usage of computers.
In many situations, the mediocrity of today's technology has become so common that it isn't even seen as an issue.

During my few years spent in industry, these sort of quotes were heard almost daily:

* *Looks like Windows crashed. What a surprise!*
* *Shoot, my PC decided to do updates... see you in an hour.*
* *This meeting room isn't working. Let's try to schedule a different one.*
* *Just keeping refreshing the page until the content renders correctly.*

Businesses and individuals lose countless hours to these types of problems every day but they are usually brushed off with a laugh.
The stability of Windows is the source of many jokes.
Backup plans are necessary for inconsistent meeting room reservations.
Are these situations an inevitability of computer systems or has something gone wrong?

# Existing Opinions
In his blog post [Software Disenchantment](https://tonsky.me/blog/disenchantment/), Nikita Prokopov (tonsky) presents evidence for "industry's lack of care for efficiency, simplicity, and excellence".
In his presentation [The 30 Million Line Problem](https://caseymuratori.com/blog_0031), Casey Muratori argues that the interaction between software and hardware has the potential to be simplified.
In his article [The Website Obesity Crisis](https://idlewords.com/talks/website_obesity.htm), Maciej Ceglowski describes how websites have gotten progressively larger and slower.
In his talk [Preventing the Collapse of Civilization](https://www.youtube.com/watch?v=pW-SOdj4Kkk), Jonathan Blow compares modern software development to historical technological advancements.
Lastly, in his 1995 essay  [A Plea for Lean Software](https://cr.yp.to/bib/1995/wirth.pdf), Niklaus Wirth describes causes for "fat software".
1995!
This problem has a history.

# My Opinion
I agree with the sentiments expressed in each of these posts.
In some sense, it leaves me with a feeling of helplessness: what can we do about this?
Is the ship of modern software development too large to turn?
What is the real problem?
Even if we knew the answer to that, how could we fix our technological trajectory?
I don't want to accept that the tools we use everyday are incapable of being fast and reliable.

I've spent a lot of time thinking about this problem.
Not only what the root cause may be, but also what can be done to mitigate it.
The videos and articles mentioned prior have all helped guide me toward the following ideas.
In short, my ideas are the following:
1. The usage of high-level languages started off great but has gone too far
2. Modern software and tools has too much code and too much complexity
3. We don't need 100% of current functionality in order for something to be useful
4. There are other routes to achieving value than an OS and a bunch of C code
5. The Bronzebeard project is my attempt to validate my claims

# Idea 1: Where Things Went Wrong
In Jonathan Blow's talk, he identifies that somewhere along the path of building higher and higher abstractions, we went too far.
At some point, the return on investment for programming in higher-level languages started to flatten out.
Shouldn't our programmers be more productive and their programs more robust?
This doesn't seem to be the case.
Sure, technology accomplishes much more now than it used to, but does the modern scale of features match the amplified scale of complexity and fragility?

Are modern operating systems the epitome of software efficiency?
By this, I mean to ask if they bring enough value to the table to justify their costs.
What costs? Linux is free!
I'm not referring to monetary costs: I'm referring to complexity and the volume of code between you and chunk of silicon at the heart of your machine.

To be fair, I'm not really a good position to answer this question.
I have zero kernel development experience at all.
However, there do exist veteran programmers who take issue with the state of modern kernel development.
[René Rebe](https://rene.rebe.de/) supports an ongoing conversation about this exact topic.
Though his [YouTube channel](https://www.youtube.com/user/renerebe), René present frequent ideas and criticisms relating to operating system design and implementation.
He explains how better languages and better designs (such as [microkernel](https://en.wikipedia.org/wiki/Microkernel)) could be an important step in OS innovation.

# Idea 2: Too Much Code, Too Much Complexity
Most of the time, code quantity is really a poor measure of anything.
Talking about a project's "lines of code" doesn't usually add much value.
However, as a project's size increases, so does it's learning curve.
In my opinion, the easiest code to understand and maintain is code that doesn't exist.
A large amount of recent language debate has pertained to matters of safety.
Many agree that C is an "unsafe" language while Go and Rust are "safer" alternatives.
While I do concur that safety is important, I think that quantity may be more significant.

I think at some scale, the safety of a language doesn't really matter.
Thirty million line of code in _any_ language is going to result in subtle, cross-cutting bugs.
As soon as a project grows to the point where no single developer can understand the whole picture, I think it's game over.
At such a scale, achieving optimal performance and reliability becomes impossible because human intercommunication is lossy and inefficient.
I think that innovation will come not from writing _safer_ code, but from writing _less_ code.

Another issue that arises from a boundless code quantity is the matter of making changes.
If your foundation is a massive, unintelligible tangle of code, how you identify where to make the change?
Most of the time, the solution is just pile some more code on top of what's already there.
After enough years of this, iteration slows to a crawl.
Every change breaks something else and even minor feature additions can take months.

One of the things I really admire about the [Go]() programming language is its conscious choice to avoid native APIs and the C language in general.
By bootstrapping itself from assembly and making syscalls directly, Go avoids inheriting the baggage of operating systems.
By baggage, I'm referring to things such as a tightly-coupled `libc` on \*nix systems and the non-standardized `msvcrt.dll` on Windows.

# Idea 3: What Makes Something Useful?
Is a fully-featured, multi-user OS necessary to be useful?
Do focused applications such as a webserver need that?
In fact, is something like a webserver made _less_ secure because of the multi-user model and features behind it?
Perhaps even the performance of a dedicated server could be better without the extra overhead.
Especially in the embedded space, is a multi-user OS necessary?
If you had an IoT device or some simple controller, does that need the concept of a user?
How much code does something like that really need at all?

When it comes to "value", there are two equally-important angles to consider.
The first and most common angle is the value added by doing more stuff.
New features are the primary embodient of this.
More code added in order to do more things.
The second and less common angle is the value added _implicity_ by doing the same stuff with less resources.
If a program can do the same thing with half the code, that's value added.
If it can do the same task in half the time, that's also value added.
If you could accomplish 90% of the existing features with only 10% of the code, that's potentially a big win.

# Idea 4: Other Paths to Success
What other routes are there to making a computer useful?
Is installing a C-based operating system a required prerequisite?
You _could_ just write some bare-metal C or Rust, of course.
But even both of those carry with them a sizable amount of complexity in terms of their compilers and toolchains ([GCC](https://gcc.gnu.org/), [LLVM](http://llvm.org/), etc).
Don't get me wrong, these are amazing tools!
But if your goal is understand your entire stack, then adding either of them to your process is a major setback.

What about writing your program or application directly in assembly?
Oof, that'd hurt portability... right?
I'd say "yes" if the goal was to support multiple architectures.
However, I'd say "no" if a calculated choice was made to target a _specific_ architecture.
Given the title of this blog, you probably know where I'm going: RISC-V.
What if you were to simple invest in RISC-V and embrace its assembly dialect as the foundation for modern programs?
Think about how much compiler and operating system cruft you could leave behind!

If you were to hit the reset button and completely start over, what would it take to achieve usefulness?
How many instructions?
How many lines of assembly?
Do we need C or Rust?
Assembly may not be the most productive language but there are simple high-level alternatives.
I think that Forth is definitely a contender and worth taking a look at.
A minimal Forth implementation can be built in a few hundred instructions and then a useful, interactive system can built in a few hundred Forth words.
That sounds pretty efficient!
Perhaps it would even be fairly secure.

**"There are two ways of constructing a software design: One way is to make it so simple that there are obviously no deficiencies, and the other way is to make it so complicated that there are no obvious deficiencies. The first method is far more difficult."** - C. A. R. Hoare

# Idea 5: The Bronzebeard Project
Bare-metal RISC-V Forth implementation  
Uses a simple assembler written in both Python
Think the "arduino ecosystem of riscv"  
Native GUIs, beginner friendly ala Thonny  

Avoid the big toolchain!  
Can I make these little devices useful without using heavy toolchains, frameworks, or SDKs?  
What toolset is best for that?  
Raw assembly is a given: you can't avoid that, but you can avoid the toolchains  

How can you quickly and minimally achieve the SICP 3: primitives, combination, abstration?  
Would this scale any better?  
Talk about exponential quantities of code to build bigger stuff  
Pull numbers from CollapseOS (everything is just a few hundred lines of Forth)  
Is it possible to keep LoC constant when building more and more advanced programs?  
plopping Linux on RISC-V would waste some potential  

You _can_ do it with regular languages but you have to make the right design choices  
Does Forth's minimalism make this easier?  
Does Forth make it harder to paint yourself into a corner?  
Forth expresses purely aggregations of smaller ideas (like a more pure Lisp)  

# Conclusion
In short, there is simply too much code and complexity underlying most modern systems.
The rise of RISC-V will enable programmers to embrace and control hardware at its lowest and most powerful level.
Programming languages such as C are valuable for their ability to paint over the machine that they target.
Perhaps that isn't always a good thing.
I think that there is sometimes value in doing just the opposite: making a conscious and well-intended choice to utilize a single platform to its fullest.

this is all heavily opinionated without a lot of science  
its my opinion based on what i've seen and the research ive done  
modern OSs and computing is definitely a marvel but its not perfect  
there are a lot of signs that point to things being quite bad  
if people just think the same way, then nothing will ever change  
if poeple assume that what we have now is "good enough", then nothing will ever change  

I want to enable and encourage other developers to explore bare-metal development.
By investing in the RISC-V ISA and its ecosystem, we no longer need an abstration over the assembly language by which it speaks.
RISC-V assembly can be a new foundation for future servers, programs, and operating systems.
