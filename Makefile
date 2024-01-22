.POSIX:
.SUFFIXES:

default: build

.PHONY: build
build:
	hugo -d docs/ --gc --minify

.PHONY: run
run:
	hugo server -D

.PHONY: clean
clean:
	rm -fr resources/
