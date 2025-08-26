# Meal Planner Go

A simple meal planner and shopping list application built with Go and HTMX.

## Features

- Weekly meal planning
- Shopping list management
- Drag-and-drop reordering of shopping list items
- Marking items as complete

## Tech Stack

- **Backend**: Go
- **Database**: MongoDB
- **Frontend**: HTMX, TailwindCSS, SortableJS

## Project Structure

```
.
├── cmd/
│   └── server/            # Application entry point
│       └── main.go        # Server initialization and configuration
├── pkg/
│   ├── db/                # Database layer
│   │   └── mongodb.go     # MongoDB connection and operations
│   ├── handlers/          # HTTP handlers
│   │   └── handlers.go    # Route handlers implementation
│   ├── models/            # Data models
│   │   └── models.go      # Application data structures
│   └── templates/         # HTML templates
│       └── index.html     # Main application template
├── public/                # Static assets
│   ├── css/               # CSS files
│   │   └── index.css      # Application styles
│   └── js/                # JavaScript files
│       └── htmx.min.js    # HTMX library
├── go.mod                 # Go module file
└── go.sum                 # Go dependencies checksum
```

## Getting Started

### Prerequisites

- Go (1.22+)
- MongoDB Atlas account (or local MongoDB instance)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mealplannergo
   ```

2. **Set up environment variables**
   ```bash
   # Copy the example environment file
   cp .env.example .env
   
   # Edit .env and add your MongoDB Atlas connection string
   # GO_SHOPPING_MONGO_ATLAS_URI=mongodb+srv://your-connection-string
   # GO_ENV=development
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Set up test database (IMPORTANT for development)**
   ```bash
   # This copies your production data to the test database
   # You only need to do this once, or when you want fresh test data
   make copy-prod-to-test
   ```

### Running the Application

**For Development (Recommended):**
```bash
# Runs with test database (GoShopping-test) - safe for development
make dev
```

**For Production:**
```bash
# ⚠️ WARNING: Uses production database (GoShopping)
make prod
```

**Alternative ways to run:**
```bash
# Using go run directly (development mode)
GO_ENV=development go run ./cmd/server

# Using go run directly (production mode)  
GO_ENV=production go run ./cmd/server
```

The application will be available at http://localhost:8080.

### Environment Modes

- **Development Mode** (`GO_ENV=development`): 
  - Uses `GoShopping-test` database
  - Safe for testing new features
  - Default mode when using `.env` file

- **Production Mode** (`GO_ENV=production`):
  - Uses `GoShopping` database
  - ⚠️ **USE WITH CAUTION** - modifies production data
  - Requires explicit confirmation

### Database Management

```bash
# Copy production data to test database
make copy-prod-to-test

# Run server in development mode (test database)
make dev

# Run server in production mode (production database)
make prod
```

## Testing

The project includes a comprehensive test suite with unit tests and integration tests.

```bash
# Run all tests
make test

# Run tests with coverage report
make coverage
```

Test coverage is a priority for this project. We aim for 100% test coverage to ensure code quality and prevent regressions.

### Test Structure

- **Unit Tests**: Test individual packages in isolation
  - `pkg/models/models_test.go`: Tests for data models
  - `pkg/handlers/handlers_test.go`: Tests for HTTP handlers using mocked DB
  - `pkg/db/mongodb_test.go`: Tests for database operations

- **Integration Tests**: Test API endpoints
  - `cmd/server/main_test.go`: Tests for API routes

The test suite uses the following libraries:
- `testify`: For assertions and mocks
- `httptest`: For HTTP testing