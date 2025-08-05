# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=gin-boilerplate
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run the application
run:
	$(GOCMD) run main.go

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Generate swagger docs
swagger:
	swag init

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Build for production
build-prod:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) -a -installsuffix cgo -o $(BINARY_UNIX) .

.PHONY: build clean run test test-coverage deps swagger dev fmt lint build-prod
