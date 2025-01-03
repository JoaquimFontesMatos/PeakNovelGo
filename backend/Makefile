# Variables
GO := go
TEST_DIR := ./test
MAIN_DIR := ./cmd/server
COVERAGE_FILE := build/coverage.out

# Default target
.PHONY: test

test:
	$(GO) test ./... -v

test-coverage:
	$(GO) test ./... -cover -coverprofile=$(COVERAGE_FILE)

test-race:
	set CGO_ENABLED=1 && $(GO) test ./... -v -race

html-coverage:
	$(GO) tool cover -html=$(COVERAGE_FILE)

test-all:
	$(GO) test ./... -v -cover -coverprofile=$(COVERAGE_FILE) -race

run:
	@$(GO) run $(MAIN_DIR)/main.go

help:
	@echo "Available targets:"
	@echo "  test             Run all tests with verbose output"
	@echo "  test-coverage    Run tests and generate a coverage report"
	@echo "  test-race        Run tests with race detection"
	@echo "  html-coverage    Generate and view HTML coverage report"
	@echo "  test-all         Run all tests with verbose output, generate a coverage report and check for race conditions"
	@echo "  run              Run the main application"
