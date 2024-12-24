---
date: 2024-12-22
title: "DerzPlates: My First WoW Addon"
slug: "derzplates-my-first-wow-addon"
draft: true
---

I've played World of Warcraft off and on since around 2007.
I spent months as a true n00b but eventually found myself a guild and a raiding spot as a balance druid (sometimes known as the "OOMkin").
As I learned more about the game, I quickly discovered addons: extra, community-built software that runs inside of the game client.
Addons can do so many things: track DPS (damage per second) numbers, show what loot a boss drops, or completely redesign the entire user interface.

Eventually, I went to college for Software Engineering.
Playing WoW since then, I've wondered: how are addons written?
What language are they written in?
What API does the client expose to addon developers?
Despite being unanswered for so many years, I finally flipped that stone over and started reading.
Then, a few weeks ago, I wrote and published my very first WoW addon!

# DerzPlates

I'm rolling a bear druid (a tanking spec) for the latest round of WoW Classic (the [20th Anniversary Edition](http://worldofwarcraft.blizzard.com/en-us/news/24156594)) so this is important.
As a tank, my goal to get all of the enemies to attack me while the damage dealers kill them and the healers keep everyone alive.
Unfortunately, the default UI doesn't provide a way to indicate which enemies are attacking you.
The addon I wrote is called [DerzPlates](https://github.com/theandrew168/derzplates) and it has single, simple goal: **make it visually clear (via nameplate color) which enemy units are attacking me**.

I'm aware that some addons already exist for this ([ThreatPlates](https://www.curseforge.com/wow/addons/tidy-plates-threat-plates) is a very popular one)
But I feel like those are too heavy-handed and make too many changes to the default UI.
I _only_ want to change the plate's color (to magenta because purple is cool) and nothing else.

![Screenshot of modified nameplate color](/images/20241222/screenshot.webp)

# Addon Structure

As it turns out, WoW addons are written in [Lua](https://www.lua.org/)!
They live as self-contained directories within the game's `Interface/AddOns` folder.
To be identified and executed by the game client, the folder needs at least two files: a table of contents (`.toc`) file and a Lua (`.lua`) source code file.
That's it!

## Table of Contents

The `.toc` file contains basic metadata about the addon: its name, description, version, and which version of WoW it was built for.
It then lists out any source files that the client should execute.

```ini
## Interface: 11505
## Title: DerzPlates
## Notes: The simplest threat plates addon on the market.
## Author: theandrew168
## Version: 0.0.2

main.lua
```

## Lua Sources

The simplest addon could be `print("Hello World!")` which would print "Hello World!" to the chat window whenver the UI loads.
DerzPlates does a bit more than that but I won't list all of its code here.
Instead, I'll just a give you a sample of what addon code looks like.
Here the main decision logic from my addon: it decides whether a given nameplate should use my special color or fallback to its default color.

```lua
-- This is a custom implementation of the CompactUnitFrame_UpdateHealthColorOverride
-- function. Once setup, this gets called whenever CompactUnitFrame_UpdateHealthColor
-- gets called in order to alter the nameplate's color. If our override function returns
-- false, then the default coloring code with run. If it returns true, then our color
-- will be used instead.
--
-- Reference:
-- https://github.com/Gethe/wow-ui-source/blob/e696432cf6c1dcf18036590b64b11c975d8f9fb9/Interface/AddOns/Blizzard_UnitFrame/Classic/CompactUnitFrame.lua#L394-L396
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

In my opinion, I think Lua is pretty readable!
I haven't written in it for a while but I did create a few projects with it back in college.
I also did a Love2D (a Lua game development framework) [game jam](https://github.com/theandrew168/lovejam2023) last year so my Lua skills aren't as dull as they otherwise would've been.

# Development Struggles

At first, I thought the addon would be quite trivial:

1. Listen for a few specific events: nameplate created and threat situation updated.
2. When these events occur, loop through every nameplate on screen.
3. If the nameplate belongs to a unit that is attacking me, change its color!

This did get me to a mostly working (but buggy) iteration of the addon.
It wasn't perfect, however, and I quickly hit a snag.
You see, in the WoW client, nameplates are "pooled".
This means that instead of creating and destroying nameplates objects for each unit, they are reused from a pre-allocted pool (similar to database connections).
A downside of this is that once you change a plate's color, it could get reused for a different unit and end up showing the wrong color.
After spending a few hours trying to find any information regarding this behavior (documentation for the WoW Lua API can be a bit sparse), I was close to giving up.

In an act of desperation, I summarized and posted my struggles onto the [WoW Interface forums](https://www.wowinterface.com/forums/showthread.php?p=344701).
Amazingly, this community showed up a in huge way!
Super big thanks to SDPhantom for providing with what might be the best and most helpful technical answer I've ever received from a stranger on the internet.
Their advice and guidance let my right to the behavior I was after.
They even pointed out how my initial usage of the event system was prone to race conditions and presented a better alternative.
Thanks to SDPhantom, I released version 0.0.2.

How do devs learn this?
From what I understand, there are two primary resources:

1. The [Warcraft Wiki](https://warcraft.wiki.gg/wiki/Warcraft_Wiki:Interface_customization)
2. The Blizzard UI [source code](https://github.com/tomrus88/BlizzardInterfaceCode/tree/classic)

When writing new addons, devs will often scour the existing code to get a sense of how it works.
Then, they can apply this knowledge to identify where and how their logic should hook in and override default behavior.

# Automation

The code itself is hosted on GitHub and pretty short.
I wrote a quick [GitHub Action](https://github.com/theandrew168/derzplates/blob/main/.github/workflows/release.yml) to zip the addon whenever a new tag is added.
I hadn't written a new action in quite a while so I always have to revisit the docs and my previous workflows to relearn.
The heavy lifting of this particular action is handled by the [TheDoctor0/zip-release](https://github.com/TheDoctor0/zip-release) and [ncipollo/release-action](https://github.com/ncipollo/release-action) actions.
Being able to compose these simple and specific building blocks together to create specialized workflows is really awesome!

If you want to download the addon and try it for yourself, just grab the latest `DerzFiles.zip` file from the [releases page](https://github.com/theandrew168/derzplates/releases) and unzip it into your addons folder!

![Screenshot of a recent release](/images/20241222/release.webp)

# Conclusion

Overall, this has been a short, fun, and only mildly frustrating experience.
If it weren't for the plate pooling, I would've had this completely done and sorted within an hour.
However, troubleshooting the pooling has cost me many hours!
Still fun, though.
I've written Lua in the past so that isn't completely new to me.

I might iterate on this in the future, still.
For version 0.0.3, I'd like to add a super simple settings page where players can tweak the "you have threat" color.
While I personally love the magenta, other players might prefer something else.
I might even write _more_ addons if I find other preferences or behaviors that the default UI lacks.
I suppose time will tell.

Thanks for reading!
