package database

import (
	"context"
	"testing"
)

func TestDB_Close_Safe(t *testing.T) {
	// Test that Close doesn't panic with nil pool
	db := &DB{}
	db.Close() // Should not panic
}

func TestNew_EmptyDSN(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, "", "")
	if err == nil {
		t.Error("Expected error with empty DSN")
	}
}

func TestNew_InvalidDSN_Format(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, "not-a-valid-dsn", "")
	if err == nil {
		t.Error("Expected error with invalid DSN format")
	}
}
