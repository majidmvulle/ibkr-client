package ibkr

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		gatewayURL string
		accountID  string
	}{
		{
			name:       "valid parameters",
			gatewayURL: "http://localhost:5000",
			accountID:  "U12345",
		},
		{
			name:       "empty gateway URL",
			gatewayURL: "",
			accountID:  "U12345",
		},
		{
			name:       "empty account ID",
			gatewayURL: "http://localhost:5000",
			accountID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.gatewayURL, tt.accountID)
			if client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}
