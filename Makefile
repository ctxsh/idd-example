APP := deployment-api
BIN_DIR := $(CURDIR)/bin
CACHE_DIR := $(CURDIR)/.cache
GOCACHE := $(CACHE_DIR)/go-build
GOMODCACHE := $(CACHE_DIR)/gomod
GOLANGCI_LINT_CACHE := $(CACHE_DIR)/golangci-lint
GO := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go
GOLANGCI_LINT_VERSION ?= v2.11.4
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint
PKGS := ./...

.DEFAULT_GOAL := help

.PHONY: help deps build test lint fmt tidy clean

help:
	@echo "Targets:"
	@echo "  deps   install local development tools into ./bin"
	@echo "  build  build $(APP) into ./bin"
	@echo "  test   run Go tests"
	@echo "  lint   run golangci-lint"
	@echo "  fmt    format Go files"
	@echo "  tidy   tidy Go modules"
	@echo "  clean  remove local build and tool outputs"

$(BIN_DIR):
	mkdir -p $@

$(GOLANGCI_LINT): | $(BIN_DIR)
	GOBIN=$(BIN_DIR) $(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

deps: $(GOLANGCI_LINT)

build: | $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP) .

test:
	$(GO) test $(PKGS)

lint: $(GOLANGCI_LINT)
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) fmt --diff $(PKGS)
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) run $(PKGS)

fmt: $(GOLANGCI_LINT)
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) fmt $(PKGS)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR) $(CACHE_DIR)
