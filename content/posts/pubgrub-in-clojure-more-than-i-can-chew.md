---
date: 2024-12-01
title: "PubGrub in Clojure: More Than I Can Chew?"
slug: "pubgrub-in-clojure-more-than-i-can-chew"
draft: true
---

One of my recent side projects has been to study, understand, and implement the [PubGrub](https://nex3.medium.com/pubgrub-2fb6470504f) version solving algorithm.
I simply find it interesting!
I also thought: maybe this would also be a good reason to explore [Clojure](https://clojure.org/) (since 2025 is the [Year of Clojure](/posts/2025-the-year-of-clojure/), after all).

However, after making steady early progress (the easy parts), I found myself somewhat stuck when implementing the algorithm’s main loop.
How do I even do a straightforward “while true” loop in an immutable, functional language?
I then realize that I might be mixing too many variables, to my own detriment.

The vars are:

1. Learning a new programming language
2. Learning a new programming paradigm
3. Learning a new algorithm (from a family I’ve not experienced: SAT solvers)

Maybe I should alter my plan a bit and isolate these variables.
To that end, I made a small pivot.
I implemented the algorithm in a familiar, “good for algorithms” language: [Python](https://www.python.org/).
I use Python for [Advent of Code](https://adventofcode.com/) so I’m a bit more familiar with whipping up quick, alg-heavy prototypes.
This isolates the third variable.

Then, once I fully understand the alg and have a reference implementation, I’ll write the code again in Clojure.
This isolates variables 1 and 2 (which are super closely related).
Lessons learned: don’t try to mix too many new things at once.
This is also why I don’t usually recommend folks attempt Advent of Code with an unfamiliar language.
If you are aiming for 25/25, you’re gonna have a bad time.
