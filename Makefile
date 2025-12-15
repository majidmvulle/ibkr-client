.PHONY: help proto proto-lint build test test-unit test-integration test-coverage go-lint sqlc-generate docker-build docker-push helm-package helm-push dev-up dev-down dev-logs dev-restart migrate-up migrate-down clean

help:
	@echo "Available targets:"
	@echo "  proto                  - Generate protobuf code"
	@echo "  proto-lint             - Lint protobuf files"
	@echo "  build                  - Build all Go modules"
	@echo "  test                   - Run all tests"
	@echo "  test-unit              - Run unit tests only"
	@echo "  test-integration       - Run integration tests only"
	@echo "  test-coverage          - Run tests with coverage report"
	@echo "  go-lint                - Lint all Go modules"
	@echo "  sqlc-generate          - Generate type-safe Go code from SQL"
	@echo "  docker-build           - Build all images (server, migrations, gateway) for multiple platforms"
	@echo "  docker-push            - Push all images to registry"
	@echo "  helm-package           - Package Helm chart"
	@echo "  helm-push              - Package and push Helm chart to OCI registry"
	@echo "  dev-up                 - Start local dev environment"
	@echo "  dev-down               - Stop local dev environment"
	@echo "  dev-logs               - Show logs from dev environment"
	@echo "  dev-restart            - Restart dev environment"
	@echo "  migrate-up             - Run database migrations"
	@echo "  migrate-down           - Rollback database migrations"
	@echo "  clean                  - Clean build artifacts"

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

# Docker configuration
REGISTRY ?= ghcr.io/majidmvulle/ibkr-client
VERSION ?= v0.1.0
PLATFORMS ?= linux/amd64,linux/arm64

# Image names
SERVER_IMAGE := $(REGISTRY)/ibkr-client
MIGRATIONS_IMAGE := $(REGISTRY)/ibkr-client-migrations
GATEWAY_IMAGE := $(REGISTRY)/ibkr-gateway

go-lint:
	cd ibkr-go && golangci-lint run ./... --fix

sqlc-generate:
	cd ibkr-go && sqlc generate

docker-build:
	@echo "Building all images for $(PLATFORMS)..."
	@echo "Building server image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(SERVER_IMAGE):$(VERSION) \
		-t $(SERVER_IMAGE):latest \
		-f ibkr-go/Dockerfile \
		ibkr-go/
	@echo "Building migrations image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(MIGRATIONS_IMAGE):$(VERSION) \
		-t $(MIGRATIONS_IMAGE):latest \
		-f ibkr-go/Dockerfile.migrations \
		ibkr-go/
	@echo "Building IBKR Gateway image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		-t $(GATEWAY_IMAGE):$(VERSION) \
		-t $(GATEWAY_IMAGE):latest \
		-f ibkr-gateway/Dockerfile \
		ibkr-gateway/
	@echo "All images built successfully!"

docker-push:
	@echo "Pushing all images to $(REGISTRY)..."
	@echo "Pushing server image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		--push \
		-t $(SERVER_IMAGE):$(VERSION) \
		-t $(SERVER_IMAGE):latest \
		-f ibkr-go/Dockerfile \
		ibkr-go/
	@echo "Pushing migrations image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		--push \
		-t $(MIGRATIONS_IMAGE):$(VERSION) \
		-t $(MIGRATIONS_IMAGE):latest \
		-f ibkr-go/Dockerfile.migrations \
		ibkr-go/
	@echo "Pushing IBKR Gateway image..."
	docker buildx build \
		--platform $(PLATFORMS) \
		--push \
		-t $(GATEWAY_IMAGE):$(VERSION) \
		-t $(GATEWAY_IMAGE):latest \
		-f ibkr-gateway/Dockerfile \
		ibkr-gateway/
	@echo "All images pushed successfully!"

helm-package:
	@echo "Packaging Helm chart..."
	helm package infra/helm/charts/ibkr-client --version $(VERSION)

helm-push:
	@echo "Pushing Helm chart to OCI registry..."
	helm package infra/helm/charts/ibkr-client --version $(VERSION)
	helm push ibkr-client-$(VERSION).tgz oci://$(REGISTRY)
	@echo "Helm chart pushed successfully!"
	@rm -f ibkr-client-$(VERSION).tgz

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
