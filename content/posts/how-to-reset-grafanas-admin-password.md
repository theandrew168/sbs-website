---
date: 2024-03-17
title: "How to Reset Grafana's Admin Password"
slug: "how-to-reset-grafanas-admin-password"
---

The other day, I went to log into my self-hosted [Grafana dashboard](https://grafana.com/oss/grafana/) (to confirm that my metrics server wasn't running out of storage) and realized that I didn't remember the password!
I usually save credentials to Firefox's builtin password manager but I must've skipped that step when I initially configured Grafana.

Either way, I needed to figure out if there was a way to reset the password without having any SMTP settings configured (meaning I couldn't use the regular "Forget your password?" email link).
I do, however, have access to the physical server that hosts Grafana.
Perhaps it ships with some sort of password reset command line utility?

## Saved by the CLI

After a bit of digging, I [discovered](https://community.grafana.com/t/admin-password-reset/19455) that Grafana _does_ include a [password reset utility](https://grafana.com/docs/grafana/latest/cli/#reset-admin-password).
On the server that hosts Grafana, simply run the following command (as the `root` or `grafana` user):

```
grafana-cli admin reset-admin-password admin
```

Once this executes successfully, the Grafana login page will revert to the "change the default admin password" workflow.
Be sure to do this quickly before someone else comes along and snags admin access to your dashboard!

After setting a new password, I made sure that it was safely stored within my password manager.
Hopefully, I won't need to do this type of reset again anytime soon.
But if I do, it's nice to know that a convenient utility already exists to save the day.
