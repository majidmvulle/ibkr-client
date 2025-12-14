# IBKR Client Helm Chart

This Helm chart deploys the IBKR Client Portal Gateway middleware service to Kubernetes (k3s).

## Architecture

The deployment consists of three containers:

1. **Init Container (migrations)**: Runs database migrations using Goose before the main application starts
2. **Main Container (ibkr-client)**: The IBKR middleware service
3. **Sidecar Container (ibkr-gateway)**: IBKR Client Portal Gateway running on localhost:5000

## Prerequisites

- Kubernetes 1.19+ (k3s)
- Helm 3.0+
- PostgreSQL database (external or in-cluster)
- IBKR account credentials
- Kubernetes Secret with sensitive data

## Installation

### 1. Create Kubernetes Secret

First, create a secret with all sensitive data:

```bash
# Generate a secure 32-byte encryption key
ENCRYPTION_KEY=$(openssl rand -base64 32 | head -c 32)

# Create secret
kubectl create secret generic ibkr-client-secrets \
  --namespace ibkr \
  --from-literal=db-write-dsn='postgres://user:pass@host:5432/dbname?sslmode=disable' \
  --from-literal=db-read-dsn='postgres://user:pass@host:5432/dbname?sslmode=disable' \
  --from-literal=encryption-key="$ENCRYPTION_KEY" \
  --from-literal=ibkr-account-id='DU123456'
```

### 2. Install Helm Chart

```bash
# Install with default values
helm install ibkr-client ./infra/helm/charts/ibkr-client \
  --namespace ibkr \
  --create-namespace \
  --set image.tag=v0.1.0 \
  --set migrations.image.tag=v0.1.0 \
  --set ibkrGateway.image.tag=stable

# Or with custom values file
helm install ibkr-client ./infra/helm/charts/ibkr-client \
  --namespace ibkr \
  --create-namespace \
  --values my-values.yaml
```

## Configuration

### Key Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `replicaCount` | int | `1` | Number of replicas (set as needed) |
| `image.repository` | string | `ghcr.io/majidmvulle/ibkr-client` | Image repository |
| `image.tag` | string | `""` | Image tag (defaults to chart appVersion) |
| `migrations.image.repository` | string | `ghcr.io/majidmvulle/ibkr-client-migrations` | Migrations image |
| `migrations.image.tag` | string | `""` | Migrations image tag |
| `ibkrGateway.image.repository` | string | `ghcr.io/majidmvulle/ibkr-gateway` | IBKR Gateway image |
| `ibkrGateway.image.tag` | string | `""` | IBKR Gateway image tag |
| `ibkrGateway.port` | int | `5000` | IBKR Gateway port |
| `config.appEnv` | string | `production` | Application environment |
| `config.appDebug` | bool | `false` | Enable debug mode |
| `config.httpPort` | int | `8080` | HTTP port |
| `config.grpcPort` | int | `50051` | gRPC port |
| `secrets.secretName` | string | `ibkr-client-secrets` | Name of Kubernetes Secret |

### Example Custom Values

```yaml
# my-values.yaml
replicaCount: 2

image:
  tag: "v0.2.0"

migrations:
  image:
    tag: "v0.2.0"

ibkrGateway:
  image:
    tag: "stable"

config:
  appEnv: "staging"
  appDebug: true
  logLevel: -4  # Debug

resources:
  limits:
    cpu: 1000m
    memory: 1Gi
  requests:
    cpu: 200m
    memory: 256Mi
```

## Upgrading

```bash
# Upgrade with new image tag
helm upgrade ibkr-client ./infra/helm/charts/ibkr-client \
  --namespace ibkr \
  --set image.tag=v0.2.0 \
  --set migrations.image.tag=v0.2.0
```

## Uninstalling

```bash
helm uninstall ibkr-client --namespace ibkr
```

## Accessing the Service

Since the service uses ClusterIP, you can access it via:

### Port Forward
```bash
# HTTP
kubectl port-forward -n ibkr svc/ibkr-client 8080:8080

# gRPC
kubectl port-forward -n ibkr svc/ibkr-client 50051:50051
```

### Via Ingress Controller
Configure your existing ingress controller to route to the service.

### Via Service Mesh
If using a service mesh, configure routing rules.

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
- Secrets stored in Kubernetes secrets (create separately)

## License

MIT
