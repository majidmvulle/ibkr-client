.PHONY: help proto proto-lint build test test-unit test-integration test-coverage go-lint docker-build docker-push gateway-build gateway-push sqlc-generate dev-up dev-down dev-logs dev-restart migrate-up migrate-down clean

# Docker image configuration
IMG ?= ibkr-client:latest
GATEWAY_IMG ?= ibkr-gateway:latest

help:
	@echo "Available targets:"
	@echo "  proto            - Generate protobuf code"
	@echo "  proto-lint       - Lint protobuf files"
	@echo "  build            - Build all Go modules"
	@echo "  test             - Run all tests"
	@echo "  test-unit        - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  go-lint          - Lint all Go modules"
	@echo "  sqlc-generate    - Generate type-safe Go code from SQL"
	@echo "  docker-build     - Build Docker image (IMG=<image:tag>)"
	@echo "  docker-push      - Push Docker image (IMG=<image:tag>)"
	@echo "  gateway-build    - Build IBKR Gateway image (GATEWAY_IMG=<image:tag>)"
	@echo "  gateway-push     - Push IBKR Gateway image (GATEWAY_IMG=<image:tag>)"
	@echo "  dev-up           - Start local dev environment"
	@echo "  dev-down         - Stop local dev environment"
	@echo "  dev-logs         - Show logs from dev environment"
	@echo "  dev-restart      - Restart dev environment"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback database migrations"
	@echo "  clean            - Clean build artifacts"

proto:
	cd proto && buf generate

proto-lint:
	cd proto && buf lint

build:
	cd ibkr-go && go build -o ../bin/ibkr-server ./cmd/server

test:
	@echo "Running all tests..."
	cd ibkr-go && go test -v -race ./...

test-unit:
	@echo "Running unit tests..."
	cd ibkr-go && go test -v -short -race ./...

test-integration:
	@echo "Running integration tests..."
	@if [ ! -f .env.test ]; then \
		echo "Creating .env.test from .env.test.example..."; \
		cp .env.test.example .env.test; \
	fi
	@set -a && . ./.env.test && set +a && cd ibkr-go && go test -v -run Integration ./test/integration/...

test-coverage:
	@echo "Running tests with coverage..."
	cd ibkr-go && go test -v -short -coverprofile=coverage.out -covermode=atomic ./...
	cd ibkr-go && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: ibkr-go/coverage.html"

test-coverage-report:
	@echo "Coverage summary:"
	cd ibkr-go && go tool cover -func=coverage.out

go-lint:
	cd ibkr-go && golangci-lint run ./... --fix

sqlc-generate:
	cd ibkr-go && sqlc generate

docker-build:
	docker build -t $(IMG) -f ibkr-go/Dockerfile .

docker-push:
	docker push $(IMG)

gateway-build:
	docker build -t $(GATEWAY_IMG) -f ibkr-gateway/Dockerfile ibkr-gateway

gateway-push:
	docker push $(GATEWAY_IMG)

dev-up:
	@echo "Starting development environment..."
	docker-compose up -d
	@echo "Waiting for services..."
	@sleep 5
	@echo "Development environment ready!"

dev-down:
	docker-compose down

dev-logs:
	docker-compose logs -f

dev-restart:
	docker-compose restart

migrate-up:
	@echo "Running migrations..."
	docker-compose run --rm migrations

migrate-down:
	@echo "Rolling back migrations..."
	cd migrations && goose postgres "$(DB_WRITE_DSN)" down

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf ibkr-go/coverage.out ibkr-go/coverage.html
