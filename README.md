# Object Cloud

`object-cloud` is a lightweight, high-performance file storage microservice written in Go. It enables backend applications to generate secure, presigned upload tokens so client browsers can upload files directly to the storage server, keeping heavy file traffic off your main application backend.

For version 1, it operates as a streamlined, single-node service that streams file uploads straight to local disk and serves them publicly, providing a foundational storage primitive that can be shared across multiple projects. The architecture is designed to be simple today, establishing the groundwork to eventually evolve into a partitioned, highly available distributed cluster in future versions.

---

## Prerequisites

- [Go](https://go.dev/doc/install) (1.24+ recommended)
- `git`

---

## Getting Started

### Clone Repository

```bash
git clone git@github.com:kzielonka/object-cloud.git
cd object-cloud
```

### Download Dependencies

```bash
go mod download
```

---

## Testing

Run the test suite using standard Go tooling:

```bash
# Run all tests in all packages
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test -v ./internal/object/...

# Run tests with race detection enabled
go test -race ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## Code Quality & Static Analysis

```bash
# Format code
go fmt ./...

# Run standard Go static analysis
go vet ./...

# Tidy module dependencies
go mod tidy
```

---

## Building

To compile the packages:

```bash
# Verify compilation across all packages
go build ./...

# When a main entrypoint (e.g. cmd/server/main.go) is added:
# go build -o bin/object-cloud ./cmd/server
```

---

## Project Structure

```text
.
├── README.md             # Project documentation
├── go.mod                # Go module definition
└── internal/             # Private application and library code
    └── object/           # Core object storage logic and interfaces
        ├── object.go     # Store implementation and Storage interface
        └── object_test.go# Unit tests for Store
```
