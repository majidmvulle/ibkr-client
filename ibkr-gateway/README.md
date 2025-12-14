# IBKR Client Portal Gateway

This directory contains the Dockerfile for building a custom IBKR Client Portal Gateway image.

## Security

We build our own image from the official Interactive Brokers download instead of using third-party images. This ensures we know exactly what's running in the container.

## Building

```bash
# Build the image
make gateway-build GATEWAY_IMG=myregistry/ibkr-gateway:v1.0.0

# Push to registry
make gateway-push GATEWAY_IMG=myregistry/ibkr-gateway:v1.0.0
```

## Running Locally

The gateway is included in the docker-compose setup. See the root `docker-compose.yml` for configuration.

## Configuration

The gateway uses environment variables for configuration:
- `TWS_USERID`: Your IBKR username
- `TWS_PASSWORD`: Your IBKR password
- `TRADING_MODE`: `paper` or `live`
- `READ_ONLY_API`: `yes` or `no`

## Official Documentation

https://www.interactivebrokers.com/campus/ibkr-api-page/cpapi-v1/#gw-step-one
