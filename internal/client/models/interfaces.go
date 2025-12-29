package models

type SecretData interface {
	Validate() error
}
