.PHONY: test build run clean coverage dev prod copy-prod-to-test

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

# Run in development mode (uses GoShopping-test database)
dev: build
	@echo "Starting server in development mode..."
	@echo "Using test database: GoShopping-test"
	@GO_ENV=development ./bin/app

# Run in production mode (uses GoShopping database)
prod: build
	@echo "⚠️  WARNING: Starting server in PRODUCTION mode!"
	@echo "Using production database: GoShopping"
	@echo "Press Ctrl+C within 5 seconds to cancel..."
	@sleep 5
	@GO_ENV=production ./bin/app

# Copy production data to test database
copy-prod-to-test:
	@echo "Copying production data to test database..."
	@echo "⚠️  This will overwrite any existing data in GoShopping-test!"
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@go run scripts/migrate.go