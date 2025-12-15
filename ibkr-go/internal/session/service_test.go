package session

import (
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name       string
		sessionTTL time.Duration
		wantTTL    time.Duration
	}{
		{
			name:       "with custom TTL",
			sessionTTL: 1 * time.Hour,
			wantTTL:    1 * time.Hour,
		},
		{
			name:       "with zero TTL uses default",
			sessionTTL: 0,
			wantTTL:    DefaultSessionTTL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(nil, []byte("test-key"), tt.sessionTTL)
			if service == nil {
				t.Fatal("NewService() returned nil")
			}
			if service.sessionTTL != tt.wantTTL {
				t.Errorf("NewService() sessionTTL = %v, want %v", service.sessionTTL, tt.wantTTL)
			}
		})
	}
}
