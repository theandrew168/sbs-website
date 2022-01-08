.POSIX:
.SUFFIXES:

default: build

.PHONY: blog
blog:
	go run scripts/blogify.go

.PHONY: css
css: blog
	tailwindcss -m -i tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: blog css
	go build -o sbs .

.PHONY: run
run: blog
	go run . &
	tailwindcss --watch -m -i tailwind.input.css -o static/css/tailwind.min.css

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
