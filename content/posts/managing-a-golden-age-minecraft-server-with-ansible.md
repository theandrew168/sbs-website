---
date: 2024-01-30
title: "Managing a Golden Age Minecraft Server with Ansible"
slug: "managing-a-golden-age-minecraft-server-with-ansible"
tags: ["ansible", "minecraft"]
draft: true
---

Lately, I've noticed a growing interest in what people like to call "Golden Age Minecraft".
This term typically refers to versions of Minecraft prior to 1.8.0: The Adventure Update.
This updated added a ton of new features and, for better or for worse, added well-defined goals to the game.
Prior to this update, Minecraft was truly a sandbox.
The only things to do were explore, build stuff, and find diamonds.
It was a simple loop and generated hundreds of hours of fun for my friends and I back in high school.

# Simpler Times

Back in the day, I ran a Minecraft server from the utility room in my parent's house (so that it could be plugged into the router).
It ran a terrible Asus Netbook with barely any specs.
Knowing next to nothing about technology, I struggled to set this up.
This was my first foray into software and technology.
The networking was especially difficult at the time.
I had no idea what port was or how to "foward" it.
I also remember telling my friends (who lived across town): "I've got it setup! Try connecting to 192.168.1.10".
This gives me a laugh in hindsight, as I'm sure it does for any other tech professional.
The joke here is that the `192.168.0.0/16` address space is reserved for local networks: outside connections can't connect to this address space on a different network (not by default, at least).

Despite all the fuss, I eventually got the server up and running and we had a blast.
We used the `gargamel` seed (of course) and spent months transforming that beautiful canyon into an epic base.
We had a storage house, a monster arena, and various redstone shenanigans (I also had to learn what a XOR gate was).
We made a giant lighthouse, multiple outposts, and underground rail systems to connect everything.
Everything we did was motivated by one thing: fun.
There was no dragon to kill, no End to explore, no Elytra to obtain.
Just mine, craft, and have a good time.

# Professional Experience

These days, I know orders of magnitude more about technology.
I've had prior jobs that involved system admin work and a nice tool for automation called [Ansible]().
This post outlines using Ansible to manage Minecraft servers "in style".
