package middleware

import (
	"context"
	"testing"
)

func TestGetAccountIDFromContext(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		want   string
		wantOk bool
	}{
		{
			name:   "context with account ID",
			ctx:    context.WithValue(context.Background(), ClientIdentityContextKey{}, "U12345"),
			want:   "U12345",
			wantOk: true,
		},
		{
			name:   "context without account ID",
			ctx:    context.Background(),
			want:   "",
			wantOk: false,
		},
		{
			name:   "context with wrong type",
			ctx:    context.WithValue(context.Background(), ClientIdentityContextKey{}, 12345),
			want:   "",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetAccountIDFromContext(tt.ctx)
			if got != tt.want {
				t.Errorf("GetAccountIDFromContext() got = %v, want %v", got, tt.want)
			}
			if ok != tt.wantOk {
				t.Errorf("GetAccountIDFromContext() ok = %v, want %v", ok, tt.wantOk)
			}
		})
	}
}
