package telemetry

import (
	"context"
	"testing"
)

func TestInitTracer_WithEndpoint(t *testing.T) {
	ctx := context.Background()

	// Test with invalid endpoint (will fail to connect but tests the code path)
	_, err := InitTracer(ctx, "test-service", "invalid-endpoint:4317")
	// We expect an error because the endpoint is invalid
	if err == nil {
		t.Error("Expected error with invalid endpoint")
	}
}

func TestInitTracer_EmptyServiceName(t *testing.T) {
	ctx := context.Background()

	// Test with empty service name
	tp, err := InitTracer(ctx, "", "")
	if err != nil {
		t.Fatalf("InitTracer() error = %v", err)
	}

	if tp == nil {
		t.Error("TracerProvider should not be nil")
	}

	// Cleanup
	Shutdown(ctx, tp)
}
