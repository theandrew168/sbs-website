---
date: 2025-03-30
title: "Brain Dump: OOP, FP, DDD, CQRS, EDA, Oh My!"
slug: "brain-dump-oop-fp-ddd-eda-cqrs-oh-my"
draft: true
---

Is DDD objectively good?
Regardless of language paradigm (imperative, OOP, functional, etc)?
I’m not sure.
Do functional / data-driven languages care about domain modeling?
Is it still a good idea to structure your data and pure functions around concepts of the domain?
What about Data-Driven Design?
Are they mutually exclusive?

Is [FCIS](https://www.destroyallsoftware.com/screencasts/catalog/functional-core-imperative-shell) objectively good?
Yes, I think so.
I think this is an amazing idea regardless of paradigm: [hoist your IO](https://www.youtube.com/watch?v=PBQN62oUnN8) and logic so that testing is easier.
I think this should always be embraced.
Isn't this just onion arch / hexagonal arch / clean arch?

Is EDA objectively good?
I’m not sure.
I need to read more about it.
I do like how it decouples the “what” from the “how”.
It also allows multiple, disparate systems to react to (and submit) events of various types.
I’ve already seen instances of how this would be nice at work.
How does EDA differ from using a queue for long-running background jobs / tasks?
Would it be redundant to have both?

Are these all related?
Are these all just individual pieces of a “good software design” puzzle?
Some facets seem at odds: OOP vs functional, DDD vs DDD, etc.

Maybe I’m just getting lost in the terminology.
What is OOP? Just a group of functions operating on a shared piece of data via single dispatch (and maybe some inheritance).
What is FP? Just structuring your code with an emphasis on pure functions and hoisting IO (basically FCIS).
What is EDA? Just decoupling the creation vs action of events in your codebase. Right?
