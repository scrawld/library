package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

// AesEncrypt aes encrypt
func AesEncrypt(key, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS5Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	ciphertext := make([]byte, len(plaintext))
	blockMode.CryptBlocks(ciphertext, plaintext)
	return hex.EncodeToString(ciphertext), nil
}

// AesDecrypt aes encrypt
func AesDecrypt(key []byte, ciphertext string) ([]byte, error) {
	cipherBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	plaintext := make([]byte, len(cipherBytes))
	blockMode.CryptBlocks(plaintext, cipherBytes)
	plaintext = PKCS5UnPadding(plaintext)
	return plaintext, nil
}

// PKCS5Padding padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding unpadding
func PKCS5UnPadding(plaintext []byte) []byte {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	return plaintext[:(length - unpadding)]
}
