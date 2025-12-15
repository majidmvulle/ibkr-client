# Integration Tests

This directory contains integration tests for the IBKR Client services.

## Overview

Integration tests verify that all components work together correctly:
- **Order Service**: Place, modify, cancel, and list orders
- **Portfolio Service**: Get portfolio, positions, and account summary
- **Market Data Service**: Get quotes, historical data, and streaming quotes
- **Session Management**: Create, validate, delete, and cleanup sessions
- **Database**: Connection health, queries, and concurrent access

## Running Tests

### All Integration Tests

```bash
make test-integration
```

### Specific Test File

```bash
cd ibkr-go
go test -v ./test/integration/order_service_test.go
```

### With Coverage

```bash
make test-coverage
```

## Test Structure

```
test/integration/
├── setup_test.go              # Test setup and shared utilities
├── order_service_test.go      # Order Service integration tests
├── portfolio_service_test.go  # Portfolio Service integration tests
├── marketdata_service_test.go # Market Data Service integration tests
├── session_test.go            # Session Management integration tests
└── database_test.go           # Database integration tests
```

## Prerequisites

Integration tests require:
1. **PostgreSQL database** running (via docker-compose)
2. **Mock IBKR Gateway** running (via docker-compose)
3. **Test environment** configured (.env.test)

### Quick Setup

```bash
# Start test environment
make dev-up

# Run integration tests
make test-integration

# Stop test environment
make dev-down
```

## Test Coverage

### Order Service Tests
- ✅ Place order (market, limit)
- ✅ List orders
- ✅ Cancel order
- ✅ Invalid symbol handling

### Portfolio Service Tests
- ✅ Get portfolio with positions
- ✅ Get positions list
- ✅ Get account summary
- ✅ Money conversion validation

### Market Data Service Tests
- ✅ Get quote for symbol
- ✅ Get historical data
- ✅ Stream quotes (server streaming)
- ✅ Invalid symbol handling

### Session Management Tests
- ✅ Create and validate session
- ✅ Invalid token rejection
- ✅ Delete and validate
- ✅ Expired session cleanup
- ✅ Multiple concurrent sessions
- ✅ Token uniqueness

### Database Tests
- ✅ Connection health (write/read pools)
- ✅ Session CRUD operations
- ✅ Health check endpoint
- ✅ Transaction isolation
- ✅ Concurrent access

## Test Helpers

### CreateTestSession
Creates a test session and returns the token:
```go
token := CreateTestSession(t, accountID)
defer DeleteTestSession(t, token)
```

### ValidateTestSession
Validates a session token:
```go
accountID := ValidateTestSession(t, token)
```

### DeleteTestSession
Deletes a test session:
```go
DeleteTestSession(t, token)
```

## CI/CD Integration

Integration tests are designed to run in CI/CD pipelines:

```yaml
# GitHub Actions example
- name: Start test environment
  run: make dev-up

- name: Run integration tests
  run: make test-integration

- name: Upload coverage
  run: make test-coverage
```

## Troubleshooting

### Tests Failing to Connect

Ensure services are running:
```bash
docker-compose ps
```

### Database Connection Issues

Check PostgreSQL logs:
```bash
docker-compose logs postgres
```

### Mock Gateway Not Responding

Check mock gateway logs:
```bash
docker-compose logs ibkr-gateway-mock
```

### Cleanup Test Data

Test data is automatically cleaned up in `TestMain`, but you can manually clean:
```bash
# Connect to database
docker-compose exec postgres psql -U ibkr -d ibkr_test

# Delete test sessions
DELETE FROM sessions WHERE created_at < NOW();
```

## Best Practices

1. **Use `testing.Short()`**: Skip integration tests in short mode
2. **Clean up resources**: Always defer cleanup (DeleteTestSession)
3. **Use test context**: Access shared resources via `testCtx`
4. **Verify responses**: Check all important fields
5. **Log useful info**: Use `t.Logf()` for debugging

## License

MIT
