# Makefile

# Variables (tweak these as needed)
APP_NAME := my-documentdb-system
DOCKER_IMG := $(APP_NAME)
DOCKER_TAG := latest

.PHONY: all build run docker-build docker-run docker-stop clean help

## Default target: build the app
all: build

## Build the Go binary
build:
	@echo "==> Building Go binary..."
	go build -o bin/$(APP_NAME) ./cmd/main.go
	@echo "==> Build complete."

## Run the Go server locally (without Docker)
run: build
	@echo "==> Running the Go server on :8080..."
	./bin/$(APP_NAME)

## Build Docker image
docker-build:
	@echo "==> Building Docker image..."
	docker build -t $(DOCKER_IMG):$(DOCKER_TAG) .
	@echo "==> Docker build complete."

## Run with Docker Compose
docker-run:
	@echo "==> Starting Docker containers..."
	docker-compose up --build

## Stop Docker containers
docker-stop:
	@echo "==> Stopping Docker containers..."
	docker-compose down

## MongoDB Shell Access
mongo-shell:
	@echo "==> Starting MongoDB shell..."
	@docker exec -it my-mongo-db mongosh

## Clean build artifacts
clean:
	@echo "==> Cleaning up..."
	rm -rf bin

## Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the Go binary"
	@echo "  make run           - Run the Go server locally (no Docker)"
	@echo "  make docker-build  - Build the Docker image"
	@echo "  make docker-run    - Run via Docker Compose"
	@echo "  make docker-stop   - Stop Docker Compose services"
	@echo "  make mongo-shell   - Access the MongoDB shell"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make help          - Show this help message"
