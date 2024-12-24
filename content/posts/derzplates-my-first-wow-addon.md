---
date: 2024-12-22
title: "DerzPlates: My First WoW Addon"
slug: "derzplates-my-first-wow-addon"
---

I've played World of Warcraft off and on since 2007.
I spent months as a true n00b but eventually found myself a guild and a raiding spot as a balance druid (sometimes known as the "OOMkin").
As I learned more about the game, I quickly discovered addons: additional, community-built software that runs inside of the game client.
Addons can do so many things: track DPS (damage per second) numbers, show what loot a boss drops, or completely redesign the user interface.

Eventually, I went to college for Software Engineering.
Playing WoW since then, I've always wondered: how are addons written?
What language are they written in?
What API does the client expose to addon developers?
Despite being unanswered for so many years, I finally flipped that stone over and started reading.
Then, a few weeks ago, I wrote and published my very first WoW addon!
This post talks about what the addon does, how it works, and the struggles I went through during development.

# DerzPlates

I'm rolling a bear druid (a tanking spec) for the latest round of WoW Classic (the [20th Anniversary Edition](http://worldofwarcraft.blizzard.com/en-us/news/24156594)).
As a tank, my goal to get all of the enemies to attack me while the damage dealers kill them and the healers keep everyone alive.
Unfortunately, the default UI in classic doesn't provide a way to indicate which enemies are attacking you.
The addon I wrote is called [DerzPlates](https://github.com/theandrew168/derzplates) and it has single, simple goal: **make it visually clear (via nameplate color) which enemy units are attacking me**.

I'm aware that some addons already exist for this ([ThreatPlates](https://www.curseforge.com/wow/addons/tidy-plates-threat-plates) is a very popular one), but I feel like those are too heavy-handed and make too many changes to the default UI.
I _only_ want to change the plate's color (to magenta because purple is cool) and nothing else.
Here's a screenshot of the addon in action.
Bonus points if you notice the Wicked reference!

![Screenshot of modified nameplate color](/images/20241222/screenshot.webp)

# Addon Structure

As it turns out, WoW addons are written in [Lua](https://www.lua.org/)!
They live as self-contained directories within the game's `Interface/AddOns` folder.
To be identified and executed by the game client, the folder needs at least two files: a table of contents (`.toc`) file and a Lua (`.lua`) source code file.
That's it!

## Table of Contents

The `.toc` file contains basic metadata about the addon: its name, description, version, and which version of WoW client (interface) it was built for.
It then lists out any source files that the client should execute.

```ini
## Title: DerzPlates
## Notes: The simplest threat plates addon on the market.
## Author: theandrew168
## Version: 0.0.2
## Interface: 11505

main.lua
```

## Lua Sources

The simplest addon could be `print("Hello World!")` which would print "Hello World!" to the chat window whenever the UI loads.
DerzPlates does a bit more than that but I won't dump all of the code here.
Instead, I'll just a give you a sample of what addon code looks like.
Here is the main decision logic: it decides whether a given nameplate should use my special color or fallback to its default color.

```lua
-- This is a custom implementation of the CompactUnitFrame_UpdateHealthColorOverride
-- function. Once setup, this gets called whenever CompactUnitFrame_UpdateHealthColor
-- gets called in order to alter the nameplate's color. If our override function returns
-- false, then the default coloring code with run. If it returns true, then our color
-- will be used instead.
local function UpdateHealthColorOverride(self)
	-- Prefer self.displayedUnit but fallback to self.unit if necessary.
	local unit = self.displayedUnit or self.unit

	-- Don't change the nameplate color for players.
	if UnitIsPlayer(unit) then
		return false
	end

	-- Don't change the nameplate color for units you cannot attack.
	if not UnitCanAttack("player", unit) then
		return false
	end

	-- UnitThreatSituation returns 3 if the player has the highest threat
	-- and is the primary target of a given unit.
	local isTanking = 3

	-- Check if the player is currently tanking the nameplate's unit.
	local status = UnitThreatSituation("player", unit)
	if status == isTanking then
		-- If they ARE tanking this unit, set the plate's color to magenta.
		self.healthBar:SetStatusBarColor(1, 0, 1)
		-- Overwrite saved color so CompactUnitFrame_UpdateHealthColor() can restore the default later.
		self.healthBar.r, self.healthBar.g, self.healthBar.b = 1, 0, 1
		-- Signal CompactUnitFrame_UpdateHealthColor() that we set the color ourselves.
		return true
	end

	-- If we got here, we want CompactUnitFrame_UpdateHealthColor() to apply its default color.
	return false
end
```

In my opinion, Lua is pretty readable!
I haven't written it for a while but I did use it for a few projects back in college.
I also did a [Love2D](https://love2d.org/) (a Lua 2D game development framework) [game jam](https://github.com/theandrew168/lovejam2023) last year so my Lua skills aren't as dull as they otherwise would've been.

# Development Struggles

At first, I thought that writing the addon would be quite trivial:

1. Listen for a few specific events: nameplate created and threat situation updated.
2. When these events occur, loop through every nameplate visible on the screen.
3. If the nameplate belongs to a unit that is attacking me, change its color!

This did get me to a mostly working (but buggy) iteration of the addon.
It wasn't perfect, however, and I quickly hit a snag.
You see, in the WoW client, nameplates are "pooled".
This means that instead of creating and destroying nameplates objects for each unit, they are reused from a pre-allocated pool (similar to database connections).
A downside of this is that once you change a plate's color, it could get reused for a different unit and end up showing the wrong color.
After spending a few hours trying to find any information regarding this behavior (documentation for the WoW Lua API can be a bit sparse), I was close to giving up.

In an act of desperation, I summarized and posted my struggles onto the [WoW Interface forums](https://www.wowinterface.com/forums/showthread.php?p=344701).
Amazingly, this community showed up a in huge way!
Super big thanks to SDPhantom for providing what might be the best and most helpful technical answer I've ever received from a stranger on the internet.
Their advice and guidance led me right to the behavior I wanted.
They even pointed out how my initial usage of the event system was prone to race conditions and presented a better alternative.
Thanks to SDPhantom, I was able to release a much more stable v0.0.2.

This did make me wonder, though: how do developers learn this stuff?
How did SDPhantom instantly know what was wrong with my code and how to fix it?
Surely experience plays a role but what about examples and documentation?
From what I understand, there are two primary resources:

1. The [Warcraft Wiki](https://warcraft.wiki.gg/wiki/Warcraft_Wiki:Interface_customization)
2. The Blizzard UI [source code](https://github.com/tomrus88/BlizzardInterfaceCode/tree/classic)

When writing new addons, programmers will often scour the existing code to get a sense of how it works.
Then, they can apply this knowledge to identify where and how their logic should hook in and override default behavior.

# Automation

The code itself is hosted on GitHub.
I wrote a quick [GitHub Action](https://github.com/theandrew168/derzplates/blob/main/.github/workflows/release.yml) to zip the addon code and create a release whenever a new tag is added.
I don't write new actions very often so I always have to revisit the docs and my previous workflows to relearn.
The heavy lifting of this particular action is handled by the [TheDoctor0/zip-release](https://github.com/TheDoctor0/zip-release) and [ncipollo/release-action](https://github.com/ncipollo/release-action) actions.
Being able to compose these simple and specific building blocks together to create specialized workflows is really awesome!

If you want to download the addon and try it for yourself, just grab the latest `DerzFiles.zip` file from the [releases page](https://github.com/theandrew168/derzplates/releases) and unzip it into your addons folder.
Let me know if you end up giving it a shot!

![Screenshot of a recent release](/images/20241222/release.webp)

# Conclusion

Overall, this has been a short, fun, and only mildly frustrating experience.
If it weren't for the nameplate pooling behavior, I would've had this completely done and sorted within an hour.
However, troubleshooting the pooling cost me quite a few hours.
Thankfully, SDPhantom and the WoW Interface forum community showed up to help me out.
With their support, I was able to get everything working before the project turned into a stressful time sink.

I might iterate on DerzPlates in the future.
For v0.0.3, I'd like to add a super simple settings page where players can tweak the "you have threat" color.
While I personally love the magenta, other players might prefer something else.
I might even write _more_ addons if I find other preferences or behaviors that the default UI lacks.
I suppose time will tell.

Thanks for reading and may the winds guide you!
