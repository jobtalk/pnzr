
release:
	@echo "[+] releasing $(VERSION)"
	@echo "[+] building"
	@$(MAKE) build
.PHONY: release
	
build:
		@gox -os="linux darwin windows openbsd" ./...
.PHONY: build


# Clean build artifacts.
clean:
	@git clean -f
.PHONY: clean
