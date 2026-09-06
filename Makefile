APP_NAME := serverigh
CMD := ./cmd/serverigh
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP_NAME)

GO ?= go
GOCACHE ?= /tmp/serverigh-go-build
GOMODCACHE ?= /tmp/serverigh-go-mod

ROOT ?= .
HOST ?= 127.0.0.1
PORT ?= 4173

GO_ENV := GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE)

.PHONY: help tidy fmt test vet check run build clean

help: ## Show available commands.
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ { \
		printf "  %-12s %s\n", $$1, $$2 \
	}' $(MAKEFILE_LIST)

tidy: ## Sync Go module metadata.
	$(GO_ENV) $(GO) mod tidy

fmt: ## Format Go source files.
	@find cmd internal web -name '*.go' -print0 | xargs -0 gofmt -w

test: ## Run the Go test suite.
	$(GO_ENV) $(GO) test ./...

vet: ## Run go vet.
	$(GO_ENV) $(GO) vet ./...

check: fmt tidy vet test ## Format, tidy, vet, and test.

run: ## Run the development server.
	$(GO_ENV) $(GO) run $(CMD) \
		--root "$(ROOT)" \
		--host "$(HOST)" \
		--port "$(PORT)"

build: ## Build the serverigh binary.
	@mkdir -p $(BIN_DIR)
	$(GO_ENV) $(GO) build -o $(BIN) $(CMD)

clean: ## Remove build artifacts.
	rm -rf $(BIN_DIR)
