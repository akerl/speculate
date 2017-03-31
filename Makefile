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
GOVEND = $(BIN)/govend

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
	files=$$($(GOFMT) -l $$(find . -type f -name '*.go' ! -path './.*')); if [ -n "$$files" ]; then \
		  echo "Error: '$(GOFMT)' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \
		  fi;

test: deps
	cd $(BASE) && $(GO) test ./...

deps: $(BASE) $(GOVEND)
	cd $(BASE) && $(GOVEND) -v

$(BASE):
	mkdir -p $(dir $@)
	rsync -ax --exclude '.gopath' --exclude '.git' $(CURDIR)/ $@

$(GOLINT): $(BASE)
	$(GO) get github.com/golang/lint/golint

$(GOX): $(BASE)
	$(GO) get github.com/mitchellh/gox

$(GOVEND): $(BASE)
	$(GO) get github.com/govend/govend
