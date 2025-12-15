package crypto

import (
	"testing"
)

func TestHashToken_Consistency(t *testing.T) {
	token := "test-token-123"

	hash1 := HashToken(token)
	hash2 := HashToken(token)

	if hash1 != hash2 {
		t.Error("HashToken should be consistent for same input")
	}

	if len(hash1) == 0 {
		t.Error("Hash should not be empty")
	}
}

func TestHashToken_Uniqueness(t *testing.T) {
	token1 := "token1"
	token2 := "token2"

	hash1 := HashToken(token1)
	hash2 := HashToken(token2)

	if hash1 == hash2 {
		t.Error("Different tokens should produce different hashes")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	plaintext := []byte("secret message")

	encrypted, err := EncryptToken(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptToken() error = %v", err)
	}

	decrypted, err := DecryptToken(encrypted, key)
	if err != nil {
		t.Fatalf("DecryptToken() error = %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("Decrypted = %v, want %v", string(decrypted), string(plaintext))
	}
}

func TestEncryptToken_InvalidKey(t *testing.T) {
	shortKey := []byte("short")
	plaintext := []byte("message")

	_, err := EncryptToken(plaintext, shortKey)
	if err == nil {
		t.Error("Expected error with invalid key length")
	}
}

func TestDecryptToken_InvalidKey(t *testing.T) {
	key := []byte("0123456789abcdef0123456789abcdef")
	wrongKey := []byte("fedcba9876543210fedcba9876543210")
	plaintext := []byte("secret")

	encrypted, _ := EncryptToken(plaintext, key)

	_, err := DecryptToken(encrypted, wrongKey)
	if err == nil {
		t.Error("Expected error with wrong decryption key")
	}
}
