package service

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/models"
)

type Servicer interface {
	Register(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) error

	SyncSecrets(ctx context.Context) error
	CreateSecret(ctx context.Context, base models.BaseSecret, data models.SecretData) (*models.LocalSecret, error)
	ReadSecret(ctx context.Context, secretID string) (*models.LocalSecret, error)
	UpdateSecret(ctx context.Context, secret *models.LocalSecret) error
	DeleteSecret(ctx context.Context, secretID string) error
	ListLocalSecrets(ctx context.Context) ([]*models.LocalSecret, error)

	Close(ctx context.Context) error
}
