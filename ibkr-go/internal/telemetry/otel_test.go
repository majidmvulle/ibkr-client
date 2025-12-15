package telemetry

import (
	"context"
	"testing"
	"time"
)

func TestInitTracer_EmptyEndpoint(t *testing.T) {
	ctx := context.Background()
	tp, err := InitTracer(ctx, "test-service", "")
	if err != nil {
		t.Fatalf("InitTracer() with empty endpoint should not error, got: %v", err)
	}
	if tp == nil {
		t.Error("Expected non-nil tracer provider")
	}

	// Clean up
	if tp != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = Shutdown(shutdownCtx, tp)
	}
}

func TestInitTracer_InvalidEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Use an invalid endpoint that will fail to connect
	tp, err := InitTracer(ctx, "test-service", "invalid-endpoint:9999")
	// The function may or may not error depending on connection timeout
	// but it should return a tracer provider
	if tp != nil {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer shutdownCancel()
		_ = Shutdown(shutdownCtx, tp)
	}

	// We mainly want to ensure the function doesn't panic
	t.Logf("InitTracer with invalid endpoint: err=%v, tp=%v", err, tp != nil)
}

func TestShutdown_NilProvider(t *testing.T) {
	ctx := context.Background()
	err := Shutdown(ctx, nil)
	if err != nil {
		t.Errorf("Shutdown(nil) should not error, got: %v", err)
	}
}

func TestShutdown_ValidProvider(t *testing.T) {
	ctx := context.Background()
	tp, err := InitTracer(ctx, "test-service", "")
	if err != nil {
		t.Fatalf("InitTracer() error = %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = Shutdown(shutdownCtx, tp)
	if err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestInitTracer_WithServiceName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		endpoint    string
	}{
		{"empty service name", "", ""},
		{"valid service name", "my-service", ""},
		{"service with spaces", "my service", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			tp, err := InitTracer(ctx, tt.serviceName, tt.endpoint)
			if err != nil {
				t.Errorf("InitTracer() error = %v", err)
			}
			if tp == nil {
				t.Error("Expected non-nil tracer provider")
			}

			if tp != nil {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_ = Shutdown(shutdownCtx, tp)
			}
		})
	}
}
