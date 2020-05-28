include golang.mk

SHELL := /bin/bash
export PATH := $(PWD)/bin:$(PATH)
APP_NAME ?= ecs-task-metadata-exporter
EXECUTABLE = $(APP_NAME)
PKG = github.com/Clever/$(APP_NAME)
PKGS := $(shell go list ./... | grep -v /vendor)

$(eval $(call golang-version-check,1.13))

.PHONY: all test build run $(PKGS) generate install_deps

all: test build

test: $(PKGS)
$(PKGS): golang-test-all-strict-deps
	$(call golang-test-all-strict,$@)

build:
	$(call golang-build,$(PKG),$(EXECUTABLE))

run: build
	IS_LOCAL=true bin/$(EXECUTABLE)

install_deps:
	go mod vendor
