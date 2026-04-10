# List all available recipes
default:
    @just --list

# Run all tests and print total coverage
test:
    gotestsum --format short --no-summary=skipped -- -race -count=1 -coverprofile=coverage.out ./...
    @go tool cover -func=coverage.out | grep "^total"

# Run golangci-lint
lint:
    golangci-lint run ./...

# Format source code with gofmt and goimports
fmt:
    gofmt -l -w .
    goimports -l -w .
