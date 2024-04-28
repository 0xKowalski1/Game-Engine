EXAMPLE ?= rotating-cubes

.PHONY: dev
dev:
	go run ./examples/$(EXAMPLE)/main.go

