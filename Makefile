.POSIX:
.SUFFIXES:

default: dist

.PHONY: build
build:
	go build -o sbs .

.PHONY: dist
dist: build
	rm -fr dist/
	mkdir dist/
	cp sbs dist/
	hugo -d dist/

.PHONY: server
server:
	hugo server -D

.PHONY: clean
clean:
	rm -fr dist/ public/ resources/
