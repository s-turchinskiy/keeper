package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/zeebo/blake3"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	storageSaltSize = 16
	keySize         = chacha20poly1305.KeySize
)

type CryptorImpl struct {
	masterPassword string
	login          string
	cachedKeys     map[string][]byte
}

func NewCryptor(masterPassword, login string) Cryptor {
	return &CryptorImpl{
		masterPassword: masterPassword,
		login:          login,
		cachedKeys:     make(map[string][]byte),
	}
}

func (c *CryptorImpl) getDeriveKey(salt []byte) []byte {
	cacheKey := base64.StdEncoding.EncodeToString(salt)

	if key, exists := c.cachedKeys[cacheKey]; exists {
		return key
	}

	key := c.genDeriveKey(salt)

	c.cachedKeys[cacheKey] = key
	return key
}

func (c *CryptorImpl) genDeriveKey(salt []byte) []byte {
	key := argon2.IDKey(
		[]byte(c.masterPassword),
		salt,
		3, 64*1024, 4, keySize,
	)
	return key
}

func (c *CryptorImpl) getSecretsKey() []byte {
	salt := []byte(c.login + "|secrets")
	return c.getDeriveKey(salt)
}

func (c *CryptorImpl) getServerKey() []byte {
	salt := []byte(c.login + "|server")
	return c.genDeriveKey(salt)
}

func (c *CryptorImpl) EncryptStorageData(plainData []byte) ([]byte, error) {
	salt := make([]byte, storageSaltSize)

	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	key := c.genDeriveKey(salt)
	encrypted, err := c.encryptWithKey(plainData, key)
	if err != nil {
		return nil, err
	}

	return append(salt, encrypted...), nil
}

func (c *CryptorImpl) DecryptStorageData(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < storageSaltSize {
		return nil, fmt.Errorf("invalid encrypted data")
	}

	salt := encryptedData[:storageSaltSize]
	ciphertext := encryptedData[storageSaltSize:]

	key := c.genDeriveKey(salt)
	return c.decryptWithKey(ciphertext, key)
}

func (c *CryptorImpl) EncryptSecretData(plainData []byte) ([]byte, error) {
	key := c.getSecretsKey()
	return c.encryptWithKey(plainData, key)
}

func (c *CryptorImpl) DecryptSecretData(encryptedData []byte) ([]byte, error) {
	key := c.getSecretsKey()
	return c.decryptWithKey(encryptedData, key)
}

func (c *CryptorImpl) CalculateDataHash(encryptedData []byte) string {
	hash := blake3.Sum256(encryptedData)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func (c *CryptorImpl) GenerateServerPassword() string {
	key := c.getServerKey()
	return base64.StdEncoding.EncodeToString(key)
}

func (c *CryptorImpl) encryptWithKey(plainData, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	encrypted := aead.Seal(nonce, nonce, plainData, nil)
	return encrypted, nil
}

func (c *CryptorImpl) decryptWithKey(encryptedData, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plainData, nil
}
