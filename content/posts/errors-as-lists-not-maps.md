---
date: 2024-07-28
title: "Errors as Lists, Not Maps"
slug: "errors-as-lists-not-maps"
---

All REST APIs must decide how to handle and represent errors.
There are many ways to accomplish this task and they all come with varying pros and cons.
This post starts by describing the strategy I've historically used when dealing with errors.
After examining some limitations with that pattern, I present an alternative.

## Errors as Maps

Errors come in many shapes and sizes.
Since errors often arise in response to input validation, it can be useful to include what specific "field" caused the problem.
Until recently, I've represented these errors as a map: the key is the problematic field and the value is the error message for that field.
This enables the frontend to display the errors right next to the invalid input which is great for user experience.

In TypeScript, this error response would look something like:

```js
type ErrorResponse = {
  errors: Record<string, string>,
};

// Example error response from a login page.
const example: ErrorResponse = {
  errors: {
    username: "must not be empty",
    password: "must be at least 8 characters",
  },
};
```

However, I've often felt that this approach was lacking in one specific area: general errors.
By "general error", I'm referring to those that aren't tied to a specific field.
Continuing with the example of a login page, the error for a failed login is often combined into an intentionally-ambiguous "invalid username or password" message.
Which field should this error be attached to? Both? Neither?
In my opinion, the answer is "neither" because this is a general error.

How can we represent this field-less error if our format enforces a key-value relationship?
Our best bet is probably to add another field called "general" (or something similar) to the map and hope that the name never collides with an actual input value:

```ts
const example: ErrorResponse = {
  // We _have_ to choose a field even for general errors. :(
  errors: {
    general: "invalid username or password",
  },
};
```

While this would probably work, it feels a bit like we are fighting against the design.
If "fields" are an optional facet of our errors, how else can we represent them?

## Errors as Lists

Let's rewind a bit and start with what we _know_ about errors:

1. A single request may yield **multiple** errors
2. They **always** have a message of some sort
3. They are **sometimes** tied to specific fields

So, errors always have a "message", optionally have a "field", and can come in multiples.
Thinking about these requirements in isolation points me toward a slightly different design.
Instead of representing errors as a map, let's represent them as a list:

```ts
// An individual error has a message and an optional field.
type Error = {
  message: string;
  field?: string;
};

// An error response contains a list of individual errors.
type ErrorResponse = {
  errors: Error[];
};
```

With this structure, we can much more easily represent general errors:

```ts
const example: ErrorResponse = {
  // No field? No problem!
  errors: [
    {
      message: "invalid username or password",
    },
  ],
};
```

And our login input validation example from before still works despite looking a bit different:

```ts
const example: ErrorResponse = {
  // If necessary, fields can be included.
  errors: [
    {
      message: "must not be empty",
      field: "username",
    },
    {
      message: "must be at least 8 characters",
      field: "password",
    },
  ],
};
```

## Benefits

Now each category of error ("general" and "specific") can be easily represented without fighting against the design of our error response.
Furthermore, this format support multiple general errors AND multiple specific errors (per field) out of the box.
If your frontend is equipped to handle multiple errors per category, go for it!
Otherwise, you can always just find the first error per category and call it a day.
I even wrote up a couple helpers to make this "pick the first out of multiple errors" logic reusable:

```js
// Find the first general error.
export function findGeneralError(errors: Error[]): string | undefined {
  return errors.find((e) => !e.field)?.message;
}

// Find the first specific error for each field.
export function findSpecificErrors(errors: Error[]): Record<string, string> {
  const errorsByField = errors.reduce(
    (acc, err) => {
      if (err.field && !acc[err.field]) {
        acc[err.field] = err.message;
      }
      return acc;
    },
    {} as Record<string, string>,
  );
  return errorsByField;
}
```

## Conclusion

Errors are a big topic with countless approaches and opinions.
The internet is full of awesome discussions about [strategies](https://stackoverflow.com/questions/39759906/validation-responses-in-rest-api) and [examples](https://www.baeldung.com/rest-api-error-handling-best-practices) of how big companies do it.
Overall, both of the approaches outlined in this post are capable of getting the job done.
At the end of the day, it mostly comes down to personal preference and the requirements of the project at hand.

Thanks for reading!
