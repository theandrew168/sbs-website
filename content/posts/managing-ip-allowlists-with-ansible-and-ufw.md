---
date: 2024-02-11
title: "Managing IP Allowlists With Ansible and UFW"
slug: "managing-ip-allowlists-with-ansible-and-ufw"
tags: ["minecraft", "ansible"]
draft: true
---

I recently stood up a "Golden Age" minecraft server for my friends and myself.
We've been having a great time so far: building houses, cave bases, and clifftop observation decks (I'm pretty proud of that one).
Since I wrote my last post, though, I knew that the server's security posture was a bit unstable.
As detailed in my [previous post](/posts/automating-a-golden-age-minecraft-server/), I was able to configure the server to only respect a specific set of usernames.

However, since the authetication servers for old versions of minecraft aren't online anymore, the server has to be run in "offline mode" which means it won't verify who is connecting: just their usernames.
This means that a malicious / hacked client could spoof the username of a known player and connect to the server.
I've been doing my best to keep the usernames secret but that is nothing more than a "security by obscurity" measure.
I outlines some options to increase the servers security in my last post and one stood out as a clear winner: only allowing known IP addresses to connect.

# UFW

[UFW](https://help.ubuntu.com/community/UFW) is the default firewall on Ubuntu servers since 8.04 LTS.
It provides a user-friendly experience above [iptables](https://linux.die.net/man/8/iptables).
I know that some people prefer to use iptables directly but I really like the balance that UFW offers.
It allows you to restrict access based on protocol, port, and source IP address.

Simple example of using UFW to restrict by proto, port, and IP:

```
ufw allow from 1.2.3.4 proto tcp to any port 22
```

In addition to "allowing" connection, you can also "limit" them which applies rate-limiting with a simple strategy: block clients (by IP address) after 6 connection attempts within 30 seconds.
All of my servers already use this for limiting SSH access (I already enforce pub-key only auth so the rate-limiting is mainly there to save the server's CPU usage).

# Ansible

Ansible has a [builtin module](https://docs.ansible.com/ansible/latest/collections/community/general/ufw_module.html) for idempotently configuring UFW.
The examples in the documentation are very thorough and even describe setting up rules for multiple source IPs.
My existing [Minecraft role](https://github.com/theandrew168/devops/tree/master/roles/minecraft) already uses UFW for allowing all incoming connections on port 25565.
This has been working for now and should still be the default if no specific IPs are listed in the playbook's vars.
For my Minecraft server, I'm using a variable called `minecraft_allowed_ips` which is a simple list of IP addresses (or ranges).

With a small bit of Ansible logic, we can split the task into two cases: one for permitting all IPs and one for permitting only specific IPs:

```yaml
- name: Limit login attempts (all IP addresses)
  ufw:
    rule: limit
    port: "25565"
    proto: tcp
  when: not minecraft_allowed_ips
  become: yes
  become_user: root

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
```

The combination of `when` and `with_items` effectively partitions the control flow: only one of these tasks will ever run during the playbook's execution.

# Conclusion

Overall, this change was [fairly easy to implement](https://github.com/theandrew168/devops/commit/9aa74693962d3e2b3a655ddde68eeb59bfcc4e12) and increased my confidence that our server won't get compromised.
It isn't a perfect fix, though, and the server still has _some_ risk with respects to unauthorized players connecting.
A bad actor could still compromise the server if they both: had a hacked client and knew a valid username AND was accessing the server from an allowed IP ([CGNAT](https://www.rapidseedbox.com/blog/cgnat) or something).

At this point, my server has three "layers" of security:

1. Only specific usernames can connect (via the server's config)
2. Only specific IP addresses can connect (via UFW)
3. Connection attempts are rate limited (via UFW)

That being said, this is just a Minecraft server.
I would be upset if anything happened to it but it isn't a critical system.
The security layers in place might not be perfect but they do allow me to sleep without worrying too much about it.
