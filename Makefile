.POSIX:
.SUFFIXES:

default: build

.PHONY: build
build:
	hugo -d docs/

.PHONY: run
run:
	hugo server -D

.PHONY: clean
clean:
	rm -fr docs/ resources/
