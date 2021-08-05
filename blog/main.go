package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	Date  time.Time
	Title string
	Slug  string
	Tags  []string
}

func main() {
	srcdir := flag.String("srcdir", "blog/", "blog markdown input dir")
	destdir := flag.String("destdir", "posts/", "blog html output dir")

	// read in the post template
	postTmpl := filepath.Join(*srcdir, "post.html.tmpl")
	ts, err := template.ParseFiles(postTmpl)
	if err != nil {
		log.Fatal(err)
	}

	// setup markdown parser
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	// keep track of post metadata (for /posts ToC)
	var posts []Post

	// find source markdown files
	files, err := os.ReadDir(*srcdir)
	if err != nil {
		log.Fatal(err)
	}

	// convert each markdown post to html
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext != ".md" {
			continue
		}

		// read in the markdown file
		infile := filepath.Join(*srcdir, file.Name())
		source, err := os.ReadFile(infile)
		if err != nil {
			log.Fatal(err)
		}

		// parse the markdown source
		var buf bytes.Buffer
		context := parser.NewContext()
		err = markdown.Convert([]byte(source), &buf, parser.WithContext(context))
		if err != nil {
			log.Fatal(err)
		}

		// grab metadata
		metaData := meta.Get(context)

		// resolve metadata types
		date, err := time.Parse("2006-01-02", metaData["date"].(string))
		if err != nil {
			log.Fatal(err)
		}
		title := metaData["title"].(string)
		slug := metaData["slug"].(string)

		var tags []string
		for _, tag := range metaData["tags"].([]interface{}) {
			tags = append(tags, strings.ToLower(tag.(string)))
		}

		// save the post metadata for rendering and the index
		post := Post{
			Date:  date,
			Title: title,
			Slug:  slug,
			Tags:  tags,
		}

		name := strings.TrimSuffix(file.Name(), ext)
		outdir := filepath.Join(*destdir, name)

		// create the output dir
		err = os.MkdirAll(outdir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// create the output html file
		outfile := filepath.Join(outdir, "index.html")
		f, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// markdown content + html template = blog post!
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
