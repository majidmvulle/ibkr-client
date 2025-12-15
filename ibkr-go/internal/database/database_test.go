package database

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	// Test with empty DSN should return error
	_, err := New(context.Background(), "", "")
	if err == nil {
		t.Error("New() with empty DSN should return error")
	}
}

func TestDB_Health(t *testing.T) {
	// Create a DB with nil pool (simulating uninitialized state)
	db := &DB{}

	err := db.Health(context.Background())
	if err == nil {
		t.Error("Health() with nil pool should return error")
	}
}
