---
date: 2024-12-01
title: "PubGrub in Clojure: More Than I Can Chew?"
slug: "pubgrub-in-clojure-more-than-i-can-chew"
---

One of my recent [side projects](https://github.com/theandrew168/pubgrub) has been to study, understand, and implement the [PubGrub](https://nex3.medium.com/pubgrub-2fb6470504f) version solving algorithm.
I'm not building any particular application or product with it, I simply find it interesting!
I also thought to myself: maybe this would also be a good reason to explore [Clojure](https://clojure.org/) (since 2025 is the [Year of Clojure](/posts/2025-the-year-of-clojure/), after all).

However, after making steady early progress (the easy parts), I found myself somewhat stuck when implementing the algorithm’s main loop.
How do I even do a straightforward “while true” loop in a functional and immutable language?
Realizing that I didn't yet understand how to translate even simple constructs into Clojure's terms, I started to feel like I'd bitten off more than I can chew.

Perhaps I am trying to work too many new concepts at once:

1. Learning a new programming language (Clojure)
2. Learning a new programming paradigm (Functional)
3. Learning a new algorithm family ([SAT Solving](https://en.wikipedia.org/wiki/SAT_solver))

Maybe I should take a beat, reflect on where I'm at, and alter my plan to better isolate these concepts.
After doing so, I decided to make a small pivot.
Instead of starting with Clojure, I'll instead implement the algorithm in a familiar and algorithm-friendly language: [Python](https://www.python.org/).
I use Python for [Advent of Code](https://adventofcode.com/) so I’m quite a bit more familiar with writing short, focused, algorithm-heavy programs.
This isolates the third concept and lets me focus on the hard part first: PubGrub.

Eventually, once I fully understand the algorithm and have written my own reference implementation, I’ll write the code again in Clojure.
This isolates concepts 1 and 2 which are similar enough to tackle at the same point (plus, you can't really separate them in the first place).
Lesson learned, though: don’t try to mix too many new things at once.

This is also why I don’t usually recommend that folks attempt Advent of Code with an unfamiliar language.
In my opinion, if you want to obtain all 50 stars with a new language, you’re gonna have a bad time.
If you care more about learning and less about the stars, then I say go for it!
