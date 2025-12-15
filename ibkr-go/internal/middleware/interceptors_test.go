package middleware

import (
	"context"
	"testing"

	"connectrpc.com/connect"
)

func TestNewValidationInterceptor(t *testing.T) {
	interceptor := NewValidationInterceptor()
	if interceptor == nil {
		t.Fatal("NewValidationInterceptor() returned nil")
	}

	// Create a mock next handler
	nextCalled := false
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		nextCalled = true
		return connect.NewResponse(&struct{}{}), nil
	}

	// Wrap with interceptor
	handler := interceptor(next)

	// Create a test request
	req := connect.NewRequest(&struct{}{})

	// Call should succeed
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify next was called
	if !nextCalled {
		t.Error("Next handler should be called")
	}
}

func TestLoggingInterceptor(t *testing.T) {
	interceptor := LoggingInterceptor(nil)
	if interceptor == nil {
		t.Fatal("LoggingInterceptor() returned nil")
	}

	// Create a mock next handler
	nextCalled := false
	next := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		nextCalled = true
		return connect.NewResponse(&struct{}{}), nil
	}

	// Wrap with interceptor
	handler := interceptor(next)

	// Create a test request
	req := connect.NewRequest(&struct{}{})

	// Call should succeed
	_, err := handler(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify next was called
	if !nextCalled {
		t.Error("Next handler should be called")
	}
}
