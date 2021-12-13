PROG:=day13

VERSION:=$(shell git describe --dirty=+WiP --always)

GOFLAGS:=-v -ldflags "-X app.version=$(VERSION) -v"

GOOS ?= $(shell uname -s | tr A-Z a-z)
GOARCH:=amd64

GO:=GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go

default: ${PROG}

${PROG}: build

build:
	$(GO) build $(GOFLAGS) -o ${PROG}

test:
	$(GO) test -v -cover

clean:
	@rm -f $(PROG)

install:
	install -b -c -s ${PROG} /usr/local/bin/

.PHONY: build clean

