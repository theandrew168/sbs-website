---
date: 2026-08-23
title: "A Positive Experience with GPT-5.6 Luna"
slug: "a-positive-experience-with-gpt-5.6-luna"
---

The other day, I was looking at [Virginia Tech's bicycle helmet ratings](https://www.helmet.beam.vt.edu/bicycle-helmet-ratings.html).
The data here is quite interesting: a breakdown of over 300 helmets including info about price, impact score (lower is better), and an overall star rating.
They recommend picking a helmet that is has a 4 or 5 star rating.
Clicking on a individual helmet gives you some extra details: how its price ranks among **all** helmets and how its score compares to the average of **only the helmets with the same number of stars**.

Naturally, I was curious to see the distribution of helmet scores across all helmets and not just each helmet against its star group's average.
The first step was to get the data into a more malleable form: ideally a CSV.
Firefox has a nice option to download a web page as raw text, so I started there.
The [raw text](https://github.com/theandrew168/helmet-ratings/blob/main/ratings.txt) was consistent enough that I knew I could write a basic parsing script to convert it to [a CSV](https://github.com/theandrew168/helmet-ratings/blob/main/ratings.csv).
If I'd had to write the code by hand, it probably would have taken 15-20 minutes.

## The Prompt

But, I didn't have 15-20 minutes.
We were about to leave and I only had 1-2 minutes.
Since I'd recently been experimenting with OpenAI's latest generation of models (GPT-5.6 Sol, Terra, and Luna), I decided to see if this was an appropriately sized / scoped task for Luna.

On thinking level "medium", I gave it the following prompt:
> The @ratings.txt file contains a list of ~300 bicycling helmets and their safety ratings. I'd like to convert it to a CSV with the following columns: name, score, stars, and price. Can you write a python script to do this conversion? Just call it main.py. It shouldn't need any deps: just the stdlib and csv module, probably.

It one-shot the task perfectly in roughly 30 seconds.
Due to my relative unfamiliarity with OpenAI's models, I was surprised.
Not just with the quality, but also with the speed and the price (it cost $0.007 total: just under a cent).
Others with more hands-on experience with OpenAI's models (and AI tech in general) probably expected this outcome, but it exceeded my expectations.
I had the data parsed, converted, and graphed within that tight 1-2 minute window.

## The Takeaway

At work, I have a Claude Code Max subscription.
Because the tokens are essentially unlimited, I feel like I tend to overlook the nuances of _how_ I'm using the LLM.
I _think_ I'm typically using Opus 4.8... or maybe Opus 5?
And I _think_ I'm using the "high" thinking level?
The point is I don't really know and I don't really pay attention to it.
I don't _think_ this is a particularly good practice, though.
Surely, I can do better.

After hearing good things about [Pi](https://pi.dev/), [Plannotator](https://plannotator.ai/), and [OpenAI's latest models](https://openai.com/index/gpt-5-6/) from the [Next Token podcast](https://www.youtube.com/@NextTokenShow) (which I highly recommend, by the way), I decided I wanted to dig a bit deeper and try out some new tools.
I liked the sound of Pi because it embraces minimalism: few features, only a handlful of tools, and a tiny system prompt.
Claude Code's system prompt is pretty large and tends to change version-by-version, from what I hear.
Plannotator is great for giving direct, line-specific feedback on plans.
It even supports line-based local code feedback with a diff-based UI.
Very cool!

And lastly, the new GPT 5.6 models.
A few days prior, I used Sol to add a new feature to Mowhen.
This was similar to how I approach things at work: using a sledgehammer for everything.
While the model did very well, it also burned like $2-3 of credits in only a few minutes.
For context, I only bought $50 worth of OpenAI credits and I'm hoping they can last a while.
I had a sudden realization: working with usage-based API credits forces me to care a lot more about how they are being spent.

So I set off on a new mission: finding out which models are the "best fit" for a given problem.
GPT-5.6 has three models and each model has a range of thinking levels: from "minimal" to "max".
How good is Luna, the little guy?
Obviously it destroys small tasks like this data conversion script but what else is it capable of?
I don't know, but finding out is the point of this.
**It seems like starting small, giving the model increasing complex tasks, and tracking when it starts to fall short is the way to go.**
I'll continue to explore this in the coming weeks.

Thanks for reading!