.PHONY: default build clean lint fmt test deps

PACKAGE = speculate
NAMESPACE = github.com/akerl
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2>/dev/null)
GOPATH = $(CURDIR)/.gopath
BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(NAMESPACE)/$(PACKAGE)
GOFILES = $(shell find . -type f -name '*.go' ! -path './.*' ! -path './vendor/*')
GOPACKAGES = $(shell echo $(GOFILES) | xargs dirname | sort | uniq)

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
	$(GOLINT) -set_exit_status $(GOPACKAGES)

fmt:
	@echo "Running gofmt on $(GOFILES)"
	@files=$$($(GOFMT) -l $(GOFILES)); if [ -n "$$files" ]; then \
		  echo "Error: '$(GOFMT)' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \
		  fi;

test: deps
	cd $(BASE) && $(GO) test $(GOPACKAGES)

deps: $(BASE) $(GOVEND)
	cd $(BASE) && $(GOVEND)

$(BASE):
	mkdir -p $(dir $@)
	rsync -ax --exclude '.gopath' --exclude '.git' $(CURDIR)/ $@

$(GOLINT): $(BASE)
	$(GO) get github.com/golang/lint/golint

$(GOX): $(BASE)
	$(GO) get github.com/mitchellh/gox

$(GOVEND): $(BASE)
	$(GO) get github.com/govend/govend
