package middleware

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"
)

func TestExtractClientIdentityFromCert_AllCases(t *testing.T) {
	tests := []struct {
		name string
		cert *x509.Certificate
		want string
	}{
		{
			name: "nil cert",
			cert: nil,
			want: "",
		},
		{
			name: "cert with CN",
			cert: &x509.Certificate{
				Subject: pkix.Name{CommonName: "test-client"},
			},
			want: "test-client",
		},
		{
			name: "cert with DNS names only",
			cert: &x509.Certificate{
				Subject:  pkix.Name{CommonName: ""},
				DNSNames: []string{"client.example.com"},
			},
			want: "client.example.com",
		},
		{
			name: "cert with both CN and DNS",
			cert: &x509.Certificate{
				Subject:  pkix.Name{CommonName: "cn-name"},
				DNSNames: []string{"dns-name.com"},
			},
			want: "cn-name",
		},
		{
			name: "cert with empty CN and empty DNS",
			cert: &x509.Certificate{
				Subject:  pkix.Name{CommonName: ""},
				DNSNames: []string{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractClientIdentityFromCert(tt.cert)
			if got != tt.want {
				t.Errorf("ExtractClientIdentityFromCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAccountIDFromContext_AllCases(t *testing.T) {
	tests := []struct {
		name   string
		ctx    context.Context
		wantID string
		wantOK bool
	}{
		{
			name:   "context with valid ID",
			ctx:    context.WithValue(context.Background(), ClientIdentityContextKey{}, "U12345"),
			wantID: "U12345",
			wantOK: true,
		},
		{
			name:   "context without ID",
			ctx:    context.Background(),
			wantID: "",
			wantOK: false,
		},
		{
			name:   "context with wrong type",
			ctx:    context.WithValue(context.Background(), ClientIdentityContextKey{}, 12345),
			wantID: "",
			wantOK: false,
		},
		{
			name:   "context with empty string",
			ctx:    context.WithValue(context.Background(), ClientIdentityContextKey{}, ""),
			wantID: "",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotOK := GetAccountIDFromContext(tt.ctx)
			if gotID != tt.wantID {
				t.Errorf("GetAccountIDFromContext() ID = %v, want %v", gotID, tt.wantID)
			}
			if gotOK != tt.wantOK {
				t.Errorf("GetAccountIDFromContext() OK = %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}
