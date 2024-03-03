---
date: 2024-03-02
title: "Functional Abstractions in TypeScript"
slug: "functional-abstractions-in-typescript"
tags: ["typescript"]
draft: true
---

For my "Bloggulus in Svelte" rewrite, I wanted to create some single-function abstrations.
One example was a `fetchPage` function which fetches an HTML page and sanitizes the output / strips out all tags.
Originally, I wrote an abstraction in the form of function types:

```ts
type PageFetcher = () => Promise<string>;
```

However, I soon realized that this didn't really accomplish my goals: providing just a function didn't allow me to cleanly provide "behind the scenes" state (like what specific page data to return).
If I wanted to bake in some specialized behavior, I'd have to write a higher-order function and hide the info within a closure.
Something like:

```ts
// capture the mocked page in a closure
function makePageFetcher(page: string): PageFetcher {
  const pageFetcher: PageFetcher = async (url: string) => Promise<string> {
    return page;
  }
  return pageFetcher;
}
```

While this _does_ technically work, I think that there is a better way.
It made me realize that I should instead write these abstractions as types / interfaces with function members instead of standalone functions.
In essence:

```ts
type PageFetcher = {
  fetchPage: (url: string) => Promise<string | undefined>;
};
```

Now, this type can be satisfied by any object / class that has a matching method signature.
Instead of a closure, I can use a class:

```ts
class MockPageFetcher {
  private page: string;
  constructor(page: string) {
    this.page = page;
  }
  async fetchPage(url: string): Promise<string> {
    return this.page;
  }
}
```

Despite being more code than the closure version, I think most developers will find classes to be easier to grok.
Additionally, you can create interfaces of _multiple_ functions and even data.
You aren't limited to a single function.
One class / object can be written that satisfies a type / inteface of multiple functions.

You can see this in action when it comes to storage interfaces:

```ts
type TagStorage = {
  create: (params: CreateTagParams) => Promise<Tag>;
  list: () => Promise<Tag[]>;
  readById: (id: string) => Promise<Tag | undefined>;
  delete: (tag: Tag) => Promise<void>;
};
```

The rest of the application depends on this interface and not any specific implementation.
Any class / object that has the matching method signatures can be used interchangeably.
The "real" impl can be a class that hides the database connection and uses SQL to talk to a postgres database.
A "fake" impl can be one that hides a simple `Record<string, Tag>` and impls each method as simple JS operations.
