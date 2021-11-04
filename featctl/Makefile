VERSION = $(shell (git describe --tags --abbrev=0 2>/dev/null || echo '0.0.0') | tr -d v)
COMMIT = $(shell git rev-parse --short HEAD)

PLATFORMS := darwin/amd64 linux/amd64 windows/amd64
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

OUT = build/

.DEFAULT_GOAL := build

.PHONY: info
info:
	@echo "VERSION: $(VERSION)"
	@echo "COMMIT:  $(COMMIT)"

.PHONY: build
build:
	@go build -o $(OUT)featctl -ldflags "-s -w \
		-X $$(go list -m)/version.Version=$(VERSION) \
		-X $$(go list -m)/version.Commit=$(COMMIT) \
		-X $$(go list -m)/version.Built=$$(date -u +%FT%T+00:00)"

release: $(PLATFORMS)

$(PLATFORMS):
	@ GOOS=$(os) GOARCH=$(arch) go build -o $(OUT)featctl-$(os)-$(arch) -ldflags "-s -w \
		-X $$(go list -m)/version.Version=$(VERSION) \
		-X $$(go list -m)/version.Commit=$(COMMIT) \
		-X $$(go list -m)/version.Built=$$(date -u +%FT%T+00:00)"

.PHONY: release $(PLATFORMS)

.PHONY: test
test:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: integration-test
integration-test: build
	@./test/test.sh

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	@rm -rf build
