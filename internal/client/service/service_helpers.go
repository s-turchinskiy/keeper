package service

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"github.com/s-turchinskiy/keeper/models/proto"
)

func (s *Service) deleteLocalSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Deleting local secret '%s'\n", secretID)

	err := s.storage.DeleteByKey(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) createRemoteSecret(ctx context.Context, secretID string) error {
	fmt.Printf("Creating remote secret '%s'\n", secretID)

	localSecret, err := s.storage.GetByKey(ctx, secretID)
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

	_, err = s.storage.Create(ctx, localSecret)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getSecretFromServer(ctx context.Context, remoteSecretID string) (*models.LocalSecret, error) {
	fmt.Printf("Get remote secret '%s'\n", remoteSecretID)

	remoteSecret, err := s.grpcClient.GetSecret(ctx, remoteSecretID)
	if err != nil {
		return nil, err
	}

	localSecret, err := models.ConvertRemoteSecretToLocalSecret(s.cryptor, remoteSecret)
	if err != nil {
		return nil, err
	}

	return localSecret, nil
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

func (s *Service) replaceRemoteSecret(ctx context.Context, localSecret *models.LocalSecret) error {
	fmt.Printf("Replacing remote secret '%s'\n", localSecret.Name)

	remoteSecret, err := models.ConvertLocalSecretToRemoteSecret(s.cryptor, localSecret)
	if err != nil {
		return err
	}

	err = s.grpcClient.UpdateSecret(ctx, remoteSecret)
	if err != nil {
		return err
	}

	return nil

}

func (s *Service) syncLocalSecret(ctx context.Context, secret *proto.Secret, connNumber string) error {

	if *secret.Deleted {

		fmt.Printf("conn %s. GetUpdatedSecrets start deleting secret \"%s\"\n", connNumber, secret.Id)
		err := s.deleteLocalSecret(ctx, secret.Id)
		fmt.Printf("conn %s. GetUpdatedSecrets end deleting secret \"%s\"\n", connNumber, secret.Id)
		return err

	}

	fmt.Printf("conn %s. GetUpdatedSecrets start updating secret \"%s\"\n", connNumber, secret.Id)

	localSecret, err := s.storage.GetByKey(ctx, secret.Id)
	if err != nil {
		fmt.Printf("conn %s. GetUpdatedSecrets end updating secret \"%s\", error: %v\n",
			connNumber, secret.Id, errorsutils.WrapError(err))
		return err
	}

	remoteSecret := models.ConvertProtoSecretToRemoteSecret(secret)
	if localSecret.LastModified.Before(remoteSecret.LastModified) {
		err = s.replaceLocalSecret(ctx, remoteSecret)
		if err != nil {
			fmt.Printf("conn %s. GetUpdatedSecrets end updating secret \"%s\", error: %v\n",
				connNumber, secret.Id, errorsutils.WrapError(err))
			return err
		}

		fmt.Printf("conn %s. GetUpdatedSecrets end updating secret \"%s\", success\n", connNumber, secret.Id)
		return nil
	}

	fmt.Printf("conn %s. GetUpdatedSecrets end updating secret \"%s\", LastModified equal\n", connNumber, secret.Id)
	return nil
}
