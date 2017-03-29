include golang.mk
.DEFAULT_GOAL := test # override default goal set in library makefile

SHELL := /bin/bash
PKG := github.com/Clever/gearman-load-logger
PKGS := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := gearman-load-logger
.PHONY: all build test $(PKGS) vendor

$(eval $(call golang-version-check,1.7))

all: test build

run:
	go run main.go

build:
	go build -o bin/$(EXECUTABLE) $(PKG)

test: $(PKGS)

$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

vendor: golang-godep-vendor-deps
	$(call golang-godep-vendor,$(PKGS))
