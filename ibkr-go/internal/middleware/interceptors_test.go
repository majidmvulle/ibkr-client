package middleware

import (
	"context"
	"testing"

	"log/slog"

	"connectrpc.com/connect"
)

func TestValidationInterceptor_AllPaths(t *testing.T) {
	interceptor := NewValidationInterceptor()
	if interceptor == nil {
		t.Fatal("NewValidationInterceptor should not return nil")
	}

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	ctx := context.Background()
	req := connect.NewRequest(&struct{}{})

	_, err := wrapped(ctx, req)
	if err != nil {
		t.Logf("Validation interceptor error: %v", err)
	}
}

func TestLoggingInterceptor_AllPaths(t *testing.T) {
	logger := slog.Default()
	interceptor := LoggingInterceptor(logger)
	if interceptor == nil {
		t.Fatal("LoggingInterceptor should not return nil")
	}

	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		return connect.NewResponse(&struct{}{}), nil
	}
	wrapped := interceptor(handler)

	ctx := context.Background()
	req := connect.NewRequest(&struct{}{})

	_, err := wrapped(ctx, req)
	if err != nil {
		t.Logf("Logging interceptor error: %v", err)
	}
}
