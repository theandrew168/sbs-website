---
date: 2024-06-09
title: "Bloggulus Outage Postmortem"
slug: "bloggulus-outage-postmortem"
draft: true
---

On June 4, 2024, [Bloggulus](https://bloggulus.com) was down for just under 7 hours.
Despite being offline for so long, the immediate fix only took me ~30 minutes to find and apply once I became aware of the issue.
This post details how I went about fixing the problem and then digs deeper into what actually happened.
I also discuss a few gaps in server monitoring and notification delivery.

![Downtime of roughly seven hours](/images/20240609/downtime.webp)

## Bloggulus is Down

I didn't actually become aware of Bloggulus being down until I visited the site myself and received a `502 Bad Gateway`.
This response is sent from [Caddy](https://caddyserver.com/) when it cannot communicate with the service it is proxying (the Bloggulus web server).
Despite my [UptimeRobot](https://uptimerobot.com/) monitor going red within 5 minutes of the outage, I have the notifications configured to be sent to my primary [shallowbrooksoftware.com](https://shallowbrooksoftware.com/) email address (managed via [Zoho](https://www.zoho.com/mail/)).
However, since I don't have this inbox configured on my phone, I had no way to know about it until checking my email on my personal laptop.

Anyhow, back to the action!
The first I did was SSH into the server and check the `systemd` logs:

```
Jun 04 14:51:39 bloggulus systemd[1]: Starting bloggulus...
Jun 04 14:51:39 bloggulus bloggulus[37755]: 2024/06/04 14:51:39 ERROR failed to connect to `host=localhost user=bloggulus database=bloggulus`: dial error (dial tcp 127.0.0.1:5432)
Jun 04 14:51:39 bloggulus systemd[1]: bloggulus.service: Main process exited, code=exited, status=1/FAILURE
Jun 04 14:51:39 bloggulus systemd[1]: bloggulus.service: Failed with result 'exit-code'.
Jun 04 14:51:39 bloggulus systemd[1]: Failed to start bloggulus.
Jun 04 14:51:45 bloggulus systemd[1]: bloggulus.service: Scheduled restart job, restart counter is at 4702.
Jun 04 14:51:45 bloggulus systemd[1]: Stopped bloggulus.
```

The error here is pretty clear: Bloggulus is unable to connect to the database running on the same machine.
Let's see what's going on with the database.

## PostgreSQL is Down

It is important to understand the architecture of Bloggulus: the Caddy reverse proxy, the Go-based web app, and the PostgreSQL database all run on a single [DigitalOcean](https://www.digitalocean.com/) droplet.
Caddy proxies traffic between the internet and Bloggulus while PostgreSQL stores all of the web app's persistent data (blogs, posts, etc).

Checking the database logs reveal that it is indeed dead as well:

```
Jun 04 14:51:36 bloggulus systemd[1]: Starting PostgreSQL Cluster 14-main...
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: perl: warning: Setting locale failed.
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: perl: warning: Please check that your locale settings:
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]:         LANGUAGE = (unset),
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]:         LC_ALL = (unset),
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]:         LANG = "en_US.UTF-8"
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]:     are supported and installed on your system.
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: perl: warning: Falling back to the standard locale ("C").
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: Error: /usr/lib/postgresql/14/bin/pg_ctl /usr/lib/postgresql/14/bin/pg_ctl start -D /mnt/bloggulus_db/data -l /var/log/postg>
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: 2024-06-04 14:51:36.657 UTC [37747] LOG:  invalid value for parameter "lc_messages": "en_US.UTF-8"
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: 2024-06-04 14:51:36.657 UTC [37747] LOG:  invalid value for parameter "lc_monetary": "en_US.UTF-8"
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: 2024-06-04 14:51:36.657 UTC [37747] LOG:  invalid value for parameter "lc_numeric": "en_US.UTF-8"
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: 2024-06-04 14:51:36.657 UTC [37747] LOG:  invalid value for parameter "lc_time": "en_US.UTF-8"
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: 2024-06-04 14:51:36.657 UTC [37747] FATAL:  configuration file "/etc/postgresql/14/main/postgresql.conf" contains errors
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: pg_ctl: could not start server
Jun 04 14:51:36 bloggulus postgresql@14-main[37742]: Examine the log output.
Jun 04 14:51:36 bloggulus systemd[1]: postgresql@14-main.service: Can't open PID file /run/postgresql/14-main.pid (yet?) after start: Operation not permitted
Jun 04 14:51:36 bloggulus systemd[1]: postgresql@14-main.service: Failed with result 'protocol'.
Jun 04 14:51:36 bloggulus systemd[1]: Failed to start PostgreSQL Cluster 14-main.
```

What's up with those `locale` warnings?
I don't recall ever seeing them before but my [Ansible automation](https://github.com/theandrew168/devops/blob/ae25c0a6333e8cfb3b6dd091b329342da7e545ba/roles/server/tasks/main.yml#L169-L187) generates and configures the `en_US.UTF-8` locale for all of my servers.
Thinking that this might be weird fluke or something, I rebooted the server.
Unfortunately, that didn't make a difference: PostgreSQL was still failing to start.

Maybe the server's baseline configuration got messed up somehow?
I decided re-run the [playbook](https://github.com/theandrew168/devops/blob/main/bloggulus.yml) that manages the Bloggulus server and all of its components.
Since all of my Ansible roles are idempotent, anything out of alignment should get straightened out.

## Running Ansible

Everything was green (unchanged) except for one task:

```yaml
- name: Ensure en_US.UTF-8 locale is available
  locale_gen:
    name: en_US.UTF-8
    state: present
  become: yes
  become_user: root
```

What does this mean?
The server's locale _setting_ didn't change, but the actual locale data itself did (the binary data that lives in `/usr/lib/locale/locale-archive`).
Very odd.
What happens if I try to perform a `locale-gen` myself?

## Manual Locale Gen

```
root@bloggulus:~# locale-gen "en_US.UTF-8"
Generating locales (this might take a while)...
  en_US.UTF-8.../usr/sbin/locale-gen: line 177: 15179 Killed                  localedef $no_archive -i $input -c -f $charset $locale_alias $locale
 done
Generation complete.
root@bloggulus:~# dmesg | grep -i kill
[115633.247115] Out of memory: Killed process 15179 (localedef) total-vm:146420kB, anon-rss:134440kB, file-rss:1756kB, shmem-rss:0kB, UID:0 pgtables:328kB oom_score_adj:0
root@bloggulus:~# locale-gen "en_US.UTF-8"
Generating locales (this might take a while)...
  en_US.UTF-8... done
Generation complete.
```

Whoa... It got OOM-Killed the first time but succeeded the second time.
The server must be _just_ nearing out of memory conditions.
To be fair, I _am_ running PostgreSQL, Caddy, and the Bloggulus web app all on a single machine with 512MB of RAM.

## The Smoking Gun

Well, if this happened to me, maybe it happened during an upgrade?
Maybe some packages require that `locales` get re-built?
Let's check the `apt` logs:

```
Log started: 2024-05-31  07:00:43
(Reading database ... 102795 files and directories currently installed.)
Preparing to unpack .../locales_2.35-0ubuntu3.8_all.deb ...
Unpacking locales (2.35-0ubuntu3.8) over (2.35-0ubuntu3.7) ...
Setting up locales (2.35-0ubuntu3.8) ...
Generating locales (this might take a while)...
  en_US.UTF-8.../usr/sbin/locale-gen: line 177: 183068 Killed                  localedef $no_archive -i $input -c -f $charset $locale_alias $locale
 done
Generation complete.
Processing triggers for man-db (2.10.2-1) ...
Log ended: 2024-05-31  07:00:48
```

There it is!
On May 31, after updating the `locales` package, running `locale-gen` was killed (likely by the OOM-Killer).
This means that the server was in an unstable
That seals the deal: this server needs more RAM!

## Moving Foward

Server needs more memory.
Bump from $4/month (512MB RAM) to $6/month (1GB RAM) droplet and update PG tuning.
Also change UptimeRobot email to something on my phone.
Stretch: re-add prom metrics and setup alerts for CPU/RAM/Storage and PG being down.
Revisit systemd / journald log sizes (50M might not be quite right).

https://gist.github.com/JPvRiel/b7c185833da32631fa6ce65b40836887
