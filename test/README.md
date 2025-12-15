# Local Testing Environment

This directory contains the local testing infrastructure for the IBKR Client project.

## Overview

The local testing environment includes:
- **PostgreSQL database** for testing
- **Mock IBKR Gateway** for API testing
- **Docker Compose** setup for easy management
- **Makefile targets** for common tasks

## Quick Start

### 1. Start the Development Environment

```bash
make dev-up
```

This will start:
- PostgreSQL on `localhost:5432`
- Mock IBKR Gateway on `localhost:5000`
- IBKR Client on `localhost:8080` (HTTP) and `localhost:50051` (gRPC)

### 2. Run Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run tests with coverage
make test-coverage
```

### 3. View Logs

```bash
make dev-logs
```

### 4. Stop the Environment

```bash
make dev-down
```

## Mock IBKR Gateway

The mock gateway (`test/mocks/mock_gateway.py`) provides mock responses for all IBKR API endpoints:

### Supported Endpoints

**Authentication:**
- `POST /v1/api/tickle` - Health check
- `POST /v1/api/iserver/auth/status` - Auth status
- `POST /v1/api/iserver/reauthenticate` - Reauthenticate
- `GET /v1/api/iserver/accounts` - Get accounts

**Orders:**
- `POST /v1/api/iserver/account/{account_id}/orders` - Place order
- `POST /v1/api/iserver/account/{account_id}/order/{order_id}` - Modify order
- `DELETE /v1/api/iserver/account/{account_id}/order/{order_id}` - Cancel order
- `GET /v1/api/iserver/account/orders` - Get live orders

**Portfolio:**
- `GET /v1/api/portfolio/{account_id}/positions/0` - Get positions
- `GET /v1/api/portfolio/{account_id}/summary` - Get account summary

**Market Data:**
- `GET /v1/api/iserver/marketdata/snapshot` - Get market data snapshot
- `GET /v1/api/iserver/marketdata/history` - Get historical data
- `POST /v1/api/iserver/secdef/search` - Search contracts

### Mock Data

- **Account ID**: `DU123456`
- **Test Symbol**: `AAPL` (conid: 265598)
- **Mock Prices**: Last: `$150.00`, Bid: `$149.50`, Ask: `$150.50`

## Environment Configuration

Copy `.env.test.example` to `.env.test` and customize:

```bash
cp .env.test.example .env.test
```

Key variables:
- `DB_WRITE_DSN` - PostgreSQL connection string
- `ENCRYPTION_KEY` - 32-byte encryption key
- `IBKR_GATEWAY_URL` - Mock gateway URL (http://localhost:5000)
- `IBKR_ACCOUNT_ID` - Test account ID (DU123456)

## Database Migrations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make dev-up` | Start development environment |
| `make dev-down` | Stop development environment |
| `make dev-logs` | Show logs |
| `make dev-restart` | Restart services |
| `make test` | Run all tests |
| `make test-unit` | Run unit tests |
| `make test-integration` | Run integration tests |
| `make test-coverage` | Generate coverage report |
| `make migrate-up` | Run migrations |
| `make migrate-down` | Rollback migrations |
| `make clean` | Clean build artifacts |

## Troubleshooting

### Port Already in Use

If ports 5432, 5000, 8080, or 50051 are already in use:

```bash
# Check what's using the port
lsof -i :5432

# Stop the service or change ports in docker-compose.yml
```

### Database Connection Issues

```bash
# Check PostgreSQL is running
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres
```

### Mock Gateway Not Responding

```bash
# Check mock gateway logs
docker-compose logs ibkr-gateway-mock

# Restart the mock gateway
docker-compose restart ibkr-gateway-mock
```

## Writing Tests

### Unit Tests

Place unit tests next to the code they test:

```
ibkr-go/
├── internal/
│   ├── session/
│   │   ├── service.go
│   │   └── service_test.go
```

### Integration Tests

Place integration tests in `ibkr-go/test/integration/`:

```
ibkr-go/
├── test/
│   ├── integration/
│   │   ├── setup_test.go
│   │   ├── order_service_test.go
│   │   └── portfolio_service_test.go
```

## CI/CD Integration

The test environment can be used in CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Start test environment
  run: make dev-up

- name: Run tests
  run: make test-coverage

- name: Stop test environment
  run: make dev-down
```

## License

MIT
