.PHONY: default build clean lint fmt test deps

PACKAGE = speculate
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null)
GOPATH = $(CURDIR)/.gopath
BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(PACKAGE)

GO = go
GOFMT = gofmt
GOLINT = $(BIN)/golint
GOCOVMERGE = $(BIN)/gocovmerge
GOCOV = $(BIN)/gocov

build: deps fmt lint test
	$(GO) build \
		-ldflags '-X $(PACKAGE)/utils.Version=$(VERSION)' \
		-o bin/$(PACKAGE)

clean:
	rm -rf $(GOPATH) bin

lint: $(GOLINT)
	$(GOLINT) -set_exit_status ./...

fmt:
	$(GOFMT) -l -w $$(find . -type f -name '*.go' ! -path './.*')

test: deps
	$(GO) test ./...

deps:
	$(GO) get -d

$(BASE):
	mkdir -p $(dir $@)
	ln -sf $(CURDIR) $@

$(GOLINT): $(BASE)
	$(GO) get github.com/golang/lint/golint

