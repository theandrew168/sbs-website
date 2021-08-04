.POSIX:
.SUFFIXES:

default: build

.PHONY: build
build:
	go run blog/main.go
	go build -o sbs .

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp sbs dist/
	cp -r posts dist/
	cp -r static dist/
	cp -r templates dist/

.PHONY: clean
clean:
	rm -fr sbs dist/ posts/
