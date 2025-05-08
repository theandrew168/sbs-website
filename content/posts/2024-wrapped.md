---
date: 2024-12-29
title: "2024 Wrapped"
slug: "2024-wrapped"
---

This year, 2024, was one of great accomplishment: both personal and professional.
From new friends to new skills, much has changed and I've experienced a huge amount of growth!
2024 was also the "Year of Balance": balancing the nuances of both work and life.
Within work, I balanced regular development work with tech leadership and project planning.
Within life, I balanced relationships and welcomed new friends who all come with unique and exciting interests, experiences, and perspectives.

## Home Ownership

My biggest personal achievement was buying a house!
My wife and I are home owners for the first time in our lives and it has been amazing.
Sure, there is a lot of work to do, but overall it has been a tremendous benefit.
We both work from home and finally have separate offices (no more working from the couch).
Our dog, Maple, now has a big yard to run around in.

The house's ground-floor living room has no TV: just a bay window, plenty of seating, and our ever-growing vinyl setup.
We jokingly refer to this as our "regen" living room.
Downstairs, we have a large, flat arrangement of short floor couches (corgis have short legs) that puts us at the perfect height for watching TV and playing video games.
This is known as the "degen" living room.

## Technical Leadership

Professionally, I saw a lot of growth in 2024.
At work, I took my first dive into technical leadership: leading teams of 2-3 developers through multi-month software projects.
I learned techniques for breaking down large, abstract features into bite-sized, actionable tasks (the [Mikado Method](https://mikadomethod.info/) and [Two Ways to Build a Pyramid](/posts/two-ways-to-build-a-pyramid/), for example).
One of the biggest and most important things I learned was this: if tickets are small, clear, and actionable, developers will be able to complete them at a quick, consistent pace.

This will, in turn, lead to increased team morale since everyone feels like they are crushing it (because they are)!
Selfishly, this reflects positively on my planning and leadership capabilities.
Since the year began, I've led three projects to successful, on-time completion.
The first few definitely caused me quite a bit of stress but things have gotten easier.
Tech leadership is a skill that gets better with experience just like anything else.

## Weekly Blogging

At the start of 2024, my buddy [Nick](https://nickherrig.com/) and I set a goal for ourselves: write a new blog post every week.
I'm very hyped to say that I've made it!
This is my last post of the year: number 52.
Some posts were good, some posts were… less good.
But overall I think these writings will serve as a great addition to my portfolio.
It also demonstrates my writing capabilities which are somewhat underrepresented by my GitHub profile.

Was it fun? Sometimes.
To be honest, writing these posts didn't often feel “more-ish” to me.
Rarely did I find myself super eager and excited to write a post (though it did happen a few times).
Initially, I put a lot of pressure on myself to make each post long, detailed, and objectively useful to any readers.
I let some of that pressure dissipate midway through the year, however.
Instead, I let myself write about anything: what I did that week, opinions on non-technical subjects, and smaller "Today I Learned" style posts.

Now that I'm at the end, it feels nice to look back and reflect on everything I've written.
Some posts fill me with joy to re-read while others fall a bit flat (they can't _all_ be winners, after all).
Here are some of my favorite posts from the year (in no particular order):

### [1. Automating a Golden Age Minecraft Server](/posts/automating-a-golden-age-minecraft-server/)

This was a really fun post about how I researched, configured, and automated the hosting process for a "Golden Age" Minecraft server.
The term "Golden Age" refers to versions of Minecraft prior to the [Adventure Update](https://minecraft.fandom.com/wiki/Adventure_Update) (I chose [Beta 1.7.3](https://minecraft.fandom.com/wiki/Java_Edition_Beta_1.7.3) for my server).
I still really love the aesthetic of early Minecraft builds: cobblestone, trees, and planks.
The post includes some nostalgic screenshots of the original server I hosted back in high school.

I ran the server on a medium-sized [DigitalOcean](https://www.digitalocean.com/) droplet for myself and a few friends.
It presented an interesting problem from a security perspective since these older servers aren't compabitle with the game's modern authentication system.
Instead, I had to implement a "best effort" security posture by allow-listing usernames and restricting IP addresses.
It definitely wasn't perfectly locked down, but we never had any issues with malicious actors showing up.

### [2. A Multi-Plateform WebGL Demo](/posts/a-multi-platform-modern-webgl-demo/)

This one was a fun throwback to the [second blog post](/posts/a-multi-platform-modern-opengl-demo-with-sdl2/) I ever wrote (back in June 2020).
That post was about writing and building a cross-platform [OpenGL](https://en.wikipedia.org/wiki/OpenGL) demo in C.
The demo itself was very basic: just a spinning red square.
Despite its simplicity, it required understanding and interacting with multiple aspects of the OpenGL API: buffers, vertex arrays, and shaders.
It used [SDL2](https://www.libsdl.org/) for handling the native window creation and input events.

A big downside of writing native OpenGL programs is that they are difficult to distribute.
As it turns out, modern computers don't want to let users run untrusted binaries (and fair enough)!
It meant that I always had to tell my friends:

> Just download, trust, and run this binary. I swear it isn't a virus!

So, for the throwback post, I wanted to recreate the same demo but in a way that could run directly in the reader's browser (no virus warnings necessary).
This is where [WebGL](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API) (a browser-based implementation of the modern OpenGL API) comes in!
In just a few hundred lines of JavaScript, I had the same demo running [within the blog post](https://github.com/theandrew168/sbs-website/blob/25a3797d8e16ef93fe49e3f068313148d6791275/content/posts/a-multi-platform-modern-webgl-demo.md?plain=1#L89-L222) itself.
I really loved how this post showed growth: both in technology (OpenGL to WebGL) and in my own experience (C to JavaScript).

### [3. Mario Kart and the Maker's Schedule](/posts/mario-kart-and-the-makers-schedule/)

This was an opinion piece I wrote as a personal reflection on Paul Graham's [Maker's Schedule, Manager's Schedule](https://www.paulgraham.com/makersschedule.html) essay.
Often, deep work requires prolonged, uninterrupted periods of focus.
A fragmented schedule full of meetings can disrupt this flow and can lead feeling anxious and stressed.
My post summarizes these sentiments and compares the two schedules to Mario Kart.
Some karts have slow acceleration but a high top speed (this is how I categorize myself).
Others have faster acceleration but at the cost of lower top speed.

This post had great pacing and flow, in my opinion.
In the beginning, I used the example of how my dog, Maple, can break my flow with her unexpected and surprising barks.
Then, near the end, I talk about how I sometimes do my most productive work late at night when work chat is quiet.
In addition, Maple is usually snoozing at night so I was able to callback to her prior disruption by including a photo of her sleeping on the couch.
It was a nice "circle back" moment which was further elevated by including a cute dog pic!

### [4. Bloggulus Outage Postmortem](/posts/bloggulus-outage-postmortem/)

This post was a technical postmortem into an outage that affected my [Bloggulus](https://bloggulus.com/) web site back in June.
While no one was actually affected by the outage (I'm the only user and I was sleeping), it was an interesting case study into both what went wrong and why it took so long for me to notice.
The post really went down the failure rabbit hole: why the app was down, why the database was down, and why the server's locale database was messed up.

Ultimately, the cause was quite simple: the server ran out of memory while trying to regenerate the locale database after a round of automatic updates.
The fix was also simple: enable and allocate swap space on the server.
I wrote a [blog post](/posts/simple-server-swap-space/) about that, too.
For what it's worth, Bloggulus has had no outages (or even OOM-killed processes) ever since!

### [5. Reinforcing Indirect Joins](/posts/reinforcing-indirect-joins/)

This post summarized an issue I kept running into at work where database queries were taking a shocking amount of time to plan.
To clarify, the actual **execution** of the query was fast, but the **planning** was taking multiple orders of magnitude longer.
I eventually identified a pattern shared between these slow queries: they were trying to join between "indirect" tables.

The post goes into more detail about what this means and how to fix it.
In the months following, I kept encountering, identifying, and fixing these indirect joins.
I can't even count how many times I would look at query, see the indirect join, and then immediately slash the planning time just by adding a single join.
I used this post to explain the problem to coworkers within PRs, too.
It came in super handy and was probably my most objectively impactful post of 2024.

## Personal Projects

2024 was also a great year for personal projects!

### Bloggulus

As always, [Bloggulus](https://bloggulus.com/) saw quite a few enhancements.
I added the ability for users to register (via OAuth) and personalize their experience with their own favorite blogs.
I made the choice to _only_ support OAuth authentication which came with some [interesting tradeoffs](/posts/oauth-auth-only/).
The current version of the app also uses a bit of [HTMX](https://htmx.org/) for reload-free interactivity.
I'd been wanting to mess with HTMX so this was a great excuse.
Overall, I like it!
I think it strikes a nice balance between server-side simplicity and client-side responsiveness.

This project has really been my personal tech stack playground.
It started as server-side rendered Go program.
Then, for fun, I wrote another version of the app in [SvelteKit](https://svelte.dev/).
After that, I migrated the original codebase to a [React](https://react.dev/) single-page application.
Eventually, I landed right back where I'd started: with server-side templating.
Despite the circle of redundancy, it was far from a waste a time because I [learned](/posts/brain-dump-bffs-and-api-calls/) so [much](/posts/has-science-gone-too-far/) from the process.
There is no free lunch and all tech stacks have pros and cons: there is rarely a "best choice" across all criteria.

I also use this project to explore multiple CSS strategies.
First came [TailwindCSS](https://tailwindcss.com/).
This was super hyped in 2024 so I figured I'd give it a shot.
While I really liked the availability of pre-made, "off the shelf" components, I found myself quickly losing track of what each class actually did and why it was necessary.
I'd eventually end up with a massive wall of classes and have no idea _why_ they were there.
The UI might've looked absolutely amazing but I didn't fully understand it.

I hate this feeling of not grokking my own code and ultimately try to avoid it.
So, I ditched TailwindCSS for something simpler (and kinda old school): vanilla CSS with [BEM](https://getbem.com/introduction/) (Block, Element, Modifier) naming.
Honestly, I love it.
I get all the power of vanilla CSS (which has come a long way and doesn't even need many extra features) with specific, granular naming.
I never worry about naming collisions or style leaks.
It's the perfect solution for the scale of web sites that I build, honestly.

### Advent of Code

2024 was the second year in a row (third year ever) that I completed [Advent of Code](https://adventofcode.com/).
This challenge of solving increasingly difficult programming puzzles is always a double-edged sword for me.
While I love the feeling of accomplishment that comes with every successful solution, I can easily get obsessive and fixated on tough problems.
At its worst, I'll spend 8+ hours of a day banging my head against the keyboard trying to work things out.
My time management suffers and my mood suffers.
Is it worth?
Is it a net positive experience for me?
Honestly, I don't know.

But! I did finish again this year and only encountered a couple days of problematic frustration.
I do check the [subreddit](https://old.reddit.com/r/adventofcode/) for hints if I'm completely stuck, though.
Since I finished, I ordered myself the real trophy: the Advent of Code coffee mug.
This will live on the shelf next to my other victory mugs.

I love the ability to write so many "throwaway" programs in quick succession.
These days, I typically only use Python for solving these puzzles and it feels good to come back home.
Python was the first language I deeply learned so I'll always have a nostalgic fondness for it.
I got to refresh my knowledge of some classic algorithms ([Dijkstra's](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm)) and even learned some new ones ([Bron-Kerbosch](https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm)).
We'll see if I can stay away from Advent of Code in 2025...

### DerzPlates

Despite playing WoW since roughly 2007, I never applied my software engineering mindset to it.
That changed just a few weeks ago when I wrote and published my first WoW addon: [DerzPlates](https://github.com/theandrew168/derzplates).
I wrote a [blog post](/posts/derzplates-my-first-wow-addon/) about the process and how the [WoW Interface](https://www.wowinterface.com/community.php) forum communityreally showed up for me.

I've written [Lua](https://www.lua.org/) in the past so that aspect of addon development was familiar to me, at least.
The WoW client API, however, was totally foreign and required so much research to even partially understand.
My initial version was kinda buggy but SDPhantom (on the forums) gave me some incredible [guidance and advice](https://www.wowinterface.com/forums/showthread.php?p=344701) to get it working properly.
They get my personal "Best Tech Support of 2024" award, for sure!

## Looking Forward

What do I hope to achieve in 2025?
I still want to blog but it likely won't be every week.
I also plan on [writing more Clojure](/posts/2025-the-year-of-clojure/) which I'm super excited for!
Software development is one of my favorite hobbies so I'm sure some new and exciting projects will reveal themselves.
I'm also hoping to go on more long walks (perhaps with some mild jogging mixed in).
Walking with music, podcasts, or audiobooks has really been my jam this year and I plan to keep it up.
Weather permitting, I also want to do more biking on the local trails and kayaking on the river.

This was an amazing year for personal and professional growth.
I honestly think it might've been my most successful and satisfying year ever.
It really set the bar high but hopefully 2025 can be even better.
Either way, I'll be sure to summarize the year's happenings in another "wrapped" post.
See you all next year!
