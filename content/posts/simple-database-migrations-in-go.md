---
date: 2021-01-12
title: "Simple Database Migrations in Go"
slug: "simple-database-migrations-in-go"
tags: ["go", "database"]
draft: true
---
Database migrations help solve the problem of track schema changes throughout the development of an application.  
New tables, altered columns, etc.  
All of these changes need to be developed and tested before being deployed.  
Manually applying these is really not a good option.  

# What needs to happen?
I was going to use a library for this but with a newfound sense of problem exploration from pymkv, I thought: how hard could this really be?  
In theory, the idea is simple: find a way to create and track sequential database changes.  
We have .sql files in the migrations/ dir and a database that maybe already have none, some, or all applied.  
Where and how do we track this and what is the process for making sure everything stays in sync?

# The basic gameplan
We've got a database, so lets just store this info in a table  
Need to store which migrations have been applied  
Name with 0001 style to preserve order when sorted  

Make sure db connection works  
Create migration table if it doesn't exist  
Get set of currently applied migrations  
Glob out the migrations in migrations/ dir  
For each of these, check if in the applied set, add to missing slice if not  
Sort the missing slice to preserve order  

For each of these files:  
Read the file  
Exec the SQL  
Insert the file name into the migration table  
