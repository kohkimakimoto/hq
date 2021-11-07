.DEFAULT_GOAL := help

# Load .env. see https://lithic.tech/blog/2020-05/makefile-dot-env
ifneq (,$(wildcard ./.env))
  include .envgit rev-parse --short HEAD
  export
endif

# Environment
GO111MODULE := on
PATH := $(CURDIR)/build/scripts:$(CURDIR)/.go-tools/bin:$(PATH)
SHELL := bash

VERSION := 2.0.0
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_LDFLAGS = -s -w \
                -X github.com/kohkimakimoto/hq/internal/version.CommitHash=$(COMMIT_HASH) \
                -X github.com/kohkimakimoto/hq/internal/version.Version=$(VERSION)


# Output help message
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: format
format: ## Format go code
	@go fmt ./...

dev: ## build dev binary
	@go build -o="build/outputs/dev/hq" ./cmd/hq

.PHONY: deps
deps: ## Install go modules
	@go mod tidy

.PHONY: tools/install
tools/install: ## Install dev tools
	@export GOBIN=$(CURDIR)/.go-tools/bin && \
		go install github.com/mitchellh/gox@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html@latest

.PHONY: tools/clean
tools/clean: ## Clean installed tools
	@rm -rf $(CURDIR)/.go-tools

.PHONY: test
test: ## Test go code
	@go test -race -timeout 30m -cover ./...

.PHONY: test/verbose
test/verbose: ## Run all tests with verbose outputting.
	@go test -race -timeout 30m -v -cover ./...

.PHONY: test/coverage
test/coverage: ## Run all tests with coverage report outputting.
	@gocov test ./... | gocov-html > coverage-report.html
