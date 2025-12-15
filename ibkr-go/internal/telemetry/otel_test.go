package telemetry

import (
	"context"
	"testing"
)

func TestInitTracer_AllPaths(t *testing.T) {
	ctx := context.Background()

	// Test with empty endpoint (no-op tracer)
	tp, err := InitTracer(ctx, "test-service", "")
	if err != nil {
		t.Fatalf("InitTracer() error = %v", err)
	}
	if tp == nil {
		t.Error("TracerProvider should not be nil")
	}
	Shutdown(ctx, tp)

	// Test with invalid endpoint
	tp2, err := InitTracer(ctx, "test", "invalid:4317")
	if err == nil {
		t.Log("InitTracer with invalid endpoint")
		if tp2 != nil {
			Shutdown(ctx, tp2)
		}
	}
}

func TestShutdown_Nil(t *testing.T) {
	ctx := context.Background()
	Shutdown(ctx, nil) // Should not panic
}
