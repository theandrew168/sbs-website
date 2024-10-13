---
date: 2024-10-12
title: "CSS Abstraction vs Duplication"
slug: "css-abstraction-vs-duplication"
draft: true
---

If multiple pages / components have the same style, CSS allows you easily apply the same styles to separate class names.
That way, instead of abstracting (giving the components a generic class name), you can just say “these two things happen to have the same style”.
In the future, if that changes, the fix is easy because the styles have only been duplicated, not abstracted.

CSS allows you get all the benefits of abstraction without the commitment.
If multiple pages have the same layout, just give them separate class names but use a comma-separated [selector list](https://developer.mozilla.org/en-US/docs/Web/CSS/Selector_list) to apply the same styles.
You didn't actually have to repeat the CSS, but you also didn't commit to the two designs always being identical.
If one of the pages changes, you can _then_ copy and paste the styles and change the second one.
The HTML itself doesn't have to be refactored to "untangle" two designs that got prematurely abstracted into one.

For example, [Bloggulus](https://bloggulus.com) has a couple groups of pages that share the same style:

1. Login and Register
2. Blogs list, Pages list (WIP) and Accounts list (admin only)

Someone feeling particularly [dry](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) might be tempted to abstract one or both of these designs into one shared name.
To be fair, the second example does brush up against the [rule of three](<https://en.wikipedia.org/wiki/Rule_of_three_(computer_programming)>).
For the first example, maybe I could create a new class called "auth" and for the second, a class named "list".
Would that be a good idea, though?

I think it boils down to intention.
Are these pages the same _by design_ or _by coincidence_?
In my case, it is definitely the latter.
And here's a spoiler alert: even pages that intentionally look the same will probably required tweaks and customizations eventually.
So, maybe advice is to save yourself the trouble and just apply the same style to multiple classes.
Let's look at some examples.

Here is how I style the login and register pages in CSS:

```css
.register,
.login {
  /* details omitted */
}

.register__form,
.login__form {
  /* details omitted */
}

.register__heading,
.login__heading {
  /* details omitted */
}

.register__label,
.login__label {
  /* details omitted */
}

.register__error,
.login__error {
  /* details omitted */
}
```

I only had to write the styles once and then apply them to multiple classes.
Pretty convenient, in my opinion.
I'm really glad CSS supports this.
I still believe that a little duplication is better than [the wrong abstraction](https://sandimetz.com/blog/2016/1/20/the-wrong-abstraction).

Thanks for reading!
