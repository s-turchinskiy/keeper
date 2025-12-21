package service

import (
	"context"
	"log"
	"sync"

	"github.com/s-turchinskiy/keeper/internal/client/crypto"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	"github.com/s-turchinskiy/keeper/internal/client/repository"
)

type Service struct {
	cryptor crypto.Cryptor

	mu sync.Mutex

	storage    repository.Repositorier
	grpcClient grpcclient.SenderReceiver
}

func NewService(ctx context.Context, cryptor crypto.Cryptor, storage repository.Repositorier, grpcClient *grpcclient.GRPCClient) *Service {
	service := &Service{
		cryptor:    cryptor,
		storage:    storage,
		grpcClient: grpcClient,
	}

	return service
}

func (s *Service) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.storage != nil {
		if err := s.storage.Close(ctx); err != nil {
			log.Printf("failed to close repository: %v", err)
		}
	}

	if s.grpcClient != nil {
		if err := s.grpcClient.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}

	return nil
}
