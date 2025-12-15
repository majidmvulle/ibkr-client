package api

import (
	"testing"
)

func TestNewOrderServiceHandler(t *testing.T) {
	handler := NewOrderServiceHandler(nil)
	if handler == nil {
		t.Fatal("NewOrderServiceHandler() returned nil")
	}
}

func TestNewPortfolioServiceHandler(t *testing.T) {
	handler := NewPortfolioServiceHandler(nil)
	if handler == nil {
		t.Fatal("NewPortfolioServiceHandler() returned nil")
	}
}

func TestNewMarketDataServiceHandler(t *testing.T) {
	handler := NewMarketDataServiceHandler(nil)
	if handler == nil {
		t.Fatal("NewMarketDataServiceHandler() returned nil")
	}
}
