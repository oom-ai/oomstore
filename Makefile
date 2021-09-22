VERSION = $(shell (git describe --tags --abbrev=0 2>/dev/null || echo '0.0.0') | tr -d v)
COMMIT = $(shell git rev-parse --short HEAD)

OUT = build/
export GOPROXY = direct
export CGO_ENABLED = 0

.DEFAULT_GOAL := build

.PHONY: info
info:
	@echo "VERSION: $(VERSION)"
	@echo "COMMIT:  $(COMMIT)"

.PHONY: build
build: featctl

.PHONY: featctl
featctl:
	@go build -o $(OUT)featctl -ldflags "-s -w \
		-X $$(go list -m)/version.Version=$(VERSION) \
		-X $$(go list -m)/version.Commit=$(COMMIT) \
		-X $$(go list -m)/version.Built=$$(date -Iseconds)" \
		featctl/*.go

.PHONY: test
test:
	@go test ./...

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	@rm -rf build
