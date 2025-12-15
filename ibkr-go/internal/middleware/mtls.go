package middleware

import (
	"context"
	"crypto/x509"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
)

// ClientIdentityContextKey is the context key for storing the client identity from mTLS.
type ClientIdentityContextKey struct{}

// NewMTLSInterceptor creates a ConnectRPC interceptor that validates mTLS client certificates.
// It extracts the client identity from the certificate's Common Name (CN) and stores it in the context.
func NewMTLSInterceptor(logger *slog.Logger) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// Extract client certificate from peer (set by TLS handshake).
			peer := req.Peer()
			if peer.Addr == "" {
				logger.Warn("No peer information available")
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no client certificate provided"))
			}

			// In a real mTLS setup, the certificate is validated during the TLS handshake.
			// The server's TLS config (ClientAuth: RequireAndVerifyClientCert) ensures
			// that only valid certificates signed by the trusted CA are accepted.
			//
			// At this point, we know the certificate is valid. We extract the identity
			// from the certificate's Common Name or Subject Alternative Names.
			//
			// Note: In Go's http.Server with TLS, the verified certificate chain is
			// available in the request's TLS.PeerCertificates field. However, Connect-RPC
			// doesn't expose this directly. We'll need to extract it from the underlying
			// HTTP request in the server setup.

			// For now, we'll add a placeholder that will be populated by the server's
			// TLS verification. The actual certificate validation happens at the TLS layer.
			clientIdentity := extractClientIdentity(ctx)
			if clientIdentity == "" {
				logger.Warn("Failed to extract client identity from certificate")
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid client certificate"))
			}

			logger.Info("mTLS authentication successful",
				slog.String("client", clientIdentity),
				slog.String("procedure", req.Spec().Procedure),
			)

			// Add client identity to context for downstream handlers.
			ctx = context.WithValue(ctx, ClientIdentityContextKey{}, clientIdentity)

			return next(ctx, req)
		}
	}
}

// extractClientIdentity extracts the client identity from the TLS certificate.
// This is a placeholder - the actual implementation will extract from the HTTP request's
// TLS.PeerCertificates field.
func extractClientIdentity(ctx context.Context) string {
	// This will be implemented when we set up the TLS server.
	// For now, return empty to indicate we need to wire this up.
	return ""
}

// ExtractClientIdentityFromCert extracts the Common Name from a client certificate.
func ExtractClientIdentityFromCert(cert *x509.Certificate) string {
	if cert == nil {
		return ""
	}

	// Use Common Name as the primary identifier.
	if cert.Subject.CommonName != "" {
		return cert.Subject.CommonName
	}

	// Fallback to first DNS name in SANs.
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}

	return ""
}
