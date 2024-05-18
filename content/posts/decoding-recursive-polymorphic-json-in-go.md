---
date: 2024-05-18
title: "Decoding Recursive Polymorphic JSON in Go"
slug: "decoding-recursive-polymorphic-json-in-go"
draft: true
---

References:

1. https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/
2. https://github.com/karaatanassov/go_polymorphic_json
3. https://alexkappa.medium.com/json-polymorphism-in-go-4cade1e58ed1
4. https://en.wikipedia.org/wiki/Polymorphism_(computer_science)
5. https://gist.github.com/theandrew168/75a5d97d794e7835670eaaa6c6d53b30

Recently, I was helping a friend design a system for matching text against a flexible system of rules.
A rule can either be a single regex pattern or a series of patterns combined with a logical operation (`AND`, `OR`, or `NOT`).
We'll call this first, pattern-based rule `Basic` and the second rule that joins other rules together `Composite`.
Since the `Composite` rule can contain both `Basic` and other `Composite` rules, we need a third type to represent "can be either basic or composite".
We'll call this type `Rule`.

In TypeScript, these types are very easy to represent:

```ts
// A Basic rule is a single regex pattern.
type Basic = {
  pattern: string;
};

// A Composite rule combines other rules with a logical operation.
type Operation = "and" | "or" | "not";
type Composite = {
  operation: Operation;
  rules: Rule[];
};

// A Rule can be either Basic or Composite.
type Rule = Basic | Composite;
```

In software design terminology, a `Rule` is a polymorphic type that can be represented by other more specific types.
However, Go doesn't (directly) support algebraic types (sum types, in this case).
Instead, we have to make clever use of Go's interface system to emulate sum types.

# Algebraic Data Types

Credit for this knowledge goes to Eli Bendersky's amazing blog post: [Go and Algebraic Data Types ](https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/).
Since we can't define a rule as the sum of two other types, we need to flip things around a bit.
Instead, we define an interface that represents a `Rule` and contains a single, private method named `isRule`:

```go
// A Rule can be either Basic or Composite.
type Rule interface {
	isRule()
}
```

Now, any specific, concrete subtypes of `Rule` can implement this private method and Go will understand that is "is a rule".
Let's fill in the two concrete rule subtypes: `Basic` and `Composite`:

```go
// A Basic rule is a single regex pattern.
type Basic struct {
	Pattern string `json:"pattern"`
}

func (Basic) isRule() {}

type Operation string

const (
	OperationAnd Operation = "and"
	OperationOr  Operation = "or"
	OperationNot Operation = "not"
)

// A Composite rule combines other rules with a logical operation.
type Composite struct {
	Operation Operation `json:"operation"`
	Rules     []Rule    `json:"rules"`
}

func (Composite) isRule() {}
```

Now that we have defined a `Basic` rule, a `Composite` rule, and a simple string type for `Operation`, we can represent arbitrary combinations of rules.
Here are a few examples:

```go
// This rule matches the text "foo" or "bar".
var fooOrBar Rule = Composite{
	Operation: OperationOr,
	Rules: []Rule{
		Basic{Pattern: "foo"},
		Basic{Pattern: "bar"},
	},
}

// This rule matches "foo" and not "bar".
var fooAndNotBar Rule = Composite{
	Operation: OperationAnd,
	Rules: []Rule{
		Basic{Pattern: "foo"},
		Composite{
			Operation: OperationNot,
			Rules: []Rule{
				Basic{Pattern: "bar"},
			},
		},
	},
}
```

I won't go into detail about how the evaluation of these rules is written, but feel free to check out the full code on [GitHub](https://gist.github.com/theandrew168/75a5d97d794e7835670eaaa6c6d53b30).
It is based on a simple and clean recursion!
This post isn't about that, though: this post is about converting these polymorphic rules to and from JSON.

# Simple JSON

# Polymorphic JSON

# All Together

```go
// Recursively parse a dynamic, polymorphic rule structure.
func ParseRule(data []byte) (Rule, error) {
	// base case

	// try to decode a basic rule
	var basic Basic
	err := json.Unmarshal(data, &basic)
	if err != nil {
		return nil, err
	}

	// if the rule has a pattern, it is basic so return it
	if basic.Pattern != "" {
		return basic, nil
	}

	// recursive case

	// decode a partial composite and its sub-rules
	var partial struct {
		Operation Operation         `json:"operation"`
		Rules     []json.RawMessage `json:"rules"`
	}
	err = json.Unmarshal(data, &partial)
	if err != nil {
		return nil, err
	}

	// decode each sub-rule recursively and collect into a slice
	var rules []Rule
	for _, ruleData := range partial.Rules {
		rule, err := ParseRule(ruleData)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	// construct and return the composite rule
	composite := Composite{
		Operation: partial.Operation,
		Rules:     rules,
	}
	return composite, nil
}
```

# Conclusion

This post explained and demonstrated how recursive, polymorphic data structures to converted to and from JSON.
While dealing with this sort of dynamic JSON isn't typically regarded as one of Go's strong suits, it is doable with a bit of recursion and partial decoding via `json.RawMessage`.
In some ways, I kind of like how specific and verbose this code needs to be in order to wrangle these complex structures.
In TypeScript, you could simply call `JSON.parse` and assert the resulting type.
While that is certainly easier, I don't necessarily think it is better.

With the approach shown by the `ParseRule` function, the incoming JSON is verified to be valid all the way down.
Furthermore, the data type we get back is the more structured `Rule` type, meaning that the rest of the code can have confidence that the rules being evaluated are well-structured.
This idea has been discussed before: see [Parse, don't validate](https://lexi-lambda.github.io/blog/2019/11/05/parse-don-t-validate/).
Thanks to that post, I had the foresight to write the function as `ParseRule` and not `ValidateRule`!

Thanks for reading!
