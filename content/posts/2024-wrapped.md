---
date: 2024-12-29
title: "2024 Wrapped"
slug: "2024-wrapped"
draft: true
---

2024 was a year of great personal accomplishment: both personal and professional.
From new friends to new skills, a lot has changed and I've experience a lot of growth!
2024 was also the year of balance: balancing work, life, and the nuances within each.
Within work, I balanced regular development work with tech leading and planning.
Within life, I balanced relationships: welcoming new friends who all come with unique and exciting interests, experiences, and perspectives.

# Home Ownership

Personally, my wife and I bought a house!
We are home owners for the first time our lives and it has been amazing.
Sure, there is a lot of work to do, but overall it has been a tremendous benefit.
We both work from home and have separate offices now.
Our dog, Maple, has a big yard to run around it.
The ground-floor living room has no TV: just a bay window, plenty of seating, and our vinyl setup.
We jokingly refer to this as our "regen" living room.
Downstairs, we have a 2x3 arrangement of short floor couches (corgis have short legs) that puts us at the perfect height for TV-realted activities.
This is the "degen" living room.

# Technical Leadership

Professionally, I saw a lot of growth in 2024.
At work, I took my first dive into technical leadership: leading teams of 2-3 through multi-month software projects.
I learned techniques ([Mikado Method](https://mikadomethod.info/) and [Two Ways to Build a Pyramid](/posts/two-ways-to-build-a-pyramid/), for exmaple) for breaking down large, abstracts features into bite-sized, actionable tasks.
The biggest thing I learned is this: if the tickets are clear and correctly-sized, developers are able to complete them at a quick, consistent pace.
This, in turn, leads to increased team morale since everyone feels like they are crushing it (because they are)!
Selfishly, this reflects positively on my planning and leadership capabilities.
Since the year began, I've led three projects to successful, on-time completion.
The first few definitely caused me quite a bit of stress but things have gotten easier.

# Weekly Blogging

At the start of 2024, my buddy Nick and I set a goal for ourselves: write a blog post every week.
I'm very hyped to say that I've made it!
This is my final, 52nd post of the year.
Some posts were good, some posts were… less good.
But overall I think these writings will serve as a great additional to my portfolio as a professional developer.
I (well, Bloggulus) got a very minimal shout out from Rachel by the Bay.
That really made my week!

Was it fun? Sometimes.
To be honest, writing these posts rarely felt “more-ish” to me.
Rarely did I find myself super eager and excited to write a post (thought it did happen a few times).
Initially, I put a lot of pressure on myself to make each post long, detailed, and objectively useful.
I let some of that pressure go midway though the year, though.
Instead, I let myself write about anything: what I did that week, opinions on non-technical subjects, and the occasional technical deep dive.

Now that I'm at the end, it feels nice to look back and reflect on everything I've written.
Some posts fill me with joy to re-read while others fall a bit flat (they can't all be winners).
Here are some of my favorite posts from the past year (in no particular order).

## [1. Automating a Golden Age Minecraft Server](/posts/automating-a-golden-age-minecraft-server/)

This was a really fun post about how I researched, configured, and automated the hosting process for a "Golden Age" Minecraft server.
This refers to servers of older versions: [Beta 1.7.3](https://minecraft.fandom.com/wiki/Java_Edition_Beta_1.7.3) for me.
I ran the server on a moderately-specced DigitalOcean droplet and a few friends played on it.
It was an interesting problem from a security perspective since these older servers aren't compabitle with the game's modern auth systems.
Instead, I had to implement a "best effort" security posture by allow-listing usernames and restricting IP addresses.
I still really love the aesthetic of early Minecraft builds: cobblestone, trees, and planks.

## [2. A Multi-Plateform WebGL Demo](/posts/a-multi-platform-modern-webgl-demo/)

This post was a fun throwback to the [second blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/) I ever wrote (back in June 2020).
That post was about writing and building a cross-platform OpenGL demo in C.
The demo itself was pretty basic: a spinning red square.
Despite its simplicity, it required dealing with quite a few aspects of the OpenGL API: buffers, vertex arrays, and shaders.
It used [SDL2](https://www.libsdl.org/) for handling the native window and input support.

One downside of writing native OpenGL programs is that they are difficult to distribute.
As it turns out, modern computers don't want to let users run untrusted binaries (and fair enough)!
It meant that I had to tell my friends: "Just download, trust, and run this binary. I swear it isn't a virus!".
So, for the throwback post, I wanted to recreate the same demo but in a way that could run in the reader's browser.
This is where [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API) (a browser-based implementation of the modern OpenGL API) comes in!
In just a few hundred lines of JavaScript, I had the same demo running within the blog post itself.
I really loved how this post showed growth: both in technology (OpenGL to WebGL) and in my own experience (C to JavaScript).

## [3. Mario Kart and the Maker's Schedule](/posts/mario-kart-and-the-makers-schedule/)

This was an opinion piece I wrote reflecting upon Paul Graham's [Maker's Schedule, Manager's Schedule](https://www.paulgraham.com/makersschedule.html) essay as well as my own experiences.
Often, deep work requires prolonged, uninterrupted periods of focus.
A fragmented schedule full of meetings can disrupt this flow and lead to stress, anxiety, and feeling behind.
This post summarizes these sentiments and compares the two schedules (Maker vs Manager) to Mario Kart.
Some karts have slow acceleration but a high top speed (this is how I categorize myself).
Others have a faster acceleration but at the cost of a lower top speed.
I'm sure some folks have both, though.
If you do, that's awesome!

This post had great pacing and flow, in my opinion.
Early on in the post, I used the example of how my dog, Maple, can break my flow while barking.
Then, near the end, I talk about how I sometimes do my most productive work late at night when work chat is quiet.
In addition, Maple is usually snoozing at night so I was able to callback to her prior disruption by including a photo of her sleeping on the couch.
She's such a cutie.
It was a nice "circling back" moment to me which was further boosted by including a cute dog pic.

## [4. Bloggulus Outage Postmortem](/posts/bloggulus-outage-postmortem/)

This post was a technical postmortem into a ~7 hour outage that affected my [Bloggulus](https://bloggulus.com/) web site.
While no one was actually affected by the outage (I'm the only user and I was sleeping), it was an interesting case study into both what went wrong and why it took so long for me to notice.
The post really went down the failure rabbit hole: the app was down, the database was down, and the server's locale was messed up.
Ultimately, the cause was quite simple: the server ran out of memory why trying to regenerate the locale database after automatic updates.
The fix was also simple: enable and allocate swap space on the server.
I even wrote a [blog post](/posts/simple-server-swap-space/) about that, too.
For what it's worth, Bloggulus has no outages (or even OOM-killed processes) since!

## [5. Reinforcing Indirect Joins](/posts/reinforcing-indirect-joins/)

This post summarized an issue I kept running into at work where database queries were taking a shocking amount of time to plan.
To clarify, the actual execution of the query was fast, but the planning was taking multiple orders of magnitude longer.
I eventually identified a pattern shared between these slow queries: they were trying to join between "indirect" tables.
The post goes into more detail about what this means and how to fix it.
In the months following the post, I kept encountering and fixing these indirect joins.
I can't even count how many times I would look at query, see the indirect join, and then immediately slash the planning time by adding a single extra join.
I used this post to explain the problem to coworkers during PRs, too.
It came in super handy was probably my most tangibly impactful post of 2024.

# Personal Projects

2024 was also a great year for my personal projects!

## Bloggulus

As always, [Bloggulus](https://bloggulus.com/) saw quite a few enhancements.
I added the ability for users to register (via OAuth) and personalize the site with their own favorite blogs.
I made the choice to _only_ support OAuth authentication which came with some [interesting tradeoffs](/posts/oauth-auth-only/).
This project has really been my personal tech stack playground.
It started as server-side rendered Go program, then I wrote another version of the app in SvelteKit, migrated the original to a React-based SPA architecture, then back to server-side templating.
It made a big circle back to the original plan but I learned so much from the process.
There is no free lunch and all tech stacks have pros and cons: there is rarely a "best choice" across all criteria.

The current version of the app also uses a bit of HTMX for reload-free interactivity.
I've always to mess with HTMX so this was a great excuse.
Overall, I like it!
I think it strikes a nice balance between server-side simplicity and client-side interactivity / responsiveness.

I also use this project to explore multiple CSS strategies.
First was [TailwindCSS](https://tailwindcss.com/).
This was super hyped in 2024 so I figured I'd give it a shot.
While I really like the amount of pre-made components out there, I found that I quickly lost my hold on "what is this CSS actually doing and why is it necessary".
I'd eventually end up with a massive wall of classes and have no idea _why_ they were there.
The UI looked correct but I didn't fully understand this.

I hate this feeling of not grokking my own code and try to avoid it.
So, I ditched TailwindCSS for something simpler (and kinda old school): vanilla CSS with [BEM](https://getbem.com/introduction/) (block, element, modified) naming.
Honestly, I love this.
I get all the power of vanilla CSS (it has come a long way and doesn't need much extra) with specific, granular naming.
I never worry about naming collisions or style leaks.
Perfect solution to my current problems, honestly.

## Advent of Code

2024 was the second year in a row (third year ever) that I completed [Advent of Code]().
This challenge of increasingly-difficult programming puzzles is always a double-edged sword for me.
While I love the feeling of accomplishment for figuring out the puzzles, I can easily get obsessive and fixated on tough problem.
At its worst, I'll spend 8+ hours of a day banging my head against the keyboard just trying to figure things out.
My time management suffers and my mood suffers.
Is it worth?
Is it a net positive experience for me?
Honestly, I don't know.

But! I did finish again this year and only had a couple days of problematic frustration.
I do check the subreddit for hints if I'm completely stuck, though.
Since I finished, I ordered myself the real trophy: the Advent of Code coffee mug.
This will live on the shelf next my mugs from victories in prior years.
I love the ability to write so many "throwaway" programs in quick succession.
These days, I only use Python for AoC and it feels good to go back to it.
Python was the first language I learned deeply so I'll always have a nostalgic fondness for it.

I got to refresh my knowledge of some classic algorithms ([Dijkstra's](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)) and even learned some new ones ([Bron-Kerbosch](https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm)).
We'll see if I can stay away from Advent of Code in 2025...

## DerzPlates

Despite playing WoW since roughly 2007, I never applied my software engineering mindset to it.
That changed just a few weeks ago when I wrote and published my first WoW addon: [DerzPlates](https://github.com/theandrew168/derzplates).
I wrote a [blog post](/posts/derzplates-my-first-wow-addon/) about the process and how the WoW Interface forum community really showed up for me.
I've written Lua in the past so that wasn't new to me.
The WoW client API, however, was totally foreign and required so much research to even partially grok.
My initial version was kinda buggy but SDPhantom gave me solid guidance and advice to get it working properly.
They get my personal "Best Tech Advice of 2024" award, for sure!

## Health and Body

I turned 30 this year!
I also decided to take health a bit more seriously, too.
I know 30 isn't even that old but I want to lock in my good habits now and focus on activites that I can still do as I get older.
I did a lot more walking this year: long walks between 60 and 90 minutes 3-4 times per week.
I like it for music, podcasts, and even audiobooks.
I didn't bike as much as I did two years prior (the new house kept us busy).
2025 will include even more walking, more biking, more kayaking, and maybe even some running.

# Looking Forward

What do I hope to achieve in 2025?
I still want to blog, but likely not every week.
I also plan on [writing more Clojure](/posts/2025-the-year-of-clojure/) which I'm super excited for!
I'm not sure if I'll do AoC again.
On December 25th, I always tell myself that it wasn't worth the stress and that I won't do it next year... but then I can't stay away.
It's the hype!
People at work are talking about it, friends are talking about it, and I _know_ that thousands of devs are having a great time with the early puzzles.
