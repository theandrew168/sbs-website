---
date: 2024-02-11
title: "Managing IP Allowlists with Ansible and UFW"
slug: "managing-ip-allowlists-with-ansible-and-ufw"
tags: ["Minecraft", "Ansible"]
---

I recently stood up a "Golden Age" minecraft server for my friends and me.
We've been having a great time so far building houses, cave bases, and clifftop observation decks.
Despite the initial fun, I've known that the server's security posture was a bit... shaky.
The server only allows a specific set of usernames to login, but, since the old authentication servers aren't online anymore, the server has no choice but to trust the client.

This means that a malicious / hacked client could spoof the username of a known player and immediately connect to the server.
I've been doing my best to keep the players' usernames secret but that is just "security by obscurity".
I outlined some options to increase the server's security in my [previous post](/posts/automating-a-golden-age-minecraft-server/) and one stood out as a clear winner: only allowing known IP addresses to connect.
Since the Minecraft server itself doesn't support this, I decided to move up a level and utilize the operating system's firewall: UFW.

# UFW

[UFW](https://help.ubuntu.com/community/UFW) has been the default firewall on Ubuntu servers since version 8.04 LTS.
It provides a user-friendly interface for managing open ports and is built on top of the classic [iptables](https://linux.die.net/man/8/iptables) tool.
I've read that some advanced users prefer to use iptables directly but I really like the balance that UFW offers.
It provides a simple syntax for restricting network access based protocol (TCP vs UDP), port, and source IP address.

Here is basic example that allows all incoming TCP connections on port 22 (SSH):

```
ufw allow 22/tcp
```

Here is an iteration on the previous example that restricts SSH connections to _only_ those coming from address `1.2.3.4`:

```
ufw allow from 1.2.3.4 to any port 22 proto tcp
```

In addition to `allow`ing connections, you can also `limit` them which applies rate limiting with a simple strategy: block clients (by IP address) after 6 connection attempts within 30 seconds.
All of my servers already use this for limiting SSH access (pub-key auth is enforced so the rate limiting is mainly there to protect the server's CPU from all of the noise).

One last thing about UFW: even though it ships with Ubuntu, it is disabled by default.
Enabling it is as easy as running `ufw enable` but be careful!
If you don't have any rules in place when the firewall comes up, you could lock yourself out of the system.
Unless you are using a hosting provider that supports native consoles, you'll be unable to connect.
I learned this lesson the hard way: be sure to allow SSH connections _before_ enabling UFW.
These days, I use automation ([Terraform](https://www.terraform.io/) and [Ansible](https://www.ansible.com/)) to ensure that these operations always occur in the correct order.

# Ansible

Ansible has a [builtin module](https://docs.ansible.com/ansible/latest/collections/community/general/ufw_module.html) for idempotently configuring UFW.
The examples in the documentation are very thorough and even describe setting up rules for multiple source IPs!
To my [Minecraft role](https://github.com/theandrew168/devops/tree/master/roles/minecraft) I added a variable named `minecraft_allowed_ips` which is a simple list of IP addresses (or ranges) that are allowed to access the server.
When this list is present and non-empty, a UFW rule should be created **per address** that allows it to connect to port 25565.
When the list is empty, all connections should be allowed regardless of where they are coming from.

With a small bit of Ansible logic we can split the task into these two cases: one for permitting _specific_ IPs and one for permitting _all_ IPs:

```yaml
- name: Limit login attempts (allowed IP addresses)
  ufw:
    rule: limit
    src: "{{ item }}"
    port: "25565"
    proto: tcp
  with_items: "{{ minecraft_allowed_ips }}"
  no_log: yes
  become: yes
  become_user: root

- name: Limit login attempts (all IP addresses)
  ufw:
    rule: limit
    port: "25565"
    proto: tcp
  when: not minecraft_allowed_ips
  become: yes
  become_user: root
```

The combination of `when` and `with_items` partitions the control flow such that only one of these tasks will run.
In tandem, these tasks cleanly handle configuring servers that are either "open to just me and my friends" or "open to anyone".

# Conclusion

Overall, these changes were [fairly straightforward to implement](https://github.com/theandrew168/devops/commit/9aa74693962d3e2b3a655ddde68eeb59bfcc4e12) and increased my confidence that my golden age Minecraft server won't get compromised.
In total, it now has three layers of security:

1. Only specific usernames can connect (via the server's config)
2. Only specific IP addresses can connect (via UFW)
3. Connection attempts are rate limited (via UFW)

Even still, though, this isn't a bulletproof solution.
A bad actor could gain access to the server if they: had a hacked client, knew a valid username, and were connecting to server from an allowed IP address / range (perhaps this could happen if both the attacker and an allowed player were behind the same [CGNAT](https://www.rapidseedbox.com/blog/cgnat)).
I don't think is very likely but it _could_ happen.

That being said, this is just a Minecraft server.
Between the low stakes of the system itself, periodic backups of the data volume, and ease of deployment, I'm not super worried.
The security layers in place might not be perfect but they do save me from losing sleep about it.
I can rest easy knowing that my precious blocks are (probably) safe and sound.
