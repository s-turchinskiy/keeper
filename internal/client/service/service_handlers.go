package service

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"io"
	"log"
)

func (s *Service) Register(ctx context.Context, login, password string) error {

	_, err := s.grpcClient.Register(ctx, login, password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) Login(ctx context.Context, login, password string) error {

	err := s.grpcClient.Login(ctx, login, password)
	if err != nil {
		return err
	}

	go func() {
		err = s.GetUpdatedSecrets(ctx)
		if err != nil {
			log.Printf("GetUpdatedSecrets failed err: %v\n", err)
		} else {
			log.Printf("GetUpdatedSecrets run succsess\n")
		}
	}()

	return nil
}

func (s *Service) SyncSecrets(ctx context.Context) error {

	localSecrets, err := s.storage.ListSecrets(ctx)
	if err != nil {
		return err
	}

	remoteSecrets := make([]*models.RemoteSecret, len(localSecrets))
	for _, localSecret := range localSecrets {
		remoteSecret, err := models.ConvertLocalSecretToRemoteSecret(s.cryptor, localSecret)
		if err != nil {
			return nil
		}
		remoteSecrets = append(remoteSecrets, remoteSecret)
	}
	err = s.grpcClient.SyncSecretsFromClient(ctx, remoteSecrets)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateSecret(ctx context.Context, base models.BaseSecret, data models.SecretData) (*models.LocalSecret, error) {

	secret, err := models.NewSecretModel(base, data, s.cryptor)
	if err != nil {
		return nil, err
	}

	_, err = s.storage.CreateSecret(ctx, secret)
	if err != nil {
		return nil, err
	}

	err = s.createRemoteSecret(ctx, secret.Name)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s *Service) UpdateSecret(ctx context.Context, secret *models.LocalSecret) error {

	_, err := s.storage.UpdateSecret(ctx, secret)
	if err != nil {
		return err
	}

	err = s.replaceRemoteSecret(ctx, secret.Name)
	return err
}

func (s *Service) ReadSecret(ctx context.Context, secretID string) (*models.LocalSecret, error) {

	//закомментировано для теста метода GetSecret сервера
	/*secret, err := s.storage.GetSecret(ctx, secretID)
	if err != nil {
		return nil, err
	}*/

	secret, err := s.getSecretFromServer(ctx, secretID)

	return secret, err
}

func (s *Service) ListLocalSecrets(ctx context.Context) ([]*models.LocalSecret, error) {

	secrets, err := s.storage.ListSecrets(ctx)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *Service) DeleteSecret(ctx context.Context, secretID string) (bool, error) {

	ok, err := s.storage.DeleteSecret(ctx, secretID)
	if err != nil {
		return false, err
	}

	err = s.grpcClient.DeleteSecret(ctx, secretID)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (s *Service) GetUpdatedSecrets(ctx context.Context) error {

	fmt.Printf("GetUpdatedSecrets starting run\n")

	var err error
	for {

		resp, err := s.grpcClient.GetStream().Recv()

		log.Printf("GetUpdatedSecrets start getting secrets %v\n", resp.Secrets)

		if err == io.EOF {
			log.Printf("GetUpdatedSecrets stopped\n")
			break
		}
		if err != nil {
			log.Printf("GetUpdatedSecrets error: %v\n", err)
			break
		}

		for _, secret := range resp.Secrets {
			if *secret.Deleted {
				err = s.deleteLocalSecret(ctx, secret.Id)
				if err != nil {
					return nil
				}
			} else {

				localSecret, err := s.storage.GetSecret(ctx, secret.Id)
				if err != nil {
					return nil
				}

				remoteSecret := models.ConvertProtoSecretToRemoteSecret(secret)
				if localSecret.LastModified.Before(remoteSecret.LastModified) {
					err = s.replaceLocalSecret(ctx, remoteSecret)
					if err != nil {
						return nil
					}
				}
			}
		}

		log.Printf("GetUpdatedSecrets end getting secrets\n")

	}

	fmt.Printf("GetUpdatedSecrets stopped\n")
	return err
}
