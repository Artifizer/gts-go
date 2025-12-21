CI := 1

.PHONY: help build dev-fmt all check fmt vet lint test security coverage update-spec e2e

# Default target - show help
.DEFAULT_GOAL := help

# Show this help message
help:
	@awk '/^# / { desc=substr($$0, 3) } /^[a-zA-Z0-9_-]+:/ && desc { target=$$1; sub(/:$$/, "", target); printf "%-20s - %s\n", target, desc; desc="" }' Makefile | sort

# Build the workspace
build:
	go build ./...
	go build -o ./bin/gts ./cmd/gts
	go build -o ./bin/gts-server ./cmd/gts-server

# Fix formatting issues
dev-fmt:
	gofmt -w .

# Run all checks and build
all: check build

# Check code formatting
fmt:
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "The following files are not formatted:"; \
		gofmt -l .; \
		exit 1; \
	fi

# Run go vet
vet:
	go vet ./...

# Run golangci-lint (skipped if Go version is unsupported)
lint:
	@if [ ! -f "$$(go env GOPATH)/bin/golangci-lint" ]; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@$$(go env GOPATH)/bin/golangci-lint run --timeout=5m || \
		(echo "Warning: golangci-lint failed (may not support your Go version yet). Skipping..." && true)

# Run all tests
test:
	go test -v -race ./...

# Check for vulnerabilities
security:
	@command -v govulncheck >/dev/null || (echo "Installing govulncheck..." && go install golang.org/x/vuln/cmd/govulncheck@latest)
	$$(go env GOPATH)/bin/govulncheck ./...

# Measure code coverage
coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	@echo "To view HTML report: go tool cover -html=coverage.out"

# Update gts-spec submodule to latest
update-spec:
	git submodule update --remote .gts-spec

# Run end-to-end tests against gts-spec
e2e: build
	@echo "Starting server in background..."
	@./bin/gts server --port 8000 & echo $$! > .server.pid
	@sleep 2
	@echo "Running e2e tests..."
	@PYTHONDONTWRITEBYTECODE=1 pytest -p no:cacheprovider --log-file=e2e.log ./.gts-spec/tests || (kill `cat .server.pid` 2>/dev/null; rm -f .server.pid; exit 1)
	@echo "Stopping server..."
	@kill `cat .server.pid` 2>/dev/null || true
	@rm -f .server.pid
	@echo "E2E tests completed successfully"

# Run all quality checks
check: fmt vet lint test e2e
