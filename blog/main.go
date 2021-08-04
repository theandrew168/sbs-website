package main

import (
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Post struct {
	Name string
}

func main() {
	srcdir := flag.String("srcdir", "blog/", "blog markdown input dir")
	destdir := flag.String("destdir", "posts/", "blog html output dir")

	// find source markdown files
	files, err := os.ReadDir(*srcdir)
	if err != nil {
		log.Fatal(err)
	}

	// read in the post template
	postTmpl := filepath.Join(*srcdir, "post.html.tmpl")
	ts, err := template.ParseFiles(postTmpl)
	if err != nil {
		log.Fatal(err)
	}

	var posts []Post

	// convert each markdown post to html
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext != ".md" {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ext)
		outdir := filepath.Join(*destdir, name)

		// create the output dir
		err := os.MkdirAll(outdir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// create the output html file
		outfile := filepath.Join(outdir, "index.html")
		f, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}

		// markdown content + html template = blog post!
		post := Post{
			Name: name,
		}
		err = ts.Execute(f, &post)
		if err != nil {
			log.Fatal(err)
		}

		// record the post for later
		posts = append(posts, post)
	}

	// read in the posts template
	postsTmpl := filepath.Join(*srcdir, "posts.html.tmpl")
	ts, err = template.ParseFiles(postsTmpl)
	if err != nil {
		log.Fatal(err)
	}

	// create the posts (ToC) html file
	outfile := filepath.Join(*destdir, "index.html")
	f, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}

	// read and apply the posts template
	err = ts.Execute(f, posts)
	if err != nil {
		log.Fatal(err)
	}
}
