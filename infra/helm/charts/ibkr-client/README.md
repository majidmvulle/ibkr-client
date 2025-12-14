# IBKR Client Helm Chart

This Helm chart deploys the IBKR Client Portal Gateway middleware service to Kubernetes.

## Architecture

The deployment consists of three containers:

1. **Init Container (migrations)**: Runs database migrations using Goose before the main application starts
2. **Main Container (ibkr-client)**: The IBKR middleware service
3. **Sidecar Container (ibkr-gateway)**: IBKR Client Portal Gateway running on localhost:5000

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- PostgreSQL database (external or in-cluster)
- IBKR account credentials

## Installation

### Development

```bash
# Install with development values
helm install ibkr-client ./infra/helm/ibkr-client \
  --namespace ibkr \
  --create-namespace \
  --values ./infra/helm/ibkr-client/values-dev.yaml
```

### Production

```bash
# Install with production values
helm install ibkr-client ./infra/helm/ibkr-client \
  --namespace ibkr \
  --create-namespace \
  --values ./infra/helm/ibkr-client/values-prod.yaml \
  --set image.tag=v0.1.0
```

## Configuration

### Required Secrets

The following secrets must be configured:

- `dbWriteDSN`: PostgreSQL write connection string
- `dbReadDSN`: PostgreSQL read connection string
- `encryptionKey`: 32-byte encryption key for session tokens
- `ibkrAccountID`: IBKR account ID

### Example Secret Creation

```bash
# Generate a secure 32-byte encryption key
ENCRYPTION_KEY=$(openssl rand -base64 32 | head -c 32)

# Create secret
kubectl create secret generic ibkr-client \
  --namespace ibkr \
  --from-literal=db-write-dsn='postgres://user:pass@host:5432/dbname' \
  --from-literal=db-read-dsn='postgres://user:pass@host:5432/dbname' \
  --from-literal=encryption-key="$ENCRYPTION_KEY" \
  --from-literal=ibkr-account-id='DU123456'
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `replicaCount` | int | `1` | Number of replicas |
| `image.repository` | string | `ghcr.io/majidmvulle/ibkr-client` | Image repository |
| `image.tag` | string | `""` | Image tag (defaults to chart appVersion) |
| `migrations.image.repository` | string | `ghcr.io/majidmvulle/ibkr-client-migrations` | Migrations image |
| `ibkrGateway.image.repository` | string | `ghcr.io/majidmvulle/ibkr-gateway` | IBKR Gateway image |
| `ibkrGateway.port` | int | `5000` | IBKR Gateway port |
| `config.appEnv` | string | `production` | Application environment |
| `config.appDebug` | bool | `false` | Enable debug mode |
| `config.httpPort` | int | `8080` | HTTP port |
| `config.grpcPort` | int | `50051` | gRPC port |
| `autoscaling.enabled` | bool | `false` | Enable HPA |
| `autoscaling.minReplicas` | int | `1` | Minimum replicas |
| `autoscaling.maxReplicas` | int | `10` | Maximum replicas |
| `ingress.enabled` | bool | `false` | Enable ingress |

## Upgrading

```bash
# Upgrade with new image tag
helm upgrade ibkr-client ./infra/helm/ibkr-client \
  --namespace ibkr \
  --values ./infra/helm/ibkr-client/values-prod.yaml \
  --set image.tag=v0.2.0
```

## Uninstalling

```bash
helm uninstall ibkr-client --namespace ibkr
```

## Health Checks

The service exposes two health check endpoints:

- `/healthz`: Liveness probe (basic health check)
- `/readyz`: Readiness probe (includes database connectivity check)

## Monitoring

The service can be integrated with OpenTelemetry by setting:

```yaml
config:
  otelCollectorEndpoint: "http://otel-collector:4317"
```

## Troubleshooting

### Check pod status
```bash
kubectl get pods -n ibkr
```

### View logs
```bash
# Main container
kubectl logs -n ibkr -l app.kubernetes.io/name=ibkr-client -c ibkr-client

# IBKR Gateway sidecar
kubectl logs -n ibkr -l app.kubernetes.io/name=ibkr-client -c ibkr-gateway

# Migration init container
kubectl logs -n ibkr -l app.kubernetes.io/name=ibkr-client -c migrations
```

### Check migrations
```bash
kubectl describe pod -n ibkr -l app.kubernetes.io/name=ibkr-client
```

## Security

- All containers run as non-root user (UID 1000)
- Read-only root filesystem
- No privilege escalation
- Capabilities dropped
- Secrets stored in Kubernetes secrets (or external secret manager in production)

## License

MIT
