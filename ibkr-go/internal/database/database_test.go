package database

import (
	"context"
	"testing"
)

func TestDB_Close(t *testing.T) {
	// Test Close doesn't panic with nil pool
	db := &DB{}
	db.Close() // Should not panic
}

func TestNew_InvalidDSN(t *testing.T) {
	ctx := context.Background()

	// Test with invalid DSN
	_, err := New(ctx, "invalid-dsn", "")
	if err == nil {
		t.Error("Expected error with invalid DSN")
	}
}
