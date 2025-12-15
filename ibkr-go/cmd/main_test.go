package main

import (
	"context"
	"testing"
)

func TestContextUsage(t *testing.T) {
	ctx := context.Background()
	if ctx == nil {
		t.Error("Context should not be nil")
	}
}
