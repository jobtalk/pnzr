HASH := $(shell git rev-parse --short HEAD)
GOVERSION := $(shell go version)
VERSION ?= $(HASH)
DATE := $(shell LC_ALL=c date)
BUILD_OS := $(shell uname)

build:
	@mkdir -p bin/darwin
	@mkdir -p bin/linux
	@echo "build linux binary"
	@GOOS=linux GOARCH=amd64 go build -ldflags '-X "main.VERSION=$(VERSION)" -X "main.BUILD_DATE=$(DATE)" -X "main.BUILD_OS=$(BUILD_OS)"' -o bin/linux/pnzr
	@echo "build darwin binary"
	@GOOS=darwin GOARCH=amd64 go build -ldflags '-X "main.VERSION=$(VERSION)" -X "main.BUILD_DATE=$(DATE)" -X "main.BUILD_OS=$(BUILD_OS)"' -o bin/darwin/pnzr
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
	@rm -rf bin
.PHONY: clean
