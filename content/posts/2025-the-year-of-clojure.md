---
date: 2024-09-29
title: "2025: The Year of Clojure?"
slug: "2025-the-year-of-clojure"
tags: ["Clojure"]
---

I’ve been thinking about [Clojure](https://clojure.org/) a lot lately.
Something about this “immutable, data-driven lisp on the JVM” has always fascinated me… but has never “clicked”.
I love its core ideas: immutable data, simple syntax, emphasis on pure functions, great interop with Java, etc.
In fact, my [modern patterns](/posts/a-better-pattern-for-go-http-handlers/) for writing web apps in Go are slowly converging on a more functional design (using higher-order function to close around dependencies).
I’d like to write more Clojure in 2025 and see if I can finally connect deeply with its philosophy.

## The Year of Clojure

My experience with the language is fairly limited: I re-implemented Alex Edwards' [snippetbox](https://github.com/theandrew168/snippetbox) application in Clojure while waiting for a delayed flight back in 2022 (you should definitely read his [Let's Go](https://lets-go.alexedwards.net/) book, by the way).
While the experience was fun, I feel like I simply translated the code directly from Go to Clojure without doing much "idiomatic alignment".
I think a Clojure expert would implement a project like that much differently.
That being said, Clojure (well, the Java ecosystem) does have some extra perks that make it appealing to me: single-file deployments and easily-embeddable resources.
Gotta love JARs!

Every talk I watch about Clojure fascinates me.
Here are just a few recent examples:

1. [Commander Pattern / CQRS](https://www.youtube.com/watch?v=B1-gS0oEtYc) - An architectural deep dive into Command Query Request Segregation (CQRS)
2. [Abstracting Data With Maps](https://www.youtube.com/watch?v=Sjb6y19YIWg) - How to model your domain with maps while avoiding "map fatigue"
3. [Everything Will Flow](https://www.youtube.com/watch?v=1bNOO3xxMc0) - An incredible talk about queueing theory and queues in practic
4. [Clojure Spec Data Constraints](https://www.youtube.com/watch?v=Xb0UhDeHzBM) - An overview of `clojure.spec` and how it can be used for driving tests with data
5. [11 Datomic Insights](https://www.youtube.com/watch?v=YSgTQzHYeLU) - A great introduction to Datomic and how it differs from other databases
6. [Code Review of a Lexer](https://www.youtube.com/watch?v=SGBuBSXdtLY) - A short and sweet code review of a lexer implemented as a single function

Even the Clojure-adjacent databases are interesting:

1. [Datomic](https://www.datomic.com/) - A closed-source, immutable database implemented with the "triple-store" data model that uses [Datalog](https://docs.datomic.com/whatis/data-model.html#datalog) for queries
2. [XTDB](https://xtdb.com/) - An open-source, immutable database that uses SQL for queries

# The Year of Closure

In my head, 2024 was the “Year of Balance”.
It was about trying to appreciate nuance: the subtle grays between black and white.
There is always good and bad: perfection is rare.
Sometimes in life, we get busy and sometimes we struggle.
In those times, some people withdraw from personal connections and focus on themselves.
I think is fine in moderation.
After all, you have do whatever you can to work things out: even if it means focusing on yourself for a while.

But, once the storm passes, should you reach back out?
In my friend group, we often refer to this as "ghost busting".
Should this idea also be a consideration for next year?
I was chatting with a friend who feels like maybe it should.
Perhaps 2025 can also be the “Year of Closure” (homophone intended).

Thanks for reading!
