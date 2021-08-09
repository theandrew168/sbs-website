package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// References for custom HTML rendering:
// https://github.com/yuin/goldmark-highlighting
// https://github.com/yuin/goldmark-highlighting/blob/master/highlighting.go

// Config struct holds options for the extension.
type Config struct {
	html.Config
}

// NewConfig returns a new Config with defaults.
func NewConfig() Config {
	return Config{
		Config: html.NewConfig(),
	}
}

// Option interface is a functional option interface for the extension.
type Option struct {
	renderer.Option
}

// HTMLRenderer struct is a renderer.NodeRenderer implementation for the extension.
type HTMLRenderer struct {
	Config
}

// NewHTMLRenderer builds a new HTMLRenderer with given options and returns it.
func NewHTMLRenderer(opts ...Option) renderer.NodeRenderer {
	r := HTMLRenderer{
		Config: NewConfig(),
	}
	return &r
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *HTMLRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ast.Image)
	src := util.EscapeHTML(util.URLEscape(n.Destination, true))
	alt := util.EscapeHTML(n.Text(source))

	element := `<img class="todo-custom-class" src="%s" alt="%s" />`
	w.WriteString(fmt.Sprintf(element, src, alt))
	return ast.WalkContinue, nil
}

type sbs struct {
	options []Option
}

// SBS is a goldmark.Extender implementation.
var SBS = &sbs{
	options: []Option{},
}

// NewSBS returns a new extension with given options.
func NewSBS(opts ...Option) goldmark.Extender {
	e := sbs{
		options: opts,
	}
	return &e
}

// Extend implements goldmark.Extender.
func (e *sbs) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(e.options...), 200),
	))
}

type Post struct {
	Date    time.Time
	Title   string
	Slug    string
	Tags    []string
	Content string
}

func main() {
	srcdir := flag.String("srcdir", "blog/", "blog markdown input dir")
	destdir := flag.String("destdir", "posts/", "blog html output dir")

	// read in the post template
	postTmpl := filepath.Join(*srcdir, "templates", "post.html.tmpl")
	ts, err := template.ParseFiles(postTmpl)
	if err != nil {
		log.Fatal(err)
	}

	// setup markdown parser
	markdown := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			highlighting.Highlighting,
			meta.Meta,
			SBS,
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
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

		content := buf.String()

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
			Date:    date,
			Title:   title,
			Slug:    slug,
			Tags:    tags,
			Content: content,
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
	postsTmpl := filepath.Join(*srcdir, "templates", "posts.html.tmpl")
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