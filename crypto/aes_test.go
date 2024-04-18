package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEncryptDecrypt(t *testing.T) {
	assert := assert.New(t)

	key := []byte("examplekey123456") // 16-byte key
	plaintext := []byte("Hello, world!")

	// Encrypt
	ciphertext, err := AesEncrypt(key, plaintext)
	assert.NoError(err, "Encryption should not return an error")

	// Decrypt
	decryptedPlaintext, err := AesDecrypt(key, ciphertext)
	assert.NoError(err, "Decryption should not return an error")

	// Verify
	assert.Equal(plaintext, decryptedPlaintext, "Decrypted plaintext should match original plaintext")
}
