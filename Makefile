GO    := GO15VENDOREXPERIMENT=1 go
PROMU := $(GOPATH)/bin/promu
GOLINT := $(GOPATH)/bin/golint
GODEP_BIN := $(GOPATH)/bin/dep
pkgs   = $(shell $(GO) list ./... | grep -v /vendor/)

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)

all: vendor format build test

$(GODEP):
	$(GO) get -u github.com/golang/dep/cmd/dep

Gopkg.toml: $(GODEP)
	$(GODEP_BIN) init

vendor: $(GODEP) Gopkg.toml Gopkg.lock
	@echo ">> No vendor dir found. Fetching dependencies now..."
	GOPATH=$(GOPATH):. $(GODEP_BIN) ensure

test: lint vet
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

lint:
	@echo ">> linting code"
	@$(GOLINT) $(pkgs)

build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

#tarball: promu
#	@echo ">> building release tarball"
#	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
	        GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
	        $(GO) get -u github.com/prometheus/promu

clean:
	@rm -f sheepdog_exporter

.PHONY: all style format build test vet promu clean
