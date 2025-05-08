---
date: 2024-03-03
title: "Behavioral Abstractions in TypeScript"
slug: "behavioral-abstractions-in-typescript"
tags: ["TypeScript"]
---

Almost every application needs to communicate with the outside world in one form or another.
It could be by scraping a web page, hitting a REST API, or simply talking to a database.
If an application depends on the specific details of any of these communications, however, then it becomes much more difficult to test.
Instead, a program's domain logic should depend on _abstract behaviors_ instead of concrete implementations.
I think that this is one of the most important facets of system design.

## Fetching Web Pages

I've recently been implementing my [Bloggulus](https://bloggulus.com/) web application in [Svelte + TypeScript](https://github.com/theandrew168/bloggulus-svelte/tree/main) for fun and experience.
One step of the syncing process involves manually retrieving a post's content if it isn't present in the RSS / Atom feed.
My initial approach to this was to simply use [fetch](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API/Using_Fetch) inline:

```ts
async function syncPost(post: Post): Promise<void> {
  let content = post.content;
  if (!content) {
    const resp = await fetch(post.url);
    content = await resp.text();
  }
  // insert new post, etc
}
```

This works, so what is the problem?
Well, the problem arises when it comes to testing this process.
If I want to verify that posts with missing content are handled correctly, I have to give the post a legitimate URL.
Otherwise, the code will attempt to fetch an invalid page and throw an error.
Sure, I _could_ test against a real web page (like [example.com](https://example.com/) or one of my own sites) but do I really want to write tests that depend on:

1. Having an active internet connection
2. The internet connection being stable and consistent
3. The permanent existence of a specific web page

That sounds like a pain.
Perhaps I could spin up a small web server as part of the testing to serve mock pages on a local port?
Again, I _could_ do that, but surely there must be a better way.

## Basic Behavior

What does my application actually depend on?
What does it need here?
It doesn't really need to know about the [Fetch API](https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API) at all.
In more abstract terms, the application just needs a bit of behavior that says: here is a URL, give me some content.
It needs a function that fetches pages:

```ts
type PageFetcher = {
  fetchPage: (url: string) => Promise<string>;
};
```

With this basic behavioral abstraction, we can make a small (but important) refactor to our code:

```ts
async function syncPost(post: Post, pageFetcher: PageFetcher): Promise<void> {
  let content = post.content;
  if (!content) {
    content = await pageFetcher.fetchPage(post.url);
  }
  // insert new posts, etc
}
```

With this small change, our code no longer depends on the implementation details of fetching pages: it only depends on the behavior in abstract.
The value added by this change becomes clear when we revisit the problem of writing tests for the sync process.
Instead of having to jump through some the painful aforementioned hoops, we can just write a mock implementation that simulates a legitimate page retrieval:

```ts
class MockPageFetcher {
  private page: string;
  constructor(page: string) {
    this.page = page;
  }
  // this method implements the PageFetcher interface
  async fetchPage(url: string): Promise<string> {
    return this.page;
  }
}

test("syncPost", async () => {
  // construct a PageFetcher that always returns "Hello, World!"
  const mockPageFetcher: PageFetcher = new MockPageFetcher("Hello, World!");

  // attempt to sync a post using our mock fetcher
  const post: Post = { url: "https://example.com/hello" };
  await syncPost(post, mockPageFetcher);

  // the sync process should have fetched our "Hello, World!" page
  expect(post.body).to.equal("Hello, World!");
});
```

I don't know about you, but I think that this is very powerful!
Because our domain logic doesn't care about _how_ we fetch pages, we can pass it a special, hand-crafted implementation of the interface and it won't know the difference.
It only knows that it needs to fetch a page and has been provided with the tools to make it happen.

## Complex Behaviors

This idea can be extended to encompass sets of behaviors, as well.
Consider how this can be applied to decoupling your application code from a database.
You code doesn't _need_ to know about database connections and SQL queries.
It only needs to know how to perform basic storage operations on domain types (create, read, update, and delete).

Bloggulus has the concept of a tag.
Tags are strings that represent topics found within a post's content.
From an application point of view, tags can be created, read (individually or in bulk), and deleted.
Since we know what the behaviors are, let's put them into an interface!

```ts
type TagStorage = {
  create: (params: CreateTagParams) => Promise<Tag>;
  list: () => Promise<Tag[]>;
  readById: (id: string) => Promise<Tag | undefined>;
  delete: (tag: Tag) => Promise<void>;
};
```

Just like before, the rest of the application depends only on this interface and not any specific implementation.
Any class or object that has the matching method signatures can be used in places where `TagStorage` is required.
The "real" implementation can be a class that hides the database connection and uses SQL to talk to a PostgreSQL database.
A "fake" implementation (for testing) can be one that stores tags in memory (likely in a `Record<string, Tag>`) and manipulates them with standard object operations.

## Conclusion

I think that it is always a good idea to separate behavior and implementation when it comes to how your app interacts with the world around it.
By representing behavior as an interface, you gain the freedom to implement the described functionality in different ways.
Despite only likely having one true implementation at runtime, the opportunities to provide finely-controlled mock implementations during testing adds a large amount of value.
It allows you to more easily test the most important (and often undertested) areas of your code: those that interact with external systems.

The concept of an interface is ubiquitous across software engineering and I use this pattern everywhere.
Even across different projects and programming languages, it is one of the handiest tools in my toolbelt.
It saves me from having to write nasty, tangled, and flaky integration tests.
Without behavioral abstractions, systems are doomed to be tightly coupled and loosely tested.
