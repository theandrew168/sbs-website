---
date: 2026-08-23
title: "A Positive Experience with GPT-5.6 Luna"
slug: "a-positive-experience-with-gpt-5.6-luna"
draft: true
---

Today, I was looking into [Virginia Tech's bicycle helmet ratings](https://www.helmet.beam.vt.edu/bicycle-helmet-ratings.html).
The data here is quite interesting: a breakdown of over 300 helmets complete with price, impact score (lower is better), and an overall star rating.
They recommend picking a helmet that is has a 4 or 5 star rating.
Clicking on a individual helmet gives you some extra details: how its price ranks among **all** helmets and how its score compares to the average of **only the helmets with the same number of stars**.

Naturally, I was curious to see the distribution of helmet scores across all helmets and not just each helmet against its star group's average.
The first step was to get the data into a more malleable form: ideally a CSV.
Firefox has a nice option to download a web page as raw text, so I started there.
The [raw text](https://github.com/theandrew168/helmet-ratings/blob/main/ratings.txt) was consistent enough that I knew I could write a basic parsing script to convert it to a CSV.
If I'd had to write the code by hand, it'd probably have taken 15-20 minutes.

But, I didn't have 15-20 minutes.
We were about to leave and I only had 1-2 minutes.
Since I'd recently been experimenting with OpenAI's latest generation of models (GPT-5.6 Sol, Terra, and Luna), I decided to see if this was an appropriately sized / scoped task for Luna.

On thinking level medium, I gave it the following prompt:
> The @ratings.txt file contains a list of ~300 bicycling helmets and their safety ratings. I'd like to convert it to a CSV with the following columns: name, score, stars, and price. Can you write a python script to do this conversion? Just call it main.py. It shouldn't need any deps: just the stdlib and csv module, probably.

In roughly 30 seconds, it was done!
Due to my relative unfamiliarity with OpenAI's models, I was surprised.
Not just with the quality, but also with the speed and the price (it cost me $0.007 total, just under a cent).
Others with more hands-on experience with OpenAI's models (and AI tech in general) probably expected this outcome, but it definitely went above my expectations.

At work, I have a Claude Code Max subscription (paid for by the company).
Because the tokens are essentially unlimited and because the money isn't coming out of my own pocket, I feel like I tend to overlook the details of what I'm doing.
I _think_ I'm typically using Opus 4.8 or maybe Opus 5?
And I _think_ I'm using the "high" thought level?
The point is I don't really know and I don't really pay attention to it.
I don't think is particularly good practice, though.

After hearing good things about [Pi](https://pi.dev/), [Plannotator](https://plannotator.ai/), and [OpenAI's latest models](https://openai.com/index/gpt-5-6/) from the [Next Token podcast](https://www.youtube.com/@NextTokenShow) (which I highly recommend, by the way), I decided I wanted to dig a bit deeper and try out some new tools.
I liked the sound of Pi because it embraces minimalism: few features, only a handlful of tools, and a tiny system prompt.
Claude Code's system prompt is pretty large and tends to change version-by-version, from what I hear.
Plannotator is great for giving direct, line-specific feedback on plans.
It even supports line-based local code feedback (with a diff-based UI).
Very cool!

And lastly, the new GPT 5.6 models.
I started with Sol (similar to how I approach things at work) for some basic web dev tasks and it did very well, but also burned like $2-3 of credits in just a few minutes (I only bought $50 worth of OpenAI credits and I want em to last).
So began a new mission: finding out which models are the "best fit" for a given problem.
I want to optimize my workflows and be more efficient.
Basic CSS work? Definitely Luna.
Planning small features? Maybe Luna?
I don't know, but finding out is the point of this.
I want to develop a better intuition for each model's capabilities and the impact of different "thinking levels".

This was just one positive experience but I'm optimistic that more will come in the future!
Thanks for reading.