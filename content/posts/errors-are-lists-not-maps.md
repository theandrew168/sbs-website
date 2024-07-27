---
date: 2024-07-27
title: "Errors Are Lists, Not Maps"
slug: "errors-are-lists-not-maps"
draft: true
---

All REST APIs must decide how to handle and represent errors.
There are many ways to represent errors and they all come with varying pros and cons.

# Errors as Maps

Historically, I've represented these errors as a map: the key is the problematic field and the value is the error message for that field.
This enables the frontend to display the errors right next to the invalid input which is great for user experience.

In TypeScript, this error response would look something like:

```js
// This could be extended to support multiple errors per field, if necessary.
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
Continuing with the example of a login page, the error for a failed login is often combined into an intentionally-ambiguous "invalid username or password".
Which field should this error be attached to? Both? Neither?
In my opinion, the answer is "neither": this is a general error.

How can we represent this general error if our format requires a key-value relationship?
Our best best is probably to add another field called "general" to the map and hope that the name never collides with an actual input value:

```js
const example: ErrorResponse = {
  // We _have_ to choose a field even for general errors. :(
  errors: {
    general: "invalid username or password",
  },
};
```

While this would probably work, it feels a bit like we are fighting against the design.
If "fields" are optional, how can we represent multiple errors without requiring specific keys?

# Errors as Lists

Let's rewind a bit and start with what we _know_ about errors:

1. A single request may yield **multiple** errors
2. They **always** have a message of some sort
3. They are **sometimes** tied to a specific fields (for input validation)

So, errors always have a "message", optionally have a "field", and can come in groups.
Thinking about these requirements in isolation point me toward a slightly different design.
Instead of representing errors as a map, let's represent them as a list:

```js
type Error = {
  message: string,
  field?: string,
};

type ErrorResponse = {
  errors: Error[],
};
```

Now, we much more easily represent basic, general errors:

```js
const example: ErrorResponse = {
  // No field? No problem!
  errors: [
    {
      message: "invalid username or password",
    },
  ],
};
```

And our login validation example from before still works despite looking a bit different:

```js
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

# Benefits

Now each category of error (general vs specific) can be easily represented without fighting against the design of our error response.
Furthermore, this format support multiple general errors AND multiple specific errors (per field) out of the box.
If your frontend is equipped to handle multiple errors per category, go for it!
Otherwise, you can always just find the first error per category in the list and call it a day.
I even wrote up a couple helpers for that:

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

# Conclusion

Errors are a big topic with countless approaches and opinions.
The internet is full of awesome discussions about [strategies](https://stackoverflow.com/questions/39759906/validation-responses-in-rest-api) and [examples](https://www.baeldung.com/rest-api-error-handling-best-practices) of how big companies do it.
Overall, both of the approaches outlined in this post are capable of getting the job done.
At the end of the day, it mostly comes down to personal preferences and the requirements of the project at hand.
For [Bloggulus](https://github.com/theandrew168/bloggulus), I'll be sticking with "errors as lists".

Thanks for reading!
