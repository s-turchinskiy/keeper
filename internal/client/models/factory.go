package models

import (
	"fmt"
	"time"

	"github.com/s-turchinskiy/keeper/internal/client/crypto"
)

func NewSecretModel(base BaseSecret, data SecretData, cryptor crypto.Cryptor) (*LocalSecret, error) {

	if err := data.Validate(); err != nil {
		return nil, fmt.Errorf("data validation failed: %w", err)
	}

	secret := &LocalSecret{
		Name:         base.Name,
		Type:         base.Type,
		LastModified: time.Now().Truncate(time.Microsecond),
		Metadata:     base.Metadata,
	}

	if err := secret.SetData(cryptor, data); err != nil {
		return nil, err
	}
	return secret, nil

}
