HASH := `git rev-parse --short HEAD`
VERSION ?= `git show --quiet --pretty=format:"%cd" HEAD | sed 's/ /_/g' | sed 's/+0900//g' | sed 's/:/-/g'`$(HASH)

build:
	@mkdir -p bin/darwin
	@mkdir -p bin/linux
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/thor
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin/thor
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
	@rm -rf bin
.PHONY: clean
