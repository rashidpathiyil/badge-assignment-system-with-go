.PHONY: test test-unit test-integration test-badge test-event test-condition test-all

# Default task when running make
all: test

# Run all tests
test-all: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test ./internal/... -v

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test ./tests/integration/... -v

# Run badge tests
test-badge:
	@echo "Running badge tests..."
	go test ./tests/integration/badge/... -v

# Run event tests
test-event:
	@echo "Running event tests..."
	go test ./tests/integration/event/... -v

# Run condition tests
test-condition:
	@echo "Running condition tests..."
	go test ./tests/integration/condition/... -v

# Run API tests
test-api:
	@echo "Running API tests..."
	go test ./tests/integration/api/... -v

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./internal/... ./tests/... -v

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	go test -bench=. ./internal/...

# Generate test coverage report
test-coverage:
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean test cache and artifacts
clean:
	@echo "Cleaning test cache and artifacts..."
	go clean -testcache
	rm -f coverage.out coverage.html

# Setup test environment
setup-test:
	@echo "Setting up test environment..."
	# Add commands to set up test database or other dependencies here

# Help target
help:
	@echo "Available targets:"
	@echo "  test-all        - Run all tests"
	@echo "  test-unit       - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-badge      - Run badge tests"
	@echo "  test-event      - Run event tests"
	@echo "  test-condition  - Run condition tests"
	@echo "  test-api        - Run API tests"
	@echo "  test-race       - Run tests with race detection"
	@echo "  benchmark       - Run benchmarks"
	@echo "  test-coverage   - Generate test coverage report"
	@echo "  clean           - Clean test cache and artifacts"
	@echo "  setup-test      - Setup test environment"
	@echo "  help            - Show this help message" 
