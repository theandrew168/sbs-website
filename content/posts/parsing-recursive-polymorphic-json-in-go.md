---
date: 2024-05-19
title: "Parsing Recursive Polymorphic JSON in Go"
slug: "parsing-recursive-polymorphic-json-in-go"
draft: true
---

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

In software design terminology, a `Rule` is a [polymorphic type](<https://en.wikipedia.org/wiki/Polymorphism_(computer_science)>) that can be represented by other more specific types.
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

Usually, encoding (marshalling) and decoding (unmarshalling) JSON in Go is quite simple.
To encode a struct, just pass it to [json.Marshal](https://pkg.go.dev/encoding/json#Marshal).
To decode a struct, pass the raw JSON (as a `[]byte`) and a variable of the incoming type and pass them to [json.Unmarshal](https://pkg.go.dev/encoding/json#Unmarshal).
Here is what that looks like when dealing with a basic rule:

```go
func main() {
    var foo Basic = Basic{
        Pattern: "foo",
    }

    encoded, err := json.Marshal(foo)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("json: %s\n", encoded)

    var decoded Basic
    err = json.Unmarshal(encoded, &decoded)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("data: %+v\n", decoded)
}
```

When executed, this program prints:

```
json: {"pattern":"foo"}
data: {Pattern:foo}
```

This is exactly what we expected: the `Basic` struct can be easily transformed to and from JSON.
However, things a get a bit tricky once we introduce the polymorphic `Rule` type.
Let's try this experiment again with a `Composite` struct and see how the [encoding/json](https://pkg.go.dev/encoding/json) package reacts.

# Polymorphic JSON

This snippet is identical to the one above except that we trying to encode and decode a `Composite` rule.
We'll use the `fooOrBar` example from earlier.

```go
func main() {
    var fooOrBar Composite = Composite{
        Operation: OperationOr,
        Rules: []Rule{
            Basic{Pattern: "foo"},
            Basic{Pattern: "bar"},
        },
    }

    encoded, err := json.Marshal(fooOrBar)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("json: %s\n", encoded)

    var decoded Composite
    err = json.Unmarshal(encoded, &decoded)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("data: %+v\n", decoded)
}
```

Now let's see what happens:

```go
json: {"operation":"or","rules":[{"pattern":"foo"},{"pattern":"bar"}]}
json: cannot unmarshal object into Go struct field Composite.rules of type main.Rule
```

Well, that didn't work.
While encoding our polymorphic `Rule` interface works just fine, Go's JSON package doesn't know how to decode it.
This is because, when encoding, the specific type and structure of each individual rule is known.
When decoding, however, the `encoding/json` package doesn't have enough context and knowledge of the data structure to figure out a way to correctly parse it into concrete structs.
We are going to need to write some code to help fill in these gaps.

# Recursive Descent

Figuring out how to make this work took quite a bit of research.
I owe credit to Kiril Karaatanassov's [go_polymorphic_json](https://github.com/karaatanassov/go_polymorphic_json) project and Alex Kalyvitis's [JSON polymorphism in Go](https://alexkappa.medium.com/json-polymorphism-in-go-4cade1e58ed1) article for helping me to arrive at a clean, working solution.

At a high level, I knew that the parsing process would require a few specific steps:

1. Try to decode a basic rule
2. If that succeeds, return it
3. Otherwise, partially decode a composite rule
4. Recursively parse and collect each sub-rule
5. Construct the composite rule and return it

Step 3 seemed the most mysterious to me: how can you "partially decode" JSON in Go?
If I know that a rule is composite, it must have two keys: `operation` and `rules`.
The `operation` field must contain a string and `rules` should be an array of... other rules.
But since I don't know _what_ those rules are yet, I need to avoid decoding them directly and instead pass their raw JSON text back into our parsing function.
How can this be done?

Thankfully, the `encoding/json` is back to help us out again!
The critical piece of missing tech is [json.RawMessage](https://pkg.go.dev/encoding/json#RawMessage).
From the docs:

> RawMessage is a raw encoded JSON value. It implements Marshaler and Unmarshaler and can be used to delay JSON decoding or precompute a JSON encoding.

Sounds like exactly what we need.
Rather than decoding directly into a `Composite` struct, we'll instead define a custom type to represent a partial composite rule.
The two types look pretty similar:

```go
type Composite struct {
    Operation Operation `json:"operation"`
    Rules     []Rule    `json:"rules"`
}

type PartialComposite struct {
	Operation Operation         `json:"operation"`
	Rules     []json.RawMessage `json:"rules"`
}
```

Okay, at this point we have everything we need to parse our polymorphic JSON structure.
Let's get recursive!

```go
// Recursively parse a dynamic, polymorphic rule structure.
func ParseRule(data []byte) (Rule, error) {
	// 1. Try to decode a basic rule
	var basic Basic
	err := json.Unmarshal(data, &basic)
	if err != nil {
		return nil, err
	}

	// 2. If that succeeds, return it
	if basic.Pattern != "" {
		return basic, nil
	}

	// 3. Otherwise, partially decode a composite rule
	var partial struct {
		Operation Operation         `json:"operation"`
		Rules     []json.RawMessage `json:"rules"`
	}
	err = json.Unmarshal(data, &partial)
	if err != nil {
		return nil, err
	}

	// 4. Recursively parse and collect each sub-rule
	var rules []Rule
	for _, ruleData := range partial.Rules {
		rule, err := ParseRule(ruleData)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	// 5. Construct the composite rule and return it
	composite := Composite{
		Operation: partial.Operation,
		Rules:     rules,
	}
	return composite, nil
}
```

This function contains all of the logic necessary to parse all rules that the system can represent.
From simple regex patterns to complex combinations of rules, this single function can handle them all.
To wrap things up, let's use `ParseRule` to fix the broken code from earlier (this is the last code snippet, I swear):

```go
func main() {
    var fooOrBar Composite = Composite{
        Operation: OperationOr,
        Rules: []Rule{
            Basic{Pattern: "foo"},
            Basic{Pattern: "bar"},
        },
    }

    encoded, err := json.Marshal(fooOrBar)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("json: %s\n", encoded)

    decoded, err := ParseRule(encoded)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("data: %+v\n", decoded)
}
```

Instead of getting an error, we are now able to encode and decode the `fooOrBar` rule successfully:

```
json: {"operation":"or","rules":[{"pattern":"foo"},{"pattern":"bar"}]}
data: {Operation:or Rules:[{Pattern:foo} {Pattern:bar}]}
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
