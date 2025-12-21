package service

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
)

func (s *Service) deleteLocalSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Deleting local secret '%s'\n", secretID)

	ok, err := s.storage.DeleteSecret(ctx, secretID)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("don't ok delete")
	}

	return nil
}

func (s *Service) createRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Creating remote secret '%s'\n", secretID)

	localSecret, err := s.storage.GetSecret(ctx, secretID)
	if err != nil {
		return err
	}

	remoteSecret, err := models.ConvertLocalSecretToRemoteSecret(s.cryptor, localSecret)
	if err != nil {
		return err
	}

	err = s.grpcClient.SetSecret(ctx, remoteSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createLocalSecret(ctx context.Context, remoteSecret *models.RemoteSecret) error {
	fmt.Printf("Creating local secret '%s'\n", remoteSecret.Name)

	remoteSecret, err := s.grpcClient.GetSecret(ctx, remoteSecret.Name)
	if err != nil {
		return err
	}

	localSecret, err := models.ConvertRemoteSecretToLocalSecret(s.cryptor, remoteSecret)
	if err != nil {
		return err
	}

	_, err = s.storage.CreateSecret(ctx, localSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) deleteRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Deleting remote secret '%s'\n", secretID)

	err := s.grpcClient.DeleteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) replaceLocalSecret(ctx context.Context, remoteSecret *models.RemoteSecret) error {
	fmt.Printf("Replacing remote secret '%s'\n", remoteSecret.Name)

	err := s.deleteLocalSecret(ctx, remoteSecret.Name)
	if err != nil {
		return err
	}

	err = s.createLocalSecret(ctx, remoteSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) replaceRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Replacing remote secret '%s'\n", secretID)

	err := s.deleteRemoteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	err = s.createRemoteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}
