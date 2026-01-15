package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// KeySize is the size of the key for NaCl SecretBox (32 bytes)
const KeySize = 32

// NonceSize is the size of the nonce (24 bytes)
const NonceSize = 24

// DeriveKey generates a 32-byte key from a password using SHA256.
// In production, Argon2 or Scrypt would be better, but SHA256 is fast for this MVP.
func DeriveKey(password string) [KeySize]byte {
	return sha256.Sum256([]byte(password))
}

// Encrypt encrypts a message using a shared password.
// It prepends the generated nonce to the output.
func Encrypt(message []byte, password string) ([]byte, error) {
	key := DeriveKey(password)
	var nonce [NonceSize]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}

	// Overhead is appended, but here we prepend nonce manually to transport it
	encrypted := secretbox.Seal(nonce[:], message, &nonce, &key)
	return encrypted, nil
}

// Decrypt decrypts a message that has the nonce prepended.
func Decrypt(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < NonceSize {
		return nil, errors.New("ciphertext too short")
	}

	key := DeriveKey(password)
	var nonce [NonceSize]byte
	copy(nonce[:], encrypted[:NonceSize])

	// The message payload follows the nonce
	payload := encrypted[NonceSize:]

	decrypted, ok := secretbox.Open(nil, payload, &nonce, &key)
	if !ok {
		return nil, errors.New("decryption failed")
	}

	return decrypted, nil
}
