.POSIX:
.SUFFIXES:

default: build

.PHONY: blog
blog:
	go run scripts/blogify.go

.PHONY: build
build: blog
	go build -o sbs .

.PHONY: run
run: blog
	go run . -conf sbs.conf

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
