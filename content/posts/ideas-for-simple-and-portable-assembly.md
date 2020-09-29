---
date: 2020-09-28
title: "Ideas for Simple and Portable Assembly"
slug: "ideas-for-simple-and-portable-assembly"
tags: ["python", "assembly", "forth"]
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

# So What Now?
I agree with those posts  
what can we do?  
Is ASM too low?  
Is C too high?  
Maybe too high with libc but not if freestanding?  
Maybe there is another path to high-level value?  
Portable ASM + Forth? Sounds fun, at least.  

# Idea 1: How High is too High?
slide from JB's talk: about high level langs and that something went wrong somewhere  
is Linux OK? Windows? or have modern OS's gone too far?  
Rene Rebe talks about this: microkernel, etc  

# Idea 2: Different for Arbitrary Reasons
again, JB's talk  
if code could be injected into a CPU, it'd run  
diff OS's have diff formats: ELF, PE, Mach-O  
is that progress?  

# Idea 3: Portability Required Planning
Rob Pike's talk on the Go assembler / compiler  
assembly dialects express mostly the same ideas  
same concepts, just maybe different syntax or approach  

# Idea 4: Generic Assembly
Magni ideas: exit, alloc, etc  
what do those look like on each target (use a table?)  
Magni ASM + OS + Arch = executable  
is that valuable or do you lose too much? (ASCII add before multiply)  
by not using specific insts, do you miss out?  
is it bad to just use a RISC subset of a CISC arch?  

# Idea 5: The Bronzebeard Project
Python-based project  
Simple, template assembler  
Includes builders for common EXE formats (and raw / hex for bare metal)  
Current use case is to build small progs for RISC-V devices  
A portable Forth impl would be cool, too  
