---
date: 2026-08-22
title: "Introducing Mowhen"
slug: "introducing-mowhen"
tags: ["Mowhen"]
draft: true
---

Mowhen is launched! What does it do? Why is it useful? What is on the roadmap?

## Overview

Mowhen is a service helps you track when to mow.
It isn't just a scheduling app, though: it uses weather data (based on the lawn's location) and a bunch of math to calculate how much lawn growth has occurred.
When a lawn is ready to mow it can notify you via email.

## Notifications

There are multiple notification types (which can be individually enabled / disabled):
* Mow Today - Your lawn is ready to be mowed today based on growth patterns and weather conditions.
* Mow Tomorrow - Get a heads-up the day before so you can plan your schedule accordingly.
* Rain Check Alert - Mowing is due in two days, but rain is forecasted. Consider mowing early to avoid wet grass.
* Schedule Updated - Your next mowing date has been determined based on extended weather forecasts.

## Lawn Settings

Lawns can be configured based on its individual settings:
* Type (cool vs warm season)
* Sunlight (full vs partial)
* Irrigation (automated vs none)
* Fertilization (regular vs none)

## Tech Stack

The website is built using SvelteKit and TypeScript.
CSS is vanilla and leverages SvelteKit's automatic scoping and isolation.
PostgreSQL is used for the database (I'm a big fan).
The app is deployed on DigitalOcean (managed by Terraform and Ansible).
Architecture follows what I [recently wrote about](/posts/architecture-md/): basics impls of DDD and CQRS.