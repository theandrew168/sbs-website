---
date: 2024-03-10
title: "Conditional Embedding in Go"
slug: "conditional-embedding-in-go"
tags: ["Go"]
---

I was recently working on [a project](https://github.com/theandrew168/bloggulus) where the frontend is a [Svelte](https://kit.svelte.dev/) SPA and the backend is a [Go](https://go.dev/) REST API.
At a high level, this is how my project is structured:

```python
main.go
backend/
  # source files for Go REST API
frontend/
  package.json
  src/
    # source files for Svelte SPA
```

Since Go is awesome, I've been using its [embed](https://pkg.go.dev/embed) feature to bake all of the compiled frontend files into the single output binary.
This works great when building for production: compile the frontend, build the backend, done!
This snapshots the frontend and backend code into a single, static binary.
For iterative local development, however, I _don't_ want this behavior.
I instead want the frontend changes to be reflected without needing to rebuild and restart the backend.
In short, I want to say:

```python
if building_for_production:
    # embed frontend output directory (lock it in)
else:
    # view frontend output directory (keep it open)
```

However, since embedding files into Go binary is a compile-time feature, you can't really say things like "embed these files _if_ some condition is true".

# Research

So, I started doing some research.
Googling for "go conditional embedding" led me to this relevant [GitHub issue](https://github.com/golang/go/issues/44484) (a proposal for direct support for conditional embedding).
The proposal was rejected for a couple reasons:

1. Embedded files are decided at build time. To quote [Ian Lance Taylor](https://github.com/golang/go/issues/44484#issuecomment-948977876):

> When we build the binary, we have to decide whether to embed the file or not. We can't postpone that decision until the binary is run.

2. There exists a workaround: [use build tags](https://github.com/golang/go/issues/44484#issuecomment-948137497)! With some clever file structure and a build tag, the answer to "should these files be embedded" can essentially be pulled up into the build process.

# The Workaround

The workaround described in the proposal involves adding a few Go files to your frontend directory.
These files turn your frontend directory into a Go package that exposes two vars: one for the `embed.FS` itself and an `IsEmbedded` boolean indicating whether or not the data has been embedded.
This basic interface is implemented across three files:

1. `frontend/frontend.go` - Defines the exported vars
2. `frontend/frontend_embed.go` - Embeds the target directory and sets `IsEmbedded` to true
   1. Starts with the comment: `//go:build embed`
   2. Only included when the `embed` tag is provided at build time
   3. Build tags can be provided as command line flags: `go build -tags embed main.go`
3. `frontend/frontend_noembed.go` - Does NOT embed anything and sets `IsEmbedded` to false
   1. Starts with the comment: `//go:build !embed`
   2. Only included when the `embed` tag is NOT provided at build time

It is then up to the importer of this frontend package to decide what to do with the exported information (some nuances are omitted from this example):

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
        // use the embedded frontend output files if they are there
        frontendFS = frontend.Frontend
    } else {
        // otherwise just open the frontend output dir directly
        frontendFS = os.DirFS("./frontend/output")
    }

    app := web.NewApplication(logger, store, frontendFS)
}
```

Despite working as described and solving my problem, I felt myself wanting something a bit cleaner.
What I really wanted was a package that encapsulated the decision of "should this directory be embedded or opened" and exposed only a single [fs.FS](https://pkg.go.dev/io/fs#FS) to the importer.

# Enhancements

So, I set about making these enhancements.
To be honest, it was easier than I expected.
It ended up being both a nicer interface _and_ less code!
Let's dive in, starting with the one-var interface that this package exposes:

### frontend/frontend.go

```go
package frontend

import "io/fs"

var Frontend fs.FS
```

As simple as it gets!
Just export an `fs.FS` that will be initialized when the package is imported.

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
}
```

This file only gets imported when the `embed` build tag is present thanks to the `//go:build embed` comment at the top.
It embeds the `./frontend/build` directory which holds all of the frontend's compiled output files.

The call to [fs.Sub](https://pkg.go.dev/io/fs#Sub) is in here because `embed.FS` includes the embedded directory's top-level folder by default.
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
Upon import, the `./frontend/build` directory will be opened and ready for reading.
Any changes to the frontend will be immediately visible to readers of the exported `fs.FS` (assuming the frontend is running in some sort of "watch" mode that recompiles the code if anything changes).

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

As you can see, the caller's code is now much cleaner: it doesn't have to know anything about _how_ the `frontend.Frontend` FS is populated.
It just imports the frontend package and passes the exposed `fs.FS` down to whatever part of the code plans to serve it.

# Conclusion

I went into this problem with some pessimism: I really didn't think I'd be able to find a clean solution to the problem.
Thankfully, I'm not the first person to consider "conditional embedding" so there was already some prior discussion and a workaround to build upon.
I started with this workaround and iterated on idea.
I transformed a "good but not great" interface for handling conditional embedding into a clean, single-variable export.

Will this solution tested and verified, I now get to enjoy the best of both worlds: static binaries with embedded frontend files and a productive local development setup where frontend changes are reflected immediately.
I'm very happy with it!
