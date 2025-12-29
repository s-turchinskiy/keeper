package models

const (
	SecretTypePassword = "password"
	SecretTypeText     = "text"
	SecretTypeBinary   = "binary"
	SecretTypeCard     = "card"
)

const (
	MaxFileSize         = 2 * 1024 * 1024 // 2 MB
	MaxTextSize         = 1 * 1024 * 1024 // 1 MB
	MaxCardHolderLength = 100
	MaxUsernameLength   = 255
	MaxPasswordLength   = 1024
	MaxURLLength        = 2048
)
