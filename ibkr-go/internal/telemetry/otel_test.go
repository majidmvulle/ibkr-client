package telemetry

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestInitTracer_NoEndpoint(t *testing.T) {
	ctx := context.Background()
	tp, err := InitTracer(ctx, "test-service", "")
	if err != nil {
		t.Fatalf("InitTracer() with empty endpoint should not error: %v", err)
	}
	if tp == nil {
		t.Fatal("InitTracer() returned nil tracer provider")
	}

	// Cleanup
	if err := Shutdown(ctx, tp); err != nil {
		t.Errorf("Shutdown() error: %v", err)
	}
}

func TestShutdown_NilProvider(t *testing.T) {
	ctx := context.Background()
	err := Shutdown(ctx, nil)
	if err != nil {
		t.Errorf("Shutdown() with nil provider should not error: %v", err)
	}
}

func TestShutdown_ValidProvider(t *testing.T) {
	ctx := context.Background()
	tp := sdktrace.NewTracerProvider()

	err := Shutdown(ctx, tp)
	if err != nil {
		t.Errorf("Shutdown() error: %v", err)
	}
}
