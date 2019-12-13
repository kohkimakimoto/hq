.PHONY:help
.DEFAULT_GOAL := help

export GO111MODULE := off
export PATH := $(CURDIR)/.go-packages/bin:$(PATH)


# This is a magic code to output help message at default
# see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY:fmt
fmt: ## Run `go fmt`
	go fmt $$(go list ./... | grep -v vendor)

.PHONY:test
test: ## Run all tests
	go test -cover $$(go list ./... | grep -v vendor)

.PHONY:testv
testv: ## Run all tests with verbose outputing.
	go test -v -cover $$(go list ./... | grep -v vendor)

.PHONY:testcov
testcov:
	gocov test $$(go list ./... | grep -v vendor) | gocov-html > coverage-report.html

.PHONY:dev
dev: ## Build dev binary
	@bash -c $(CURDIR)/build/scripts/dev.sh

.PHONY:dist
dist: bindata ## Build dist binaries
	@bash -c $(CURDIR)/build/scripts/dist.sh

.PHONY:clean
clean: ## Clean build outputs
	@bash -c $(CURDIR)/build/scripts/clean.sh

.PHONY:bindata
bindata: ## bindata
	yarn run clean && yarn run build
	go-bindata -o ./res/ui/bindata.go -ignore bindata.go -prefix res/ui -pkg views res/ui/...
	go-bindata -o ./res/views/bindata.go -ignore bindata.go -prefix res/views -pkg views res/views/...

.PHONY:devbindata
devbindata: ## devbindata
	yarn run clean && yarn run dev
	go-bindata -o ./res/ui/bindata.go -ignore bindata.go -prefix res/ui -pkg views res/ui/...
	go-bindata -o ./res/views/bindata.go -ignore bindata.go -prefix res/views -pkg views res/views/...

.PHONY:packaging
packaging: ## Create packages (now support RPM only)
	@bash -c $(CURDIR)/build/scripts/packaging.sh


.PHONY: installtools
installtools: ## Install dev tools
	GOPATH=$(CURDIR)/.go-packages && \
      go get -u github.com/mattn/go-bindata/... && \
      go get -u github.com/golang/dep/cmd/dep && \
      go get -u github.com/mitchellh/gox && \
      go get -u github.com/axw/gocov/gocov && \
      go get -u gopkg.in/matm/v1/gocov-html

.PHONY:deps
deps: ## Install dependences.
	PATH=$(CURDIR)/.go-packages/bin:${PATH} && dep ensure

.PHONY:updatedeps
updatedeps: ## update all dependences.
	rm -rf Gopkg.*
	rm -rf vendor
	dep init
	dep ensure



