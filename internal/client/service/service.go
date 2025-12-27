package service

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/crypto"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	"github.com/s-turchinskiy/keeper/internal/client/repository"
	"log"
	"sync"
)

type OptionService func(*Service)

type Service struct {
	cryptor crypto.Cryptor

	mu sync.Mutex

	storage    repository.Repositorier
	grpcClient grpcclient.SenderReceiver
}

func NewService(ctx context.Context, storage repository.Repositorier, grpcClient *grpcclient.GRPCClient, opts ...OptionService) *Service {
	service := &Service{
		storage:    storage,
		grpcClient: grpcClient,
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func WithCrypto(cryptor crypto.Cryptor) OptionService {

	return func(s *Service) {

		s.cryptor = cryptor
	}
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
