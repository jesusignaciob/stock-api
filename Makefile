# Variables
APP_NAME := stock-api
MAIN_FILE := ./cmd/main.go
GO_FILES := $(shell find . -type f -name '*.go')
BUILD_DIR := ./bin
BINARY := $(BUILD_DIR)/$(APP_NAME)

# Default target
.PHONY: all
all: build

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE) --mode=api

.PHONY: run-data
run-data:
	go run $(MAIN_FILE) --mode=data

# Build the application
.PHONY: build
build: $(BINARY)

$(BINARY): $(GO_FILES)
	@mkdir -p $(BUILD_DIR)
	go build -o $(BINARY) $(MAIN_FILE)
	@echo "Build complete: $(BINARY)"

# Run tests
.PHONY: test
test:
	go test ./... -v

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Analyze code
# This target runs go vet, staticcheck, and golangci-lint
# to check for potential issues in the codebase.
.PHONY: analyze
analyze:
	@go vet ./...
	@staticcheck ./...
	@golangci-lint run

# Format code
# This target formats the code using gofumpt, goimports, and gci.
.PHONY: format
format:
	@echo "Formatting code..."
	@gofumpt -w .
	@goimports -w -local $(APP_NAME) .
	@gci write -s standard -s default -s "prefix($(APP_NAME))" .

# Fix lint issues
# This target runs golangci-lint with the --fix option to automatically fix lint issues.
.PHONY: fix
fix:
	@echo "Fixing lint issues..."
	@golangci-lint run --fix ./...
	@go vet ./...
	@go mod tidy

# Database migration targets
# These targets use the migrate tool to manage database migrations.
.PHONY: migrate-up
migrate-up:
	go run $(MAIN_FILE) --migrate=up

.PHONY: migrate-down
migrate-down:
	go run $(MAIN_FILE) --migrate=down

.PHONY: help
help:
	@echo "Makefile for $(APP_NAME)"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all            Build the application"
	@echo "  run            Run the application"
	@echo "  run-data       Run the data mode of the application"
	@echo "  build          Build the application"
	@echo "  test           Run tests"
	@echo "  clean          Clean build artifacts"
	@echo "  fmt            Format code"
	@echo "  lint           Lint code"
	@echo "  deps           Install dependencies"
	@echo "  analyze        Analyze code"
	@echo "  format         Format code"
	@echo "  fix            Fix lint issues"
	@echo "  migrate-up     Run database migrations up"
	@echo "  migrate-down   Run database migrations down"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Environment Variables:"
	@echo "  DB_HOST        Database host"
	@echo "  DB_PORT        Database port"
	@echo "  DB_USER        Database user"
	@echo "  DB_PASSWORD    Database password"
	@echo "  DB_NAME        Database name"
	@echo "  DB_SSLMODE     Database SSL mode"
	@echo ""
	@echo "Note: Make sure to set the environment variables before running the make commands."
	@echo ""
	@echo "For more information, visit the project repository."
	@echo ""
	@echo "Happy coding!"
	@echo ""