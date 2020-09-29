---
date: 2020-09-28
title: "Bloggulus: A Personalized Blog Aggregator"
slug: "bloggulus-personalized-blog-aggregator"
tags: ["python", "flask", "peewee", "sqlite"]
---
what it is  
why i built it  
replaced feeder  

# The Stack
python  
flask  
peewee  
sqlite  
tailwind CSS  

# Deployment
NGINX reverse proxy  
LE via certbot + cron  
GitHub actions  
zipapp works great!  
structure the app as a CLI (gunicorn, syncfeeds, etc)  

# SQLite Value
Charles Leifer's blog and peewee  
FTS is amazing  
entire DB is a single file ~20MB  
sync process runs hourly  
