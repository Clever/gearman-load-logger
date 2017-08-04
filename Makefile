include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
PKG := github.com/Clever/gearman-load-logger
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := gearman-load-logger
.PHONY: all build test $(PKGS) vendor

$(eval $(call golang-version-check,1.8))

all: test build

run: build
	bin/$(EXECUTABLE)

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

$(GOPATH)/bin/glide:
	@go get github.com/Masterminds/glide

install_deps: $(GOPATH)/bin/glide
	@$(GOPATH)/bin/glide install
