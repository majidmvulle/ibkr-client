package middleware

import (
	"context"
	"testing"

	"connectrpc.com/connect"
)

func TestInterceptorChaining(t *testing.T) {
	// Test that interceptors can be chained
	interceptor1 := NewValidationInterceptor()
	interceptor2 := LoggingInterceptor(nil)

	if interceptor1 == nil || interceptor2 == nil {
		t.Fatal("Interceptors should not be nil")
	}

	// Create a test handler
	called := false
	handler := func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		called = true
		return connect.NewResponse(&struct{}{}), nil
	}

	// Chain interceptors
	chained := interceptor1(interceptor2(handler))
	req := connect.NewRequest(&struct{}{})

	_, err := chained(context.Background(), req)
	if err != nil {
		t.Errorf("Chained interceptors error: %v", err)
	}

	if !called {
		t.Error("Handler should have been called")
	}
}

func TestContextPropagation(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ClientIdentityContextKey{}, "test-client")

	clientID, ok := GetAccountIDFromContext(ctx)
	if !ok {
		t.Error("Should find client ID in context")
	}

	if clientID != "test-client" {
		t.Errorf("Client ID = %v, want test-client", clientID)
	}
}
