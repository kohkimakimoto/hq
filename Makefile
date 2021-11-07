.DEFAULT_GOAL := help

# Load .env file. see https://lithic.tech/blog/2020-05/makefile-dot-env
ifneq (,$(wildcard ./.env))
  include .env
  export
endif

# Environment
GO111MODULE := on
PATH := $(CURDIR)/dev/scripts:$(CURDIR)/dev/.external-tools/bin:$(PATH)
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
	@grep -E '^[/0-9a-zA-Z_-]+:.*?## .*$$' Makefile | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: format
format: format/go format/ui ## Format code

.PHONY: format/go
format/go: ## Format go code
	@find . -print | grep --regex '.*\.go' | xargs goimports -w -l -local "github.com/kohkimakimoto/hq"

.PHONY: format/ui
format/ui: ## Format ui js code
	@cd ui && yarn format

.PHONY: deps
deps: ## Install go modules
	@go mod tidy

.PHONY: build/dev
build/dev: ## build dev binary
	@cd ui && if [[ ! -e node_modules ]]; then yarn install; fi
	@cd ui && yarn dev
	@go build -o="dev/build/outputs/dev/hq" ./cmd/hq

.PHONY: build/release
build/release: ## build release binaries
	@cd ui && if [[ ! -e node_modules ]]; then yarn install; fi
	@cd ui && yarn clean && yarn build
	@rm -rf dev/build/outputs/release && rm -rf dev/build/outputs/archives
	@gox -os="linux darwin" -arch="amd64 arm64" -ldflags="${BUILD_LDFLAGS}" -output "dev/build/outputs/release/hq_{{.OS}}_{{.Arch}}" ./cmd/hq
	@mkdir -p dev/build/outputs/archives
	@cd dev/build/outputs && cp -f release/hq_darwin_amd64 archives/hq && cd archives && zip hq_darwin_amd64.zip hq && rm hq
	@cd dev/build/outputs && cp -f release/hq_darwin_arm64 archives/hq && cd archives && zip hq_darwin_arm64.zip hq && rm hq
	@cd dev/build/outputs && cp -f release/hq_linux_amd64 archives/hq && cd archives && zip hq_linux_amd64.zip hq && rm hq
	@cd dev/build/outputs && cp -f release/hq_linux_arm64 archives/hq && cd archives && zip hq_linux_arm64.zip hq && rm hq

.PHONY: build/rpm
build/rpm: build/release ## build RPM packages
	@rm -rf dev/build/outputs/rpm && mkdir -p dev/build/outputs/rpm
	@rpmtool.py build \
		--pre="cp -pr dev/build/outputs/archives/hq_linux_amd64.zip dev/build/rpmbuild/SOURCES/" \
		--post="rm -rf dev/build/rpmbuild/SOURCES/hq_linux_amd64.zip" \
		--out="$(CURDIR)/dev/build/outputs/rpm" \
		dev/build/rpmbuild/SPEC/hq.spec

.PHONY: build/yum
build/yum: build/rpm ## build yum repository (Experimental)
	@rm -rf docs/rhel/7/x86_64 && mkdir -p docs/rhel/7/x86_64
	@cp -f dev/build/outputs/rpm/x86_64/*.el7.x86_64.rpm docs/rhel/7/x86_64
	@rpmtool.py createrepo -- -v docs/rhel/7/x86_64

.PHONY: build/clean
build/clean: ## clean build outputs
	@cd dev/build && rm -rf outputs

.PHONY: dev/tools/install
dev/tools/install: ## Install dev tools
	@export GOBIN=$(CURDIR)/dev/.external-tools/bin && \
		go install golang.org/x/tools/cmd/goimports@latest && \
		go install github.com/mitchellh/gox@latest && \
		go install github.com/axw/gocov/gocov@latest && \
		go install github.com/matm/gocov-html@latest && \
		go install github.com/cosmtrek/air@latest

.PHONY: dev/tools/clean
dev/tools/clean: ## Clean installed tools
	@rm -rf $(CURDIR)/dev/.external-tools

.PHONY: dev/start
dev/start: prepare-tmp-dir ## start dev server
	@if [[ ! -e dev/.tmp/hq.toml ]]; then cp -f dev/hq.example.toml dev/.tmp/hq.toml; fi
	@process-starter.py --pre "cd ui && yarn dev" --run "cd ui && yarn watch" "air"

.PHONY: test
test: test/go test/ui ## Test all code

.PHONY: test/go
test/go: ## Test go code
	@go test -race -timeout 30m ./...

.PHONY: test/go/verbose
test/go/verbose: ## Run all tests with verbose outputting.
	@go test -race -timeout 30m -v ./...

.PHONY: test/coverage
test/go/coverage: prepare-tmp-dir ## Run all tests with coverage report outputting.
	@gocov test ./... | gocov-html > dev/.tmp/coverage-report.html

.PHONY: test/ui
test/ui: ## Test UI js code
	@cd ui && yarn test

.PHONY: test/ui/coverage
test/ui/coverage: ## Test UI js code with coverage report outputting.
	@cd ui && yarn test:coverage

# This is a utility for checking variable definition
# see https://lithic.tech/blog/2020-05/makefile-wildcards/
guard-%:
	@if [[ -z '${${*}}' ]]; then echo 'ERROR: variable $* not set' && exit 1; fi

prepare-tmp-dir:
	@if [[ ! -e dev/.tmp ]]; then mkdir dev/.tmp; fi
