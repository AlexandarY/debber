GO?=$(shell which go)
BUILD_OS?=linux
ARCH?=amd64
BUILD_DIR=target
BUILD_NAME=debber

all: clean build
build:
	@echo "Doing new production build"
	GOOS=$(BUILD_OS) GOARCH=$(ARCH) go build -ldflags "-s -w" -o $(BUILD_DIR)/$(BUILD_NAME) cmd/debber/main.go
clean:
	@echo "Cleaning old build"
	rm -rf $(BUILD_DIR)
help:
	@echo "Available commands:"
	@echo " * clean"
	@echo "   Cleaning old builds"
	@echo " * build"
	@echo "   Build a release version"
	@echo " * help"
	@echo "   Show this message"
