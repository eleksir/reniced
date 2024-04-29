#!/usr/bin/env gmake -f

BUILDOPTS=-ldflags="-s -w" -a -gcflags=all=-l -trimpath

BINARY=reniced

all: clean build

build:
	go build ${BUILDOPTS} -o ${BINARY}
clean:
	go clean

upgrade:
	$(RM) -r vendor
	go get -d -u -t ./...
	go mod tidy
	go mod vendor

# vim: set ft=make noet ai ts=4 sw=4 sts=4:
