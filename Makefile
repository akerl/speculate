.PHONY: default build clean lint fmt test deps

PACKAGE = $(shell basename $(shell pwd))
NAMESPACE = github.com/akerl
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null)
GOPATH = $(CURDIR)/.gopath
BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(NAMESPACE)/$(PACKAGE)

GO = go
GOFMT = gofmt
GOX = $(BIN)/gox
GOLINT = $(BIN)/golint

build: deps $(GOX) fmt lint test
	$(GOX) \
		-ldflags '-X $(NAMESPACE)/utils.Version=$(VERSION)' \
		-gocmd="$(GO)" \
		-output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" \
		-os="darwin linux" \
		-arch="amd64"

clean:
	rm -rf $(GOPATH) bin

lint: $(GOLINT)
	$(GOLINT) -set_exit_status ./...

fmt:
	$(GOFMT) -l -w $$(find . -type f -name '*.go' ! -path './.*')

test: deps
	$(GO) test ./...

deps: $(BASE)
	rsync -ax /Users/akerl/src/akerl/speculate/.gopath .gopath

$(BASE):
	mkdir -p $(dir $@)
	ln -vsf $(CURDIR) $@

$(GOLINT): $(BASE)
	$(GO) get github.com/golang/lint/golint

$(GOX): $(BASE)
	$(GO) get github.com/mitchellh/gox

