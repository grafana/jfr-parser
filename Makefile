GO ?= go

.phony: generate-types
generate-types:
	$(GO) run ./internal/cmd/gen
	$(GO) fmt ./parser/types


TEST_PACKAGES := ./...

.PHONY: test
test:
	go test -v -race $(shell go list $(TEST_PACKAGES))

