# Go parameters
PROJECT_NAME := $(shell echo $${PWD\#\#*/})
PKG_LIST := $(shell go list ./...)
GO_FILES := $(shell find . -name '*.go' | grep -v _test.go)

.PHONY: all
all: lint install tags

.PHONY: install
install: ## Run install
	@go install ./... && echo Installed `date` && echo

.PHONY: dep
dep:
	@go get -u ./...

.PHONY: lintAll
lintAll: # run golangci-lint
	@golangci-lint run

.PHONY: lint
lint: ## Run lint
	@golangci-lint run --fast

.PHONY: test
test: ## Run unittests
	@go test -short ${PKG_LIST}

.PHONY: race
race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

.PHONY: msan
msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

.PHONY: build
build: ## Build the binary file
	@go build -i -v

.PHONY: clean
clean: ## Remove previous build
	@go clean ./...

.PHONY: linux
linux:
	@env GOOS=linux GOARCH=amd64 go build -v -o $(CURDIR)/build/$(PROJECT_NAME)

.PHONY: lines
lines:
	@find . -name "*.go" -exec wc -l {} \+

.PHONY: tags
tags:
	@-find . -name '*.go' -exec gotags {} \+ > tags

.PHONY: help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
