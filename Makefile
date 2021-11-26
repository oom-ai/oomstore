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

.PHONY: codegen
codegen:
	@docker run --rm -v $$(pwd):/oomstore -w /oomstore rvolosatovs/protoc \
		--experimental_allow_proto3_optional \
		-I=proto \
		--go_out=oomd \
		--go-grpc_out=oomd \
		proto/oomd.proto
	@docker run --rm -v $$(pwd):/oomstore -w /oomstore rvolosatovs/protoc \
		--experimental_allow_proto3_optional \
		-I=proto \
		--python_out=sdk/python/src/oomstore/codegen \
		--grpc-python_out=sdk/python/src/oomstore/codegen \
		proto/oomd.proto
	@perl -pi -e s,"import status_pb2","from . import status_pb2",g sdk/python/src/oomstore/codegen/oomd_pb2.py
	@perl -pi -e s,"import oomd_pb2","from . import oomd_pb2",g sdk/python/src/oomstore/codegen/oomd_pb2_grpc.py

.PHONY: build
build: oomctl oomd

.PHONY: oomctl
oomctl:
	$(MAKE) -C oomctl build

.PHONY: oomd
oomd:
	$(MAKE) -C oomd build

.PHONY: test
test:
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

.PHONY: integration-test
integration-test: codegen
	$(MAKE) -C oomctl integration-test
	$(MAKE) -C oomd    integration-test

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	$(MAKE) -C oomctl clean
