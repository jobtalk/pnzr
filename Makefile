HASH := $(shell git rev-parse --short HEAD)
GOVERSION := $(shell go version)
VERSION ?= $(HASH)
DATE := $(LC_ALL=c shell date)
BUILD_OS := $(shell uname)

build:
	@mkdir -p bin/darwin
	@mkdir -p bin/linux
	@echo "build linux binary"
	@GOOS=linux GOARCH=amd64 go build -ldflags '-X "main.VERSION=$(VERSION)" -X "main.BUILD_DATE=$(DATE)" -X "main.BUILD_OS=$(BUILD_OS)"' -o bin/linux/thor
	@tar -zcvf bin/linux/thor-linux-amd64.tar.gz bin/linux/thor
	@echo "build darwin binary"
	@GOOS=darwin GOARCH=amd64 go build -ldflags '-X "main.VERSION=$(VERSION)" -X "main.BUILD_DATE=$(DATE)" -X "main.BUILD_OS=$(BUILD_OS)"' -o bin/darwin/thor
	@tar -zcvf bin/darwin/thor-darwin-amd64.tar.gz bin/darwin/thor
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
	@rm -rf bin
.PHONY: clean
