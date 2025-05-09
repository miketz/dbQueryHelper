.DEFAULT_GOAL := build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint: fmt
	golint ./...

.PHONY: vet
vet: fmt
	go vet ./...

# use -ldflags to omit symbol table, debug info, and dwarf symbol table. (smaller binary).
# use -gcflags=-B to eliminate bounds checks
.PHONY: build
build: vet
	go build -gcflags=-B -ldflags="-s -w"
