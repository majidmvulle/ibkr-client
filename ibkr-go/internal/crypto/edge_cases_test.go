package crypto

import (
	"testing"
)

func TestEncryptDecrypt_EdgeCases(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte("")},
		{"single byte", []byte("a")},
		{"long text", []byte("this is a very long text that should be encrypted and decrypted successfully without any issues")},
		{"special chars", []byte("!@#$%^&*()_+-=[]{}|;':\",./<>?")},
		{"unicode", []byte("Hello 世界 🌍")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptToken(tt.plaintext, key)
			if err != nil {
				t.Fatalf("EncryptToken() error = %v", err)
			}

			decrypted, err := DecryptToken(encrypted, key)
			if err != nil {
				t.Fatalf("DecryptToken() error = %v", err)
			}

			if string(decrypted) != string(tt.plaintext) {
				t.Errorf("Decrypted = %v, want %v", string(decrypted), string(tt.plaintext))
			}
		})
	}
}

func TestHashToken_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"empty", ""},
		{"short", "a"},
		{"long", "this-is-a-very-long-token-string-for-testing-purposes"},
		{"special", "!@#$%^&*()"},
		{"unicode", "token-世界-🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := HashToken(tt.token)
			if len(hash) == 0 {
				t.Error("Hash should not be empty")
			}

			// Same input should produce same hash
			hash2 := HashToken(tt.token)
			if hash != hash2 {
				t.Error("Hash should be deterministic")
			}
		})
	}
}
