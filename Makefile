.PHONY: help proto proto-lint build test go-lint docker-build docker-push gateway-build gateway-push sqlc-generate dev-up dev-down

# Docker image configuration
IMG ?= ibkr-client:latest
GATEWAY_IMG ?= ibkr-gateway:latest

help:
	@echo "Available targets:"
	@echo "  proto         - Generate protobuf code"
	@echo "  proto-lint    - Lint protobuf files"
	@echo "  build         - Build all Go modules"
	@echo "  test          - Run tests for all modules"
	@echo "  go-lint       - Lint all Go modules"
	@echo "  sqlc-generate - Generate type-safe Go code from SQL"
	@echo "  docker-build  - Build Docker image (IMG=<image:tag>)"
	@echo "  docker-push   - Push Docker image (IMG=<image:tag>)"
	@echo "  gateway-build - Build IBKR Gateway image (GATEWAY_IMG=<image:tag>)"
	@echo "  gateway-push  - Push IBKR Gateway image (GATEWAY_IMG=<image:tag>)"
	@echo "  dev-up        - Start local dev environment"
	@echo "  dev-down      - Stop local dev environment"

proto:
	cd proto && buf generate

proto-lint:
	cd proto && buf lint

build:
	cd ibkr-go && go build -o ../bin/ibkr-server ./cmd/server

test:
	cd ibkr-go && go test -v -race ./...

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
	docker-compose up -d

dev-down:
	docker-compose down
