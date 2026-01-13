# eval-hub-backend-svc

A Go API service built with net/http.

## Getting Started

### Prerequisites

- Go 1.21 or higher

### Running the Service

#### Using Make (Recommended)

1. Install dependencies:
```bash
make install-deps
```

2. Run the server:
```bash
make run
```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable:

```bash
PORT=3000 make run
```

#### Using Go directly

1. Install dependencies:
```bash
go mod download
```

2. Run the server:
```bash
go run cmd/eval_hub/main.go
```

The server will start on port 8080 by default. You can change this by setting the `PORT` environment variable:

```bash
PORT=3000 go run cmd/eval_hub/main.go
```

### API Endpoints

- `GET /api/v1/health` - Health check endpoint
- `GET /api/v1/status` - Service status endpoint
- `GET /metrics` - Prometheus metrics endpoint
- `GET /openapi.yaml` - OpenAPI 3.1.0 specification
- `GET /docs` - Interactive API documentation (Swagger UI)

### Building

#### Using Make

Build the binary:
```bash
make build
```

Run the binary:
```bash
./bin/eval-hub-backend-svc
```

#### Using Go directly

Build the binary:
```bash
go build -o bin/eval-hub-backend-svc ./cmd/eval_hub
```

Run the binary:
```bash
./bin/eval-hub-backend-svc
```

### Makefile Targets

The project includes a Makefile with common development tasks:

- `make help` - Display all available targets
- `make clean` - Remove build artifacts
- `make build` - Build the binary
- `make run` - Run the application
- `make lint` - Lint the code (runs go vet)
- `make fmt` - Format code with go fmt (NOTE: converts to tabs per Go standard)
- `make vet` - Run go vet
- `make test` - Run unit tests
- `make test-fvt` - Run FVT (Functional Verification Tests) using godog
- `make test-all` - Run all tests (unit + FVT)
- `make test-coverage` - Run unit tests with coverage report
- `make install-deps` - Install and tidy dependencies

## Project Structure

This project follows the [standard Go project layout](https://github.com/golang-standards/project-layout):

```
eval-hub-backend-svc/
├── cmd/
│   └── eval_hub/          # Main application entry point
│       └── main.go
├── internal/               # Private application code
│   ├── handlers/          # HTTP handlers
│   │   └── handlers.go
│   ├── metrics/           # Prometheus metrics
│   │   ├── metrics.go
│   │   └── middleware.go
│   └── server/            # Server setup and configuration
│       └── server.go
├── api/                 # API specifications
│   └── openapi.yaml     # OpenAPI 3.1.0 specification
├── tests/               # Test files
│   └── features/        # BDD feature files and step definitions
│       ├── health.feature
│       ├── status.feature
│       ├── metrics.feature
│       ├── suite_test.go
│       └── step_definitions_test.go
├── go.mod
└── README.md
```

## Testing

The project includes comprehensive test coverage:

### Unit Tests

Unit tests are located alongside the code in `*_test.go` files:
- `internal/handlers/handlers_test.go` - Handler unit tests
- `internal/metrics/middleware_test.go` - Metrics middleware tests
- `internal/server/server_test.go` - Server unit tests

Run unit tests:
```bash
make test
```

### FVT (Functional Verification Tests)

FVT tests use [godog](https://github.com/cucumber/godog) for BDD-style testing:
- Feature files in `tests/features/*.feature`
- Step definitions in `tests/features/step_definitions_test.go`

Run FVT tests:
```bash
make test-fvt
```

Run all tests:
```bash
make test-all
```