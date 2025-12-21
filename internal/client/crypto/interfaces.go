package crypto

type Cryptor interface {
	EncryptStorageData(plainData []byte) ([]byte, error)
	DecryptStorageData(encryptedData []byte) ([]byte, error)

	EncryptSecretData(plainData []byte) ([]byte, error)
	DecryptSecretData(encryptedData []byte) ([]byte, error)

	CalculateDataHash(data []byte) string

	GenerateServerPassword() string
}
