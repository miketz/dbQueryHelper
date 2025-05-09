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
# use -trimpath to remove filepaths. (smaller binary).
# use -gcflags=-B to eliminate bounds checks
.PHONY: build
build: vet
	CGO_ENABLED=0 go build -gcflags=-B -ldflags="-s -w" -trimpath

# a typical go build. nothing stripped out.
buildFat: vet
	go build

# Make things small, but keep the bounds checks for safety.
buildSmall: vet
	go build -ldflags="-s -w" -trimpath

# disable optimizations to help Delve debugger.
buildDebug: vet
	go build -gcflags="all=-N -l"

