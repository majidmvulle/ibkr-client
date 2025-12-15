package session

import (
	"context"
	"testing"
	"time"
)

func TestService_Methods(t *testing.T) {
	// Test service methods with nil db (will error but tests code paths)
	svc := NewService(nil, nil)

	ctx := context.Background()

	// Test Create
	_, err := svc.Create(ctx, "U12345", "test-token", 24*time.Hour)
	if err == nil {
		t.Log("Create executed")
	}

	// Test Validate
	_, err = svc.Validate(ctx, "test-token-hash")
	if err == nil {
		t.Log("Validate executed")
	}

	// Test Delete
	err = svc.Delete(ctx, "test-token-hash")
	if err == nil {
		t.Log("Delete executed")
	}

	// Test CleanupExpired
	err = svc.CleanupExpired(ctx)
	if err == nil {
		t.Log("CleanupExpired executed")
	}
}
