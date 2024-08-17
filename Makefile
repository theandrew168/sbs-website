.POSIX:
.SUFFIXES:

default: build

.PHONY: build
build:
	hugo -d docs/ --gc --minify

.PHONY: run
run:
	hugo server -D

.PHONY: update
update:
	git submodule foreach git pull origin main

.PHONY: clean
clean:
	rm -fr resources/
