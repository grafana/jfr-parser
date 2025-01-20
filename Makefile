GO ?= go

.phony: generate-types
generate-types:
	$(GO) run ./gen
	$(GO) fmt ./parser/types


TEST_PACKAGES := ./... ./pprof/...

.PHONY: test
test:
	go test -race $(shell go list $(TEST_PACKAGES))

# Test if dependencies are up to date, without go workspaces pinning
.PHONY: test-dependency
test-dependency:
	cd ./pprof && GOWORK=off go test ./...
