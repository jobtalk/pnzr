HASH := `git rev-parse --short HEAD`
VERSION ?= `git show --quiet --pretty=format:"%cd" HEAD | sed 's/ /_/g' | sed 's/+0900//g' | sed 's/:/-/g'`$(HASH)

build:
	@mkdir -p bin/darwin
	@mkdir -p bin/linux
	@echo "build linux binary"
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/thor
	@tar -zcvf bin/linux/thor-linux-amd64.tar.gz bin/linux/thor
	@echo "build darwin binary"
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin/thor
	@tar -zcvf bin/darwin/thor-darwin-amd64.tar.gz bin/darwin/thor
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
	@rm -rf bin
.PHONY: clean
