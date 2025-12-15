```go
package main

import (
	"context"
	"testing"
)

				AppName:            "test",
				HTTPPort:           8080,
				MTLSEnabled:        true,
				MTLSServerCertPath: "/tmp/cert.pem",
				MTLSServerKeyPath:  "/tmp/key.pem",
				MTLSCACertPath:     "/tmp/ca.pem",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cfg.AppName == "" {
				t.Error("AppName should not be empty")
			}
			if tt.cfg.HTTPPort == 0 {
				t.Error("HTTPPort should not be zero")
			}
		})
	}
}
