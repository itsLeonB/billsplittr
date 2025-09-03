.PHONY: http http-hot lint test test-verbose test-coverage test-coverage-html test-clean security

http:
	go run cmd/http/main.go

http-hot:
	@echo "🚀 Starting HTTP server with hot reload..."
	air --build.cmd "go build -o bin/http cmd/http/main.go" --build.bin "./bin/http"

lint:
	golangci-lint run ./...

test:
	@echo "Running all tests..."
	go test ./internal/test/...

test-verbose:
	@echo "Running all tests with verbose output..."
	go test -v ./internal/test/...

test-coverage:
	@echo "Running all tests with coverage report..."
	go test -v -cover -coverprofile=coverage.out -coverpkg=./internal/... ./internal/test/...

test-coverage-html:
	@echo "Running all tests and generating HTML coverage report..."
	go test -v -cover -coverprofile=coverage.out -coverpkg=./internal/... ./internal/test/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-clean:
	@echo "Cleaning test cache and running tests..."
	go clean -testcache && go test -v ./internal/test/...
