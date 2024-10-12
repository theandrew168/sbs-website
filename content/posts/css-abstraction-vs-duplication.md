---
date: 2024-10-12
title: "CSS Abstraction vs Duplication"
slug: "css-abstraction-vs-duplication"
draft: true
---

If multiple pages / components have the same style, CSS allows you easily apply the same styles to separate class names.
That way, instead of abstracting (giving the components a generic class name), you can just say “these two things happen to have the same style”.
In the future, if that changes, the fix is easy because the styles have only been duplicated, not abstracted.
