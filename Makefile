GO ?= go

.phony: generate-types
generate-types:
	$(GO) run ./gen
	$(GO) fmt ./parser/types