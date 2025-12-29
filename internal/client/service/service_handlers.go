package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"strconv"
)

var ErrSecretAlreadyExist = errors.New("secret already exist")

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

	localSecrets, err := s.storage.GetAll(ctx)
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

	_, err = s.storage.GetByKey(ctx, secret.Name)
	if err == nil {
		return nil, ErrSecretAlreadyExist
	}

	_, err = s.storage.Create(ctx, secret)
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

	secret, err := s.storage.UpdateByKey(ctx, secret.Name, secret)
	if err != nil {
		return err
	}

	err = s.replaceRemoteSecret(ctx, secret)
	return err
}

func (s *Service) ReadSecret(ctx context.Context, secretID string) (*models.LocalSecret, error) {

	_, err := s.storage.GetByKey(ctx, secretID)
	if err != nil {
		return nil, err
	}

	secret, err := s.getSecretFromServer(ctx, secretID)

	return secret, err
}

func (s *Service) ListLocalSecrets(ctx context.Context) ([]*models.LocalSecret, error) {

	secrets, err := s.storage.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *Service) DeleteSecret(ctx context.Context, secretID string) error {

	err := s.deleteLocalSecret(ctx, secretID)
	if err != nil {
		return err
	}

	err = s.deleteRemoteSecret(ctx, secretID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUpdatedSecrets(ctx context.Context) error {

	connNumber := strconv.FormatUint(s.grpcClient.ConnectionNumber(), 10)
	fmt.Printf("conn %s. GetUpdatedSecrets starting run\n", connNumber)

	var err error
	for {

		resp, err := s.grpcClient.GetStream().Recv()

		fmt.Printf("conn %s. GetUpdatedSecrets start getting secrets %v\n", connNumber, resp.Secrets)

		if err == io.EOF {
			fmt.Printf("conn %s. GetUpdatedSecrets stopped EOF\n", connNumber)
			break
		}
		if err != nil {
			fmt.Printf("conn %s. GetUpdatedSecrets error: %v\n", connNumber, err)
			break
		}

		grp, ctx := errgroup.WithContext(ctx)
		for _, secret := range resp.Secrets {
			secret := secret
			grp.Go(func() error {
				return s.syncLocalSecret(ctx, secret, connNumber)
			})
		}

		if err = grp.Wait(); err != nil {
			return err
		}

		fmt.Printf("conn %s. GetUpdatedSecrets end getting secrets\n", connNumber)

	}

	fmt.Printf("conn %s. GetUpdatedSecrets stopped\n", connNumber)
	return err
}
