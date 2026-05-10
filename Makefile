# goz — Makefile
# Run `make help` to list all targets. Group headers come from `##@` comments,
# per-target docs come from inline `## description` after the target name.

BINARY  := goz
CMD     := ./cmd/goz
PKG     := ./...
BIN_DIR := bin

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

LDFLAGS := -s -w -X goz/internal/version.Version=$(VERSION)

PREVIEW_SIZE ?= 140x40
SNAPSHOT     := internal/tui/testdata/snapshot_$(PREVIEW_SIZE).txt

GO          ?= go
GOLANGCI    ?= golangci-lint
GOVULNCHECK ?= govulncheck
AIR         ?= air

.DEFAULT_GOAL := help

.PHONY: help \
        build run install dev clean \
        test test-race test-pkg cover cover-html bench \
        fmt fmt-check vet lint lint-fix tidy mod-verify \
        check ci audit tools \
        preview snapshot

# ============================================================================

help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} \
		/^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 } \
		/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Build & run

build: ## Build the binary into bin/goz
	@$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(CMD)
	@echo "→ $(BIN_DIR)/$(BINARY)  ($(VERSION))"

run: ## Run the TUI from source
	@$(GO) run $(CMD)

install: ## Install the binary to $$GOBIN (or $$GOPATH/bin)
	@$(GO) install -trimpath -ldflags '$(LDFLAGS)' $(CMD)

dev: ## Hot-reload via air (install with `make tools`)
	@command -v $(AIR) >/dev/null 2>&1 || { echo "air not installed; run: make tools" >&2; exit 1; }
	@$(AIR)

clean: ## Remove build + coverage artifacts
	@rm -rf $(BIN_DIR) coverage.out coverage.html

##@ Test

test: ## Run all tests
	@$(GO) test -count=1 $(PKG)

test-race: ## Run all tests with the race detector
	@$(GO) test -race -count=1 $(PKG)

test-pkg: ## Run tests verbosely in PKG (e.g. make test-pkg PKG=./internal/store)
	@$(GO) test -count=1 -v $(PKG)

cover: ## Generate coverage.out and print summary
	@$(GO) test -count=1 -covermode=atomic -coverprofile=coverage.out $(PKG)
	@$(GO) tool cover -func=coverage.out | tail -n 1

cover-html: cover ## Render coverage to coverage.html
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "→ coverage.html"

bench: ## Run benchmarks
	@$(GO) test -bench=. -benchmem -run=^$$ $(PKG)

##@ Quality

fmt: ## Format with gofmt -s
	@gofmt -s -w .

fmt-check: ## Fail if anything needs gofmt (no writes)
	@out=$$(gofmt -s -l .); \
	if [ -n "$$out" ]; then echo "unformatted files:"; echo "$$out"; exit 1; fi

vet: ## go vet
	@$(GO) vet $(PKG)

lint: ## golangci-lint
	@$(GOLANGCI) run $(PKG)

lint-fix: ## golangci-lint --fix
	@$(GOLANGCI) run --fix $(PKG)

tidy: ## go mod tidy
	@$(GO) mod tidy

mod-verify: ## Verify go.sum and that go.mod is tidy
	@$(GO) mod verify
	@diff=$$($(GO) mod tidy -diff 2>/dev/null); \
	if [ -n "$$diff" ]; then echo "go.mod is not tidy; run: make tidy"; exit 1; fi

##@ Checks

check: fmt-check vet test ## Quick pre-commit check (fmt-check + vet + test)

ci: fmt-check vet lint test-race build ## Full CI pipeline

audit: ## Scan for known vulnerabilities (govulncheck)
	@command -v $(GOVULNCHECK) >/dev/null 2>&1 || { echo "govulncheck not installed; run: make tools" >&2; exit 1; }
	@$(GOVULNCHECK) $(PKG)

##@ Tools & TUI

tools: ## Install dev tools (golangci-lint, govulncheck, air)
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	@$(GO) install github.com/air-verse/air@latest

preview: build ## Render a frame with --preview (override with PREVIEW_SIZE=120x32)
	@$(BIN_DIR)/$(BINARY) --preview $(PREVIEW_SIZE)

snapshot: build ## Regenerate the TUI snapshot at testdata/snapshot_<size>.txt
	@$(BIN_DIR)/$(BINARY) --preview $(PREVIEW_SIZE) > $(SNAPSHOT)
	@echo "→ $(SNAPSHOT)"
