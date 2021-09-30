VERSION = $(shell (git describe --tags --abbrev=0 2>/dev/null || echo '0.0.0') | tr -d v)
COMMIT = $(shell git rev-parse --short HEAD)

OUT = build/
# export GOPROXY = direct
# export CGO_ENABLED = 0

.DEFAULT_GOAL := build

.PHONY: info
info:
	@echo "VERSION: $(VERSION)"
	@echo "COMMIT:  $(COMMIT)"

.PHONY: build
build: featctl

.PHONY: featctl
featctl:
	$(MAKE) -C featctl build

.PHONY: test
test:
	@go test ./...

.PHONY: integration-test
integration-test:
	$(MAKE) -C featctl integration-test

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	$(MAKE) -C featctl clean
