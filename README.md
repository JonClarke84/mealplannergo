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

- Go (1.16+)
- MongoDB

### Environment Variables

Set the following environment variable with your MongoDB connection string:

```
export GO_SHOPPING_MONGO_ATLAS_URI="mongodb+srv://your-connection-string"
```

### Running the Application

```bash
# Build and run the application
go build -o app ./cmd/server
./app

# Or run directly
go run ./cmd/server
```

The application will be available at http://localhost:8080.

## Development

```bash
# Run with hot reload using air (if installed)
air -c .air.toml

# Or use go run for development
go run ./cmd/server
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