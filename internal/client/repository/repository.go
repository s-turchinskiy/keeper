package repository

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/models"
)

type Repositorier interface {
	CreateSecret(ctx context.Context, secret *models.LocalSecret) (secretID string, err error)
	GetSecret(ctx context.Context, secretID string) (*models.LocalSecret, error)
	UpdateSecret(ctx context.Context, secret *models.LocalSecret) (bool, error)
	DeleteSecret(ctx context.Context, secretID string) (bool, error)
	ListSecrets(ctx context.Context) ([]*models.LocalSecret, error)

	Close(ctx context.Context) error
}
