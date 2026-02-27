.PHONY: help build test run clean fmt vet

help:
	@echo "Phone Number Normalizer - Available targets:"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  run         - Build and run the application"
	@echo "  fmt         - Format code with go fmt"
	@echo "  vet         - Analyze code with go vet"
	@echo "  clean       - Remove build artifacts"

build:
	go build -o phone.exe

test:
	go test -v ./...

run: build
	./phone.exe

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f phone.exe
	go clean
