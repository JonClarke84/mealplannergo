.PHONY: test build run clean coverage

# Default Go build flags
GOFLAGS := -v

# Default target
all: test build

# Build the application
build:
	@echo "Building mealplannergo..."
	@go build $(GOFLAGS) -o bin/app ./cmd/server

# Run the application
run: build
	@echo "Running mealplannergo..."
	@./bin/app

# Run all tests
test:
	@echo "Running tests..."
	@go test ./... -v

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f bin/app
	@rm -f coverage.out coverage.html

# Install test dependencies
test-deps:
	@echo "Installing test dependencies..."
	@go get -u github.com/stretchr/testify/assert
	@go get -u github.com/stretchr/testify/mock
	@go get -u github.com/golang/mock/mockgen