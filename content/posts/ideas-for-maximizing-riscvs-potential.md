---
date: 2020-11-02
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

# Idea 1: Where Things Went Wrong
slide from JB's talk: about high level langs and that something went wrong somewhere  
is Linux OK? Windows? or have modern OS's gone too far?  
Rene Rebe talks about this: microkernel, etc  
Is it okay that all OS API's are C?  
Or is that an issue for security / maintainability?  

# Idea 2: Too Much Code, Too Much Complexity
how do code quantity and complexity relate? what is their relationship?  
on their own, they are both sort of empty measures, just words  
but at a point, it does start to detract  
even if the code is great, at a scale where no single person can understand it all, its GG  
its why I appreciate Rust, but I don't think "safer" languages is the answer  

any solution that involves piling more stuff on top is NOT the answer  
I like Go because it cut stuff out! Bootstraps from ASM and itself  
is there just too much code  
no one can understand it all  
how could cross-cutting bugs ever be fixed  

# Idea 3: What Makes Something Useful?
do you need a full multi-user OS to be useful  
what about a webserver?  
is that _less_ secure because of the multi-user model  
what about performance?  

especially in the embedded space  
if you need to control a motor, do you need a multi-user OS?  
do you need MQTT and lambda functions and cloud-based infra?  

how do you express that net value added by achieving the same with less?  
or even the value added if you accomplish 90% features with 10% of the code  
its def a win for readability and maintainability  
easiest code to test / maintain is code that doesn't exist  

# Idea 4: Other Paths to Success
talk about investing in RISC-V assembly  
does portability matter as much here?  
plopping Linux on RISC-V would waste some potential  
what other ways are there to achieve interactivity and solve problems?  

can we get there with a few thousand lines of code instead of 30 million?  
or can we get to 90% of the usefulness with only 10% of the code  
devs could actually understand the entire system  
maybe that is better for security than choice of language  

# Idea 5: The Bronzebeard Project
Bare-metal RISC-V Forth implementation  
Uses a simple assembler co-written in both Python and Go  
Can I make these little devices useful without using heavy toolchains, frameworks, or SDKs?  
What toolset is best for that?  
Raw assembly is a given: you can't avoid that, but you can avoid the toolchains  

How can you quickly and minimally achieve the SICP 3: primitives, combination, abstration?  
Would this scale any better?  
Talk about exponential quantities of code to build bigger stuff  
Pull numbers from CollapseOS (everything is just a few hundred lines of Forth)  
Is it possible to keep LoC constant when building more and more advanced programs?  

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

A similar viewpoint can be taken when it comes to databases and ORMs.
While the the ability to target multiple databases may be valuable in some ways, there is also something lost.
Maybe SQLite has a specific feature that'd benefit your project.
By choosing an abstration layer, however, you limit your project to whatever functionality is common between each supported database.
You are now limited to the "common denominator" of database functionality.
That useful SQLite-only feature is now lost.
It is possible that explicitly choosing SQLite for the project may have been a better option.
