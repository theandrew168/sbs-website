---
title: "Device Ownership: Introduction"
date: 2020-02-21
tags: []
---
### Is there a problem with modern software?
There has been a lot of buzz surrounding the slowness, instability, and complexity of modern software systems.
It seems as though more and more people are feeling and observing bloat and bugginess in their day-to-day usage of computers.
In many situations, the mediocrity of today's technology-based working environment has become commonplace.

During my few years spent in industry, these sort of quotes were heard almost daily:

* *Looks like [Windows](https://en.wikipedia.org/wiki/Microsoft_Windows) crashed. What a surprise!*
* *Shoot, my PC decided to do updates... see you in an hour.*
* *This meeting room isn't working. Let's try to schedule a different one.*
* *Just keeping refreshing the page until the content renders correctly.*

Businesses and individuals lose countless hours to these types of problems every day but they are usually endured with a laugh.
The stability of Windows is the source of many jokes.
Backup plans are necessary for inconsistent meeting room reservations.
Are these situations an inevitability of computer systems or has something gone wrong?

### Some people say "yes!"
Nikita Prokopov (tonsky) presents evidence for "industry's lack of care for efficiency, simplicity, and excellence" in blog post [Software Disenchantment](https://tonsky.me/blog/disenchantment/).
Casey Muratori argues that the interaction between software and hardware has the potential to be simplified in his presentation [The 30 Million Line Problem](https://caseymuratori.com/blog_0031).
Maciej Ceglowski describes how websites have gotten progressively larger and slower in his article [The Website Obesity Crisis](https://idlewords.com/talks/website_obesity.htm).
Jonathan Blow compares software development to a historical context of technological advancement in his talk [Preventing the Collapse of Civilization](https://www.youtube.com/watch?v=pW-SOdj4Kkk).
Lastly, Niklaus Wirth described causes for "fat software" back in his 1995 essay [A Plea for Lean Software](https://cr.yp.to/bib/1995/wirth.pdf).

### What role does this series play?
I agree with the analysis and opinions laid out by the individuals listed above.
I want to find out if (and where) modern implementations can be simplified by starting at the bottom and tracing new paths up to what we have today.
This series serves as a journal of the research and development necessary to make computers perform useful tasks.
Due to its openness and simplicity, the [RISC-V](https://en.wikipedia.org/wiki/RISC-V) [ISA](https://en.wikipedia.org/wiki/Instruction_set_architecture) has been chosen as the canvas on which to explore these ideas.

### What questions do I want to explore?
How much code does it take to utilize all of the features of a given piece of hardware?
How much code does it take to start from scratch and completely "own" the device?
At such a low level, how portable is software from one RISC-V device to another?
Are there high-level languages other than [C](https://en.wikipedia.org/wiki/C_(programming_language)) that provide useful and performant abstractions of a system's underlying hardware?
Do any other models of [operating system](https://en.wikipedia.org/wiki/Operating_system) design have valuable tradeoffs over traditional [monolithic kernels](https://en.wikipedia.org/wiki/Monolithic_kernel)?
Will the lack of per-core licensing costs on RISC-V chips enable innovations that displace traditional kernels with leaner [hardware abstraction layers](https://en.wikipedia.org/wiki/HAL_(software)) and [schedulers](https://en.wikipedia.org/wiki/Scheduling_(computing))?

### What to expect
In preparation for this series, I've written a simple RISC-V assembler in the form of a [Python](https://www.python.org/) package named [simpleriscv](https://pypi.org/project/simpleriscv/).
Though potentially unsuitable for large projects, this package lowers the barrier to entry for those who are curious about RISC-V.
It saves beginners from needing to [build custom GCC toolchains](https://github.com/riscv/riscv-gnu-toolchain) for their specific RISC-V target(s).

I've chosen a few reasonably-priced RISC-V development boards to kick the series off.
At the lower end, we have the [Sipeed Longan Nano](https://www.seeedstudio.com/Sipeed-Longan-Nano-RISC-V-GD32VF103CBT6-Development-Board-p-4205.html) and [Wio Lite RISC-V](https://www.seeedstudio.com/Wio-Lite-RISC-V-GD32VF103-p-4293.html).
These chips are exciting because they have a small LCD screen and WiFi module, respectively.
In the middle of the road is SiFive's [HiFive1 Rev B](https://www.sifive.com/boards/hifive1-rev-b).
This chip has a WiFi module as well and comes with a completely different CPU than the first two.
I'm confident that this selection will provide enough variety to compare and contrast what it takes to write portable code at such a low level.

I hope you find this series interesting and informative!