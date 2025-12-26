package grpcclient

import (
	"context"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/grpc"

	"github.com/s-turchinskiy/keeper/internal/client/models"
)

type SenderReceiver interface {
	Connect(ctx context.Context) error
	Close() error

	GetConnectionNumber(ctx context.Context) (uint64, error)
	ConnectionNumber() uint64
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) (string, error)

	SetSecret(ctx context.Context, secret *models.RemoteSecret) error
	GetSecret(ctx context.Context, secretID string) (*models.RemoteSecret, error)
	UpdateSecret(ctx context.Context, secret *models.RemoteSecret) error
	DeleteSecret(ctx context.Context, secretID string) error
	ListSecrets(ctx context.Context) ([]*models.RemoteSecret, error)

	SyncSecretsFromClient(ctx context.Context, secrets []*models.RemoteSecret) error
	GetStream() grpc.ServerStreamingClient[proto.GetUpdatedSecretsResponse]
}
