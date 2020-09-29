.POSIX:
.SUFFIXES:

default: dist

.PHONY: dist
dist:
	hugo -d docs/

.PHONY: server
server:
	hugo server -D
