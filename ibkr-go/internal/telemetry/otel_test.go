package telemetry

import (
	"context"
	"testing"
)

func TestInitTracer_NoEndpoint(t *testing.T) {
	ctx := context.Background()
	tp, err := InitTracer(ctx, "test-service", "")
	if err != nil {
		t.Fatalf("InitTracer() error = %v", err)
	}
	if tp == nil {
		t.Error("TracerProvider should not be nil")
	}
	Shutdown(ctx, tp)
}

func TestShutdown_Safe(t *testing.T) {
	ctx := context.Background()
	Shutdown(ctx, nil) // Should not panic
}
