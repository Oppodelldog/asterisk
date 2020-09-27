GOCMD=go
GOBUILD=$(GOCMD) build -race
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

setup: ## Install tools
	$(GOGET) -u golang.org/x/tools/...
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s v1.30.0
	mv bin/golangci-lint $(GOPATH)/bin/golangci-lint && rm -rf bin

lint: ## Run the linters
	golangci-lint run

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

tests:
	cd test && $(GOTEST) -race -covermode=atomic -coverpkg=../../asterisk/... ./...
# Self-Documented Makefile see https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.DEFAULT_GOAL := help