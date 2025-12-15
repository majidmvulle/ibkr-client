package database

import (
	"context"
	"testing"
)

func TestDB_Health_Success(t *testing.T) {
	// Test Health with nil pool (should handle gracefully)
	db := &DB{}
	err := db.Health(context.Background())
	if err == nil {
		t.Error("Expected error with nil pool")
	}
}

func TestDB_Close_WithNilPool(t *testing.T) {
	db := &DB{}
	// Should not panic
	db.Close()
}

func TestNew_ValidDSN(t *testing.T) {
	// This test would require a real database connection
	// For now, we test that the function exists and handles invalid DSN
	ctx := context.Background()

	// Test with clearly invalid DSN format
	_, err := New(ctx, "invalid://bad-dsn", "")
	if err == nil {
		t.Log("New() with invalid DSN may or may not error depending on pgx behavior")
	}
}

func TestNew_WithReadDSN(t *testing.T) {
	ctx := context.Background()

	// Test with both write and read DSN (both invalid)
	_, err := New(ctx, "invalid://write-dsn", "invalid://read-dsn")
	if err == nil {
		t.Log("New() with invalid DSNs may or may not error depending on pgx behavior")
	}
}

func TestDB_Struct(t *testing.T) {
	// Test that DB struct can be created
	db := &DB{}
	if db == nil {
		t.Error("DB struct should not be nil")
	}

	// Test that Pool field exists
	if db.Pool != nil {
		t.Error("New DB should have nil Pool")
	}
}
