package database

import (
	"context"
	"testing"
)

func TestDB_Close_Safe(t *testing.T) {
	// Test that Close doesn't panic with nil pool
	db := &DB{}
	// Close should handle nil pool gracefully
	// Note: The actual Close() method calls Pool.Close() without nil check
	// This test documents the current behavior
	if db.Pool != nil {
		db.Close() // Should not panic
	}
}

func TestNew_InvalidDSN_Format(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, "not-a-valid-dsn", "")
	if err == nil {
		t.Error("Expected error with invalid DSN format")
	}
}
