---
date: 2024-10-26
title: "OAuth Auth Only"
slug: "oauth-auth-only"
draft: true
---

I recently made a decision around how Bloggulus handles authentication.
In short, it now only supports signing in with OAuth (via GitHub or Google).
Users can no longer sign up with an arbitrary username and password.
Is this a terrible or big brain strategy?
I think the jury is still out but, personally, I'm feeling pretty good about it!
This approach (like everything in software development) comes with tradeoffs.

![Sign in options with OAuth services](/images/20241027/auth.webp)

## Tradeoffs

Pros:

- Users don't have to worry about another password.
- Getting started with the app is much simpler: just a couple clicks.
- Bloggulus is not responsible for storing any sensitive user data.
- Password resets are no longer an issue (don't need to send any reset emails).

Cons:

- Users without a GitHub or Google account can't use the app (Who would this affect? Could always support more services).

## Data Privacy

The only data I store is a message address code of the OAuth service’s unique ID for the authenticated user.
The process looks like this:

1. The OAuth services gives me a unique user identifier (which is public information)
2. Bloggulus appends the service name to the ID to make it unique across all OAuth services
3. The user ID is passed through an HMAC-SHA256 function alongside an application secret key
4. The resulting hash is stored in the database and can be used to deterministically identify return users

I really appreciate the info-sec posture that this approach enables.
What does the security posture look like?
I think this approach has great respect for user data privacy.
In the event of a database leak (due to a SQL injection or something like that), the attacker would not have enough information to brute force the hashes.
Even if they started from a list of OAuth service IDs, they wouldn't be able to arrive at the same hash without also having the application secret key.

Now, let's say the user _does_ get ahold of both the database AND the secret key.
Even in this absolute catastrophe of a scenario, the only information the attacker could obtain is the user's public service ID.
While this _does_ constitute personally-identifiable information (PII), it isn't sensitive and poses no risk to the user's upstream OAuth account.

## Local Development

Obviously, going OAuth only means that, by default, logging in locally requires settings up the config settings for at least one OAuth service.
This is kinda inconvenient, especially working on the auth system itself is pretty rare (compared to other things).
So how can we bridge this gap and avoid sacrificing developer experience?

I ended up implementing a second, local-only debug authentication method.
It only shows when the app is run with the `ENABLE_DEBUG_AUTH` environment variable set.
The way it works is quite simple: generate a random user ID, create an account, then create a session.
This means that each debug login will result in a new account being created (which could be a downside).
However, I find this to be a good balance!

![Sign in options with local debug auth](/images/20241027/debug.webp)

## Conlusion

Well reader, what do you think?
Is this an amazing idea that respects user privacy?
Or is this a selfish and inconvenient bar to set: requiring users to have an account on either GitHub or Google?
I'm still undecided (but I think it's a good idea overall).
At the very least, I learned something about how to integrate Go-based web apps with OAuth.
That alone is worth the trouble.
Investing in yourself is one of the safest investments you can make!
