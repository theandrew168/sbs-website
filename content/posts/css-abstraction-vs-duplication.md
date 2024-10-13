---
date: 2024-10-13
title: "CSS Abstraction vs Duplication"
slug: "css-abstraction-vs-duplication"
---

CSS allows you get all the benefits of abstraction but without the commitment.
By applying the same style to multiple classes (via comma-separated [selector lists](https://developer.mozilla.org/en-US/docs/Web/CSS/Selector_list)), you can get the best of both worlds: the convenience of abstraction with the flexibility of duplication.
You don't actually have to repeat the CSS _and_ you don't have to commit to multiple designs always being identical.
If and when one design changes, you only need to split the selectors (a simple copy and paste) and update the modified one.

## Duplication

The main benefit of this approach is that HTML itself doesn't have to be refactored to "untangle" two designs that got prematurely abstracted into one.
For example, the "login" and "register" pages for [Bloggulus](https://bloggulus.com) currently look identical and share the same styles.
This is how I define their styles:

```css
.register,
.login {
  /* details omitted */
}

.register__form,
.login__form {
  /* details omitted */
}

.register__label,
.login__label {
  /* details omitted */
}
```

## Abstraction

What are our other options and why is this approach better?
Good question!
Someone feeling particularly [dry](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself) might be tempted to abstract these designs into one shared name (like a shared class named "auth" or something like that).
Here's what that'd look like:

```css
.auth {
  /* details omitted */
}

.auth__form {
  /* details omitted */
}

.auth__label {
  /* details omitted */
}
```

Would this be a good idea, though?
I think it boils down to intention.
Are these pages the same _by design_ or _by coincidence_?
In my case, it is definitely the latter: I'm sure these pages will eventually grow apart.
Plus, here's a spoiler alert: even pages that intentionally look the same now will likely require individual tweaks at some point in the future.

## Conclusion

Using this approach, I only had to write my styles once and then apply them to multiple classes.
If one of the designs diverges, it'll be a simple matter of splitting the selectors and modifying the one that changed.
I'm really glad that CSS supports this!
I still believe that a little duplication is better than [the wrong abstraction](https://sandimetz.com/blog/2016/1/20/the-wrong-abstraction).

Thanks for reading!
