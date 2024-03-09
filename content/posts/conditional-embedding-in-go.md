---
date: 2024-03-09
title: "Conditional Embedding in Go"
slug: "conditional-embedding-in-go"
tags: ["Go"]
draft: true
---

I was recently working on [a project](https://github.com/theandrew168/bloggulus) where the frontend is a [Svelte](https://svelte.dev/) SPA and the backend is a [Go](https://go.dev/) REST API.
Since Go is awesome, I've been using its [embed](https://pkg.go.dev/embed) feature to bake all of the pre-built frontend files into the single output binary.
This works fine when building for production: build the frontend, build the backend, done!
This snapshots the frontend and backend code into a single, static binary.

But for iterative development, I _don't_ want this behavior: I want the frontend changes to be reflected without needing to rebuild / restart the backend.
In short, I want to say:

```python
if building_for_production:
    # embed frontend output directory
else:
    # view frontend output directory
```

Since Go's embed feature is a compile-time feature, you can't really say things like "embed these files _if_ some condition is true".

# Research

So, I started doing some research.
Googling for "go conditional embedding" led me to this relevant [GitHub issue](https://github.com/golang/go/issues/44484) (a proposal for direct support for conditional embedding).
The proposal was rejected for a couple reasons:

1. Embedded files are decided at build time. To quote [Ian Lance Taylor](https://github.com/golang/go/issues/44484#issuecomment-948977876):

> When we build the binary, we have to decide whether to embed the file or not. We can't postpone that decision until the binary is run.

2. There is a workaround: [use build tags](https://github.com/golang/go/issues/44484#issuecomment-948137497)! With some clever file structure and a build tag, the question of "should these files be embedded" can effectively be pushed up into the build step.

# The Workaround

The workaround mentioned in the issue involves adding a few Go files to your frontend directory.
These files turn your frontend directory into a package that exposes two vars: one for the `embed.FS` itself and an `IsEmbedded` boolean indicating whether or not the data has been embedded.
This basic interface is implemented across three files:

1. `frontend/frontend.go` - Defines the exported vars
2. `frontend/frontend_embed.go` - Embeds the target directory and sets `IsEmbedded` to true
   1. Only included when the `embed` tag is provided at build time
   2. Build tags can be provided as command line flags: `go build -tags embed main.go`
3. `frontend/frontend_noembed.go` - Does NOT embed anything and sets `IsEmbedded` to false
   1. Only included when the `embed` tag is NOT provided at build time

It is then up to the importer of this frontend package to decide what to do with the exported information (some nuances omitted):

```go
package main

import (
    "io/fs"
    "os"

    "myproject/frontend"
    "myproject/backend/web"
)

func main() {
   // ...

    var frontendFS fs.FS
    if frontend.IsEmbedded {
        // used the embedded files if they are there
        frontendFS = frontend.Frontend
    } else {
        // otherwise just open the output dir directly
        frontendFS = os.DirFS("./frontend/output")
    }

    app := web.NewApplication(logger, store, frontendFS)
}
```

While this did work and solved my problem, I felt myself wanting something a bit cleaner.
What I really wanted was a package that made the decision of "embed or open" for itself and simply exposed a single `fs.FS` to the importer.

# Enhancements

I set about making these enhancements.
To be honest, it was quite easy to make the transition.
It ended up being both a nicer interface _and_ less code!
Let's dive in, starting with the one-var interface that this package exposes:

### frontend/frontend.go

```go
package frontend

import "io/fs"

var Frontend fs.FS
```

As simple as it gets!
Just export an `fs.FS` name `Frontend` that will be initialized when the package is imported.

### frontend/frontend_embed.go

```go
//go:build embed

package frontend

import (
    "embed"
    "io/fs"
)

//go:embed all:build
var frontend embed.FS

func init() {
    var err error
    Frontend, err = fs.Sub(frontend, "build")
    if err != nil {
        panic(err)
    }

    Frontend = os.DirFS("./frontend/build")
}
```

This file only gets imported when the `embed` build tag is present thanks to the `//go:build embed` comment at the top.
It embeds the `./frontend/build` directory which holds all of the frontend's static output files.
The extra [fs.Sub](https://pkg.go.dev/io/fs#Sub) is in here because `embed.FS` includes the embedded directory's top-level folder by default.
Calling `fs.Sub` lets us drill down one level so that the exported `fs.FS` includes just the internal files and not the top-level folder.

### frontend/frontend_noembed.go

```go
//go:build !embed

package frontend

import (
    "os"
)

func init() {
    Frontend = os.DirFS("./frontend/build")
}
```

Another easy one!
This file only gets imported when the `embed` build tag is _not_ present thanks to the `//go:build !embed` comment at the top.
Upon import, the `./frontend/build` directory will be opened and ready for reading (any changes to the frontend will be visible to readers of the exported `fs.FS`).
Note that we had to specify the full path to the `./frontend/build` directory because [os.DirFS](https://pkg.go.dev/os#DirFS) is rooted in the program's runtime directory (where the program was executed) versus the source file's directory (which is how the `embed` package works).

### main.go

```go
package main

import (
    "myproject/frontend"
    "myproject/backend/web"
)

func main() {
    // ...
    app := web.NewApplication(logger, store, frontend.Frontend)
}
```

As you can see, now the caller's code is much cleaner: it doesn't have to know anything about _how_ the `frontend.Frontend` FS is implemented.
It just imports the frontend package and passes the FS down to whatever part of the code plans to serve it.

# Conclusion

I went into this idea with some hesitance: I really didn't think I'd be able to find a clean solution to this problem.
Thankfully, I'm not the first person to consider "conditional embedding" so there was already some prior discussion and a workaround.
I took this workaround and iterated on it.
I transformed a "mehhh" interface for handling "which files do I need" into a clean, single-var import.

Will this solution tested and verified, I now get to have the best of both worlds: static binaries with embedded frontend files and a productive local development setup where frontend changes are immediately available.
I'm very happy with it!
