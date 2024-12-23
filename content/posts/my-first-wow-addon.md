---
date: 2024-12-22
title: "My First WoW Addon"
slug: "my-first-wow-addon"
draft: true
---

Talk about DerzPlates, my first WoW addon.
Super simple but effective!
It has one goal: make it visually clear (via nameplate color) which units I'm actively tanking.
I'm rolling a bear druid for classic classic so this is important.

I know that some addons already exist for this: ThreatPlates and such.
But I feel like those are too heavy-handed and make too many changes to the default UI.
I _only_ want to change the plate's color (to magenta because purple is cool) and nothing else.

I quickly hit a snag: in the WoW client, nameplates are "pooled".
This means that instead of creating and destroying nameplates objects for each unit, they are reused from a pre-allocted pool (similar to database connections).
A downside of this is that once you change a plate's color, it could get reused for a different unit and end up showing the wrong color.
I'm still on the hunt for a way to "reset a nameplate's color" (I've posted on the WoW Addons Discord and the wow-interface forums).
Hopefully I can hear back because documentation for the WoW Lua API is a bit... sparse.

WoW Interface forums to the rescue!
Super big thanks to SDPhantom for providing with what might be the best and most helpful technical answer I've ever received from a stranger on the internet.
Their advice and guidance let my right to the behavior I was after.
They even pointed out how my initial usage of the event system was prone to race conditions and presented a better alternative.
Thanks for SDPhantom, I released version 0.0.2.

https://www.wowinterface.com/forums/showthread.php?p=344701

The code itself is hosted on GitHub and pretty short.
I wrote a quick action to zip the addon whenever a new tag is added.

Overall, this has been a small, fun, and frustrating project.
If it weren't for the plate pooling, I would've had this completely done and sorted within an hour.
However, troubleshooting the pooling has cost me many hours!
Still fun, though.
I've written Lua in the past so that isn't completely new to me.

I might iterate on this in the future, still
For version 0.0.3, I'd like to add a super simple settings page where players can tweak the "you have threat" color.
While I personally love the magenta, other players might prefer something else.

https://github.com/theandrew168/derzplates
