.PHONY: build install test lint run clean demo-gifs

BINARY=slk
BUILD_DIR=bin

build:
	go build -ldflags="-s -w" -trimpath -o $(BUILD_DIR)/$(BINARY) ./cmd/slk

# fork: install replaces PREFIX/bin/slk (default /usr/local/bin, ahead of
# ~/go/bin on PATH) instead of `go install`, and stamps the fork's git
# version so `slk --version` shows which build is running.
PREFIX ?= /usr/local
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)-fork
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
STAMP = -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

install:
	go build -ldflags="-s -w $(STAMP)" -trimpath -o $(BUILD_DIR)/$(BINARY) ./cmd/slk
	install -d $(PREFIX)/bin
	install -m 0755 $(BUILD_DIR)/$(BINARY) $(PREFIX)/bin/$(BINARY)
	@echo "installed $(PREFIX)/bin/$(BINARY) ($(VERSION))"

test:
	go test ./... -v -race

lint:
	golangci-lint run ./...

run: build
	./$(BUILD_DIR)/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

# Renders docs/assets/demo/*.gif from their VHS tapes (settings.tape is
# shared, not rendered). Needs vhs: run inside `nix develop`, or see
# https://github.com/charmbracelet/vhs. Not run in CI.
demo-gifs: build
	@command -v vhs >/dev/null || { echo "vhs not found: run 'nix develop' or install charmbracelet/vhs"; exit 1; }
	@for tape in docs/assets/demo/*.tape; do \
		case $$tape in */settings.tape) continue;; esac; \
		echo "vhs $$tape"; vhs $$tape || exit 1; \
	done
