package repository

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/models"
)

type Repositorier interface {
	Create(ctx context.Context, secret *models.LocalSecret) (*models.LocalSecret, error)
	GetAll(ctx context.Context) ([]*models.LocalSecret, error)
	GetByKey(ctx context.Context, id string) (*models.LocalSecret, error)
	UpdateByKey(ctx context.Context, id string, secret *models.LocalSecret) (*models.LocalSecret, error)
	DeleteByKey(ctx context.Context, id string) error

	Close(ctx context.Context) error
}
