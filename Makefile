.PHONY: all build install test clean fmt lint demo help

BINARY_NAME=zap
BINARY_PATH=./cmd/zap
INSTALL_PATH=/usr/local/bin

all: build

## build: build the binary
build:
	@echo "building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(BINARY_PATH)

## install: install the binary to system
install: build
	@echo "installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "installed successfully!"

## test: run all tests
test:
	@echo "running tests..."
	@go test -v -race -cover ./...

## test-coverage: run tests with coverage report
test-coverage:
	@echo "running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "coverage report generated: coverage.html"

## clean: remove build artifacts
clean:
	@echo "cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@echo "clean complete!"

## fmt: format code
fmt:
	@echo "formatting code..."
	@go fmt ./...

## lint: run linter
lint:
	@echo "running linter..."
	@golangci-lint run

## run: build and run (for testing)
run: build
	@./$(BINARY_NAME)

## demo: generate demo GIF with VHS
demo: build
	@echo "generating demo GIF..."
	@command -v vhs >/dev/null 2>&1 || { echo "error: vhs is not installed. install with: brew install vhs"; exit 1; }
	@vhs assets/demo.tape
	@echo "demo generated: assets/demo.gif"

## demo-compress: compress demo GIF with gifsicle
demo-compress:
	@echo "compressing demo GIF..."
	@command -v gifsicle >/dev/null 2>&1 || { echo "error: gifsicle is not installed. install with: brew install gifsicle"; exit 1; }
	@gifsicle -O3 --colors 256 assets/demo.gif -o assets/demo.gif
	@echo "demo compressed!"

## help: show this help message
help:
	@echo "available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/ /'
