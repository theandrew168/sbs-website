---
date: 2024-07-21
title: "Simple Server Swap Space"
slug: "simple-server-swap-space"
tags: ["Ansible"]
---

Since the recent [Bloggulus outage](/posts/bloggulus-outage-postmortem/), I've been keeping a close eye on things.
While the server has mostly been stable, I still noticed the occasional OOM kill after creating backups via [pg2s3](https://github.com/theandrew168/pg2s3).
Here is an example from `journalctl -u pg2s3` logs (notice that this was happening nearly a month after the most-outage memory increase):

```
Jul 01 09:00:03 bloggulus pg2s3[32400]: created bloggulus_2024-07-01T09:00:00Z.backup.age
Jul 01 09:00:04 bloggulus pg2s3[32400]: deleted bloggulus_2024-05-04T09:00:00Z.backup.age
Jul 02 09:00:05 bloggulus systemd[1]: pg2s3.service: A process of this unit has been killed by the OOM killer.
Jul 02 09:00:05 bloggulus systemd[1]: pg2s3.service: Main process exited, code=killed, status=9/KILL
Jul 02 09:00:05 bloggulus systemd[1]: pg2s3.service: Failed with result 'oom-kill'.
Jul 02 09:00:05 bloggulus systemd[1]: pg2s3.service: Consumed 9.452s CPU time.
```

Okay, so it seems like the 1GB of RAM isn't quite enough when backups are taking place.
The server works just fine under normal operation, however, but backups push it over the edge.
If only there was a way to "download more RAM" and give the server a bit more breathing room...

# Swap Space

Enter the [swap space](https://wiki.archlinux.org/title/Swap)!
This is a Linux concept for giving servers additional memory capabilities without increasing the amount of physical RAM installed.
It works by using a regular file (on the filesystem) for memory overflow.
It isn't as fast as regular RAM, but it is better than having processes getting OOM killed!
From [All about Linux swap space](https://www.linux.com/news/all-about-linux-swap-space/):

> Linux divides its physical RAM (random access memory) into chunks of memory called pages.
> Swapping is the process whereby a page of memory is copied to the preconfigured space on the hard disk, called swap space, to free up that page of memory.
> The combined sizes of the physical memory and the swap space is the amount of virtual memory available.

I didn't realize that my Digital Ocean droplets don't have swap configured and enabled by default.
Thankfully, it is quite easy to set up and will hopefully help eliminate those pesky remaining OOM kills.
Digital Ocean has a [great guide](https://www.digitalocean.com/community/tutorials/how-to-add-swap-space-on-ubuntu-22-04) for setting up and configuring swap space on an Ubuntu server.
Initially, I followed these steps manually.
Then, once I got things working, I decided to "lock it in" via my Ansible automation.

# Ansible Tasks

The automation was quite simple: it only takes five tasks!
I used Jeff Geerling's awesome [ansible-role-swap](https://github.com/geerlingguy/ansible-role-swap/tree/master) for inspiration and guidance.
In short, these tasks create, initialize, and enable a 1GB swap file on the root filesystem (at `/swapfile`).
It also adds an entry to [fstab](https://wiki.archlinux.org/title/Fstab) so that the swap space is automatically enabled on subsequent restarts.

Let's take a look:

```yml
- name: Create swapfile
  command:
    cmd: fallocate -l 1G /swapfile
    creates: /swapfile
  register: create_swapfile
  become: yes
  become_user: root

- name: Set swapfile permissions
  file:
    path: /swapfile
    mode: "0600"
  become: yes
  become_user: root

- name: Initialize swapfile
  command:
    cmd: mkswap /swapfile
  when: create_swapfile is changed
  become: yes
  become_user: root

- name: Enable swapfile
  command:
    cmd: swapon /swapfile
  when: create_swapfile is changed
  become: yes
  become_user: root

- name: Add swapfile to fstab
  mount:
    name: none
    src: /swapfile
    fstype: swap
    opts: sw
    state: present
  become: yes
  become_user: root
```

# Bonus Tuning

The Digital Ocean guide also [details](https://www.digitalocean.com/community/tutorials/how-to-add-swap-space-on-ubuntu-22-04#step-6-tuning-your-swap-settings) a few "bonus" settings that can help a server manage its swap more efficiently.

## swappiness

This [setting](https://docs.kernel.org/admin-guide/sysctl/vm.html#swappiness) controls how eager a system is to boot data out of main memory and into the swapfile.
Since we only want the server to use the swap space when absolutely necessary, we adjust this setting to a low value.
The default value is `60`.

```yml
- name: Configure swappiness
  sysctl:
    name: vm.swappiness
    value: 10
  become: yes
  become_user: root
```

## vfs_cache_pressure

This [setting](https://docs.kernel.org/admin-guide/sysctl/vm.html#vfs-cache-pressure) controls how quickly the server releases directory and inode (file) information from the cache.
Lowering this value causes filesystem data (which can be expensive to retrieve) to remain in the cache for longer periods of time.
The default value is `100`.

```yml
- name: Configure vfs_cache_pressure
  sysctl:
    name: vm.vfs_cache_pressure
    value: 50
  become: yes
  become_user: root
```

# Conclusion

You can view all of these tasks together in my [devops repo](https://github.com/theandrew168/devops/blob/main/roles/server/tasks/swap.yml).
I also decided that my server role had gotten a bit messy and decided to split it into separate files using Ansible's [include_tasks](https://docs.ansible.com/ansible/latest/collections/ansible/builtin/include_tasks_module.html) directive.
Now I've got a clean separation of tasks which is much easier to navigate and understand.

Since I added swap space to the Bloggulus server, I haven't seen a single OOM kill.
Things have been running smoothly for weeks and only a small portion of the swap is being used:

```
               total        used        free      shared  buff/cache   available
Mem:           957Mi       226Mi        77Mi        92Mi       653Mi       481Mi
Swap:          1.0Gi        41Mi       982Mi
```

Let's hope things stay that way! Thanks for reading.
