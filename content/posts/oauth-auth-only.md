---
date: 2024-10-27
title: "OAuth Auth Only"
slug: "oauth-auth-only"
---

I recently made a decision around how Bloggulus handles authentication.
In short, it now only supports signing in with OAuth (via GitHub or Google).
Users can no longer sign up with an arbitrary username and password.
Is this a terrible idea or big brain play?
I think the jury is still out but, personally, I'm feeling pretty good about it!
This approach, like everything in software development, comes with tradeoffs.

![Sign in options with OAuth services](/images/20241027/auth.webp)

## Tradeoffs

Let's start with some of the benefits:

- Users don't have to worry about another password.
- Getting started with the app is much simpler: just a couple clicks.
- Bloggulus is not responsible for storing any sensitive user data.
- Password resets are no longer an issue (don't need to send any reset emails).

How about downsides?
I can't think of many but there is at least one:

- Users without a GitHub or Google account can't use the app

I'm unsure how big of a downside this is, though.
I don't personally know anyone who doesn't have either a GitHub (friends) or Google (family) account.
Plus, this approach is somewhat future proof because support for other OAuth services can always be added down the line.

## Data Privacy

The only data I store is a [message address code](https://en.wikipedia.org/wiki/Message_authentication_code) of the OAuth service's unique ID for the authenticated user.
The process looks like this:

1. The OAuth services gives me a unique user identifier (which is public information)
2. Bloggulus appends the service name to the ID to make it unique across all OAuth services
3. The user ID is passed through an HMAC-SHA256 function alongside an application secret key
4. The resulting hash is stored in the database and can be used to deterministically identify users

I really appreciate the "information security" posture of this approach.
I think it offers great respect for users and protects their data.
In the event of a database leak (due to SQL injection or something like that), the attacker would not have enough information to brute force the hashes.
Even if they started from a list of known, public OAuth service IDs, they wouldn't be able to arrive at the same hash without also having the application secret key.

Now, let's say the attacker gets ahold of both the database _and_ the secret key.
Even in this absolute catastrophe of a scenario, the only information the attacker could obtain are the users' service IDs.
While this _does_ constitute personally-identifiable information (PII), it isn't sensitive and poses no risk to each user's upstream OAuth account.

## Local Development

Obviously, going OAuth only means that, by default, logging in locally requires setting up the config for at least one OAuth service.
This is somewhat inconvenient, especially because working on the auth system itself is pretty rare (compared to other things).
So how can we bridge this gap and avoid sacrificing developer experience?

I ended up implementing another, local-only debug authentication method.
It only shows when the app is run with the `ENABLE_DEBUG_AUTH` environment variable set.
The way it works is quite simple: generate a random user ID, create an account, then create a session.
This means that each debug login will result in a new account being created.
While this could be a downside, I find it to be a good balance!

![Sign in options with local debug auth](/images/20241027/debug.webp)

## Conlusion

Well reader, what do you think?
Is this an amazing idea that respects user privacy?
Or is this a selfish and inconvenient precedent to set: requiring users to have an account on either GitHub or Google?
I'm still undecided but I ultimately think it's a good idea overall.
At the very least, I learned something about how to integrate Go-based web apps with OAuth services.
That alone is worth the trouble.

You know what they say: investing in yourself is one of the safest investments you can make!
