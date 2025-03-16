# Makefile for app-gateway Go project with gqlgen

# Configuration
BINARY_NAME=app-gateway
MAIN_PACKAGE=./commander
GENERATED_DIR=./graph/generated

# Go commands
GQLGEN=go run github.com/99designs/gqlgen
GO_RUN=go run
GO_BUILD=go build
GO_TEST=go test
GO_MOD=go mod
GO_FMT=gofmt

# Default target
.PHONY: all
all: tidy generate build

# Tools management
.PHONY: tools
tools:
	@echo "Creating tools.go file if it doesn't exist..."
	@if [ ! -f tools.go ]; then \
		echo '//go:build tools\n// +build tools\n\npackage tools\n\nimport (\n\t_ "github.com/99designs/gqlgen"\n\t_ "github.com/99designs/gqlgen/cmd"\n)' > tools.go; \
	fi
	$(GO_MOD) tidy

# GraphQL code generation
.PHONY: generate
generate: tools
	@echo "Generating GraphQL code..."
	@mkdir -p $(GENERATED_DIR)
	$(GQLGEN) generate

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	$(GO_BUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	$(GO_RUN) $(MAIN_PACKAGE)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO_TEST) ./... -v

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO_TEST) ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Clean up generated files and binaries
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO_FMT) -w .

# Update dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GO_MOD) tidy

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all            - Run tidy, generate, and build"
	@echo "  tools          - Create tools.go file for managing tool dependencies"
	@echo "  generate       - Generate GraphQL code using gqlgen"
	@echo "  build          - Build the application binary"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove binary and coverage files"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy up Go module dependencies"
	@echo "  help           - Display this help message"

