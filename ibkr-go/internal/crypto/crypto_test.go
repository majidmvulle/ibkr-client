package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptToken(t *testing.T) {
	validKey := []byte("12345678901234567890123456789012") // 32 bytes for AES-256

	tests := []struct {
		name      string
		plaintext []byte
		key       []byte
		wantErr   bool
	}{
		{
			name:      "successful encryption and decryption",
			plaintext: []byte("Hello, World!"),
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "empty plaintext",
			plaintext: []byte(""),
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "long plaintext",
			plaintext: []byte("This is a very long plaintext message that should still be encrypted and decrypted correctly without any issues."),
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "special characters",
			plaintext: []byte("!@#$%^&*()_+-=[]{}|;:',.<>?/~`"),
			key:       validKey,
			wantErr:   false,
		},
		{
			name:      "unicode characters",
			plaintext: []byte("Hello 世界 🌍"),
			key:       validKey,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			ciphertext, err := EncryptToken(tt.plaintext, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify ciphertext is not empty and different from plaintext
			if len(ciphertext) == 0 {
				t.Error("EncryptToken() returned empty ciphertext")
			}
			if bytes.Equal(ciphertext, tt.plaintext) {
				t.Error("EncryptToken() returned plaintext as ciphertext")
			}

			// Decrypt
			decrypted, err := DecryptToken(ciphertext, tt.key)
			if err != nil {
				t.Errorf("DecryptToken() error = %v", err)
				return
			}

			// Verify round-trip
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("DecryptToken() = %v, want %v", string(decrypted), string(tt.plaintext))
			}
		})
	}
}

func TestEncryptTokenErrors(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
		key       []byte
		wantErr   bool
	}{
		{
			name:      "invalid key length - too short",
			plaintext: []byte("test"),
			key:       []byte("short"),
			wantErr:   true,
		},
		{
			name:      "invalid key length - too long",
			plaintext: []byte("test"),
			key:       []byte("123456789012345678901234567890123"), // 33 bytes
			wantErr:   true,
		},
		{
			name:      "empty key",
			plaintext: []byte("test"),
			key:       []byte(""),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := EncryptToken(tt.plaintext, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecryptTokenErrors(t *testing.T) {
	validKey := []byte("12345678901234567890123456789012")
	validCiphertext, _ := EncryptToken([]byte("secret"), validKey)

	tests := []struct {
		name       string
		ciphertext []byte
		key        []byte
		wantErr    bool
	}{
		{
			name:       "wrong key",
			ciphertext: validCiphertext,
			key:        []byte("00000000000000000000000000000000"),
			wantErr:    true,
		},
		{
			name:       "too short ciphertext",
			ciphertext: []byte("short"),
			key:        validKey,
			wantErr:    true,
		},
		{
			name:       "empty ciphertext",
			ciphertext: []byte(""),
			key:        validKey,
			wantErr:    true,
		},
		{
			name:       "invalid key length",
			ciphertext: validCiphertext,
			key:        []byte("short"),
			wantErr:    true,
		},
		{
			name:       "corrupted ciphertext",
			ciphertext: append([]byte{}, validCiphertext[:len(validCiphertext)-1]...), // Truncate
			key:        validKey,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptToken(tt.ciphertext, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecryptToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHashToken(t *testing.T) {
	tests := []struct {
		name string
		data string
		want string // Expected hash (SHA-256)
	}{
		{
			name: "simple string",
			data: "hello",
			want: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
		},
		{
			name: "empty string",
			data: "",
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashToken(tt.data)

			// Verify hash is not empty
			if got == "" {
				t.Error("HashToken() returned empty string")
			}

			// Verify hash length (SHA-256 produces 64 hex characters)
			if len(got) != 64 {
				t.Errorf("HashToken() length = %d, want 64", len(got))
			}

			// For known hashes, verify exact match
			if got != tt.want {
				t.Errorf("HashToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashTokenConsistency(t *testing.T) {
	data := "test data for consistency"

	hash1 := HashToken(data)
	hash2 := HashToken(data)

	if hash1 != hash2 {
		t.Errorf("HashToken() not consistent: %v != %v", hash1, hash2)
	}
}

func TestHashTokenUniqueness(t *testing.T) {
	data1 := "data1"
	data2 := "data2"

	hash1 := HashToken(data1)
	hash2 := HashToken(data2)

	if hash1 == hash2 {
		t.Error("HashToken() produced same hash for different inputs")
	}
}

func TestHashTokenDifferentInputs(t *testing.T) {
	inputs := []string{
		"test",
		"Test",
		"TEST",
		"test ",
		" test",
		"test1",
		"test2",
	}

	hashes := make(map[string]bool)
	for _, input := range inputs {
		hash := HashToken(input)
		if hashes[hash] {
			t.Errorf("Duplicate hash for input: %s", input)
		}
		hashes[hash] = true
	}
}

func TestEncryptDecryptLargeData(t *testing.T) {
	key := []byte("12345678901234567890123456789012")

	// Test with 1MB of data
	largeData := make([]byte, 1024*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	ciphertext, err := EncryptToken(largeData, key)
	if err != nil {
		t.Fatalf("EncryptToken() failed for large data: %v", err)
	}

	decrypted, err := DecryptToken(ciphertext, key)
	if err != nil {
		t.Fatalf("DecryptToken() failed for large data: %v", err)
	}

	if !bytes.Equal(decrypted, largeData) {
		t.Error("Large data round-trip failed")
	}
}
