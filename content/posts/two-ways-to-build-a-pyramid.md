---
date: 2024-05-12
title: "Two Ways to Build a Pyramid"
slug: "two-ways-to-build-a-pyramid"
---

Recently at work, I've been overseeing my first project as a tech lead / project manager.
Myself and two other developers are working together to rewrite and revamp a section of the application's frontend.
In preparation, I needed to break down a set of requirements and design mocks into small, actionable tasks.

Originally, I tried to plan the entire project upfront (we have roughly eight weeks of runway to complete it).
This required a huge amount of time!
The further into the future I planned, the harder it became to carve off specific tasks.
This was due to the uncertainty surrounding how the project would look at that point in time.
As another downside, this approach would've deferred delivering a usable frontend until the last few weeks of the project.
I eventually discovered that there was a better way to go about all of this.

## Iterative Development

What I realized was the importance of "iterative development".
In short, this means delivering a tangible, usable artifact at much smaller intervals throughout the project's lifecycle.
As soon as possible, try to get some meaningful subset of the project working.
Then, add features incrementally such that the application is still usable (but with more functionality) after each step.
Despite being loosely familiar with iterative development, I didn't truly appreciate its important until holding a project management position.

What I had done before is sometimes called "traditional development".
I was trying to plan the entire project upfront, build it from the bottom up, and only deliver something usable at the very end.
This isn't a great situation to be in!
Waiting until the end (or near the end) of a project to get feedback from stakeholders is rife with problems.

The sooner you can get a demo in front of the product and design folks, the better off everyone will be.
Allowing stakeholders to see and use the features allows them to see if their ideas are taking shape as expected.
Does the design work in practice?
Does the frontend feel good to use?
Does the addition of these features make the application better?
Getting early answers to these crucial questions and adapting as soon as possible is paramount to a project's success.

## Pyramids

How does all of this relate to building pyramids?
Well, while reading Daniel Hooper's amazing [blog post](https://danielchasehooper.com/posts/shapeup/) about building a 3D modeling tool, he explained this concept with a very simple image.
In his own words:

> I tried to work in such a way that I always had a working 3D modeler, and progressively improved it as time allowed.
> I think about it like building a pyramid.
> If you build layer by layer, you don’t have a pyramid until the very end.
> On the other hand you can build it so that stopping at any step is a complete pyramid.

<img src="/images/20240512/pyramid.webp" alt="Two ways to build a pyramid" style="width: 100%" />

On the left, you build the entire pyramid from the ground up and have nothing to present until the very end.
On the right, you first build the MVP (minimum viable pyramid) and then iteratively enhance it until all the work is complete.
The critical invariant with the latter approach is that you always have a tangible (albeit smaller) pyramid to inspect and appreciate.

This analogy of comparing project management to pyramid construction has been used a [few](https://www.informationweek.com/it-leadership/two-ways-to-build-a-pyramid) [times](https://blog.platypus.solutions/how-not-to-build-a-pyramid-314f8838a61f) before.
Additionally, it corresponds quite strongly to the third principle of the [Agile Manifesto](https://www.agilealliance.org/agile101/12-principles-behind-the-agile-manifesto/):

> Deliver working software frequently, from a couple of weeks to a couple of months, with a preference to the shorter timescale.

## Conclusion

This line of thinking impacted how I planned my first multi-developer, multi-month software development project.
Despite trying to plan the entire project upfront (waterfall style), I adjusted my strategy to prioritize delivering an MVP of the new features as soon as possible.
This put the project in a position where early input could be collected from the product and design teams.
Over the next few weeks, we will be iteratively enhancing the frontend with additional functionality.
These changes will be usable and testable by the project's stakeholders at each step along the way.
