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

help:
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ { \
		printf "  %-12s %s\n", $$1, $$2 \
	}' $(MAKEFILE_LIST)

tidy:
	$(GO_ENV) $(GO) mod tidy

fmt:
	@find cmd internal -name '*.go' -print0 | xargs -0 gofmt -w

test:
	$(GO_ENV) $(GO) test ./...

vet:
	$(GO_ENV) $(GO) vet ./...

check: fmt tidy vet test

run:
	$(GO_ENV) $(GO) run $(CMD) \
		--root "$(ROOT)" \
		--host "$(HOST)" \
		--port "$(PORT)"

build:
	@mkdir -p $(BIN_DIR)
	$(GO_ENV) $(GO) build -o $(BIN) $(CMD)

clean:
	rm -rf $(BIN_DIR)
