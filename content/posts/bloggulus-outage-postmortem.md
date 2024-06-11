---
date: 2024-06-09
title: "Bloggulus Outage Postmortem"
slug: "bloggulus-outage-postmortem"
draft: true
---

Bloggulus was down for ~7 hours.
What happened?
How did I fix it?
Why did it take me so long to notice?
How long was the locale corrupted before I noticed (til PG or the server rebooted).
Which upgrade was most likely to break it?
What does `locale-gen` really do?
Calls `localedef` and gens `/usr/lib/locale/locale-archive` binary file (database of locale data).

Tried to do my own `localedef`... OOM-Killed!
Saw some other oom kills for PG and Redis...
Theory: an upgrade came in (something `l10n`) that required locales to be regenerated.
However, they failed to due to being OOM killed, which left the locale DB in a bad / empty state.

Found the smoking gun!
Apt logs from 5/31/2024 show the `locale-gen` OOM kill after upgrading `locales` package.
Server was now in a bad state and would break upon the next reboot (didn't happen until 6/4/2024).

## Bloggulus is Down

```
Jun 04 14:51:39 bloggulus systemd[1]: Starting bloggulus...
Jun 04 14:51:39 bloggulus bloggulus[37755]: 2024/06/04 14:51:39 ERROR failed to connect to `host=localhost user=bloggulus database=bloggulus`: dial error (dial tcp 127.0.0.1:543>
Jun 04 14:51:39 bloggulus systemd[1]: bloggulus.service: Main process exited, code=exited, status=1/FAILURE
Jun 04 14:51:39 bloggulus systemd[1]: bloggulus.service: Failed with result 'exit-code'.
Jun 04 14:51:39 bloggulus systemd[1]: Failed to start bloggulus.
Jun 04 14:51:45 bloggulus systemd[1]: bloggulus.service: Scheduled restart job, restart counter is at 4702.
Jun 04 14:51:45 bloggulus systemd[1]: Stopped bloggulus.
```

## PostgreSQL is Down

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

Take note of those `locale` warnings...
I first tried restarting the server (couldn't hurt, right) but that didn't make a difference: PostgreSQL was still failing to start.
Maybe the server's baseline configuration got messed up somehow?
Since I use Ansible to manage my servers, I'll re-run the playbook that handles the Bloggulus server.

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

![Downtime of roughly seven hours](/images/20240609/downtime.webp)
