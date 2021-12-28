# sbs-website
Shallow Brook Software's main website

## Design
Credit for this site's overall look and feel goes to [John Komarnicki](https://github.com/johnkomarnicki).
Check him out!
He does some great design work.
Internet Illustrations come from [Storyset](https://storyset.com/internet).

## Setup
This project depends on the [Go programming language](https://golang.org/dl/) and the [TailwindCSS CLI](https://tailwindcss.com/blog/standalone-cli).

## Running
Build the blog, let Tailwind watch for CSS changes (in a background process), and run the web server:
```bash
# make run
go run scripts/blogify.go
tailwindcss --watch -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css &
go run .
```
