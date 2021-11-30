.POSIX:
.SUFFIXES:

default: build

.PHONY: blog
blog:
	go run scripts/blogify.go

.PHONY: 
build: blog
	go build -o sbs main.go

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp sbs dist/
	cp -r posts dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: test
test:
	go test -count=1 -v ./...

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
