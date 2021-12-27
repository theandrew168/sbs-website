.POSIX:
.SUFFIXES:

default: build

.PHONY: blog
blog:
	go run scripts/blogify.go

.PHONY: css
css:
	tailwindcss --minify -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: blog css
	go build -o sbs .

.PHONY: watch
watch:
	tailwindcss --watch -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: run
run: blog
	go run .

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
