# Output file name
BINARY=wpchg

# Version
VERSION=0.1.0

# Default target
default: build

# Build the binary
build:
	go build -o ./bin/$(BINARY) -ldflags "-X main.version=$(VERSION)" ./cmd/wpchg

# Build for windows
build-win:
	GOOS=windows GOARCH=amd64 go build -o ./bin/$(BINARY).exe -ldflags "-X main.version=$(VERSION)" ./cmd/wpchg

# Clean the binary
clean:
	rm -f ./bin/*

# Install the binary
install:
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/wpchg

# Print version
version:
	@echo v$(VERSION)

# Help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build     Build the binary"
	@echo "  build-win Build the binary for windows"
	@echo "  clean     Clean the binary"
	@echo "  install   Install the binary"
	@echo "  version   Print version number"
	@echo "  help      Show this help message"
	@echo ""