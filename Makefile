.POSIX:
.SUFFIXES:

default: build

.PHONY: blog
blog:
	go run scripts/blogify.go

.PHONY: css
css:
	tailwindcss -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css

.PHONY: build
build: blog css
	go build -o sbs .

.PHONY: run
run: blog
	tailwindcss --watch -m -i static/css/tailwind.input.css -o static/css/tailwind.min.css &
	go run .

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
