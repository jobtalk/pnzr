HASH := `git rev-parse --short HEAD`
VERSION ?= `git show --quiet --pretty=format:"%cd" HEAD | sed 's/ /_/g' | sed 's/+0900//g' | sed 's/:/-/g'`$(HASH)

release:
	@echo "[+] releasing $(VERSION)"
	@echo "[+] building"
	@$(MAKE) build
	@echo "[+] comitting"
	@git add thor_darwin_amd64
	@git add thor_linux_amd64
	@git tag -a $(VERSION) -m $(VERSION)
	@git push origin $(VERSION)
	@rm -rf thor_darwin_amd64
	@rm -rf thor_linux_amd64
	@echo "[+] complete"
.PHONY: release
	
build:
		@gox -os="linux darwin" -arch="amd64" `glide novendor`
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
.PHONY: clean
