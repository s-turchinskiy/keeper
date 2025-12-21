package service

import (
	"context"
	"errors"
	"github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"sync/atomic"
)

const maxSecretSize = 5 * 1024 * 1024 // 5MB

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrSecretTooLarge       = errors.New("secret data too large")
	ErrUnknownTypeOperation = errors.New("unknown type operation")
)

func (s *Service) GetNewConnectionNumber(ctx context.Context) uint64 {
	return atomic.AddUint64(&s.currentConnectionNumber, 1)

}

func (s *Service) Register(ctx context.Context, login, password string) (*models.User, error) {

	existing, _ := s.usersRepository.GetUserByLogin(ctx, login)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.usersRepository.CreateUser(ctx, login, string(passwordHash))
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, *models.User, error) {
	user, err := s.usersRepository.GetUserByLogin(ctx, login)
	if err != nil {
		return "", nil, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	token, err := s.TokenManager.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *Service) SetSecret(ctx context.Context, secret *models.Secret) error {
	if len(secret.Data) > maxSecretSize {
		return ErrSecretTooLarge
	}

	return s.secretRepository.SetSecret(ctx, secret)
}

func (s *Service) GetSecret(ctx context.Context, userID, secretID string) (*models.Secret, error) {
	return s.secretRepository.GetSecret(ctx, userID, secretID)
}

func (s *Service) DeleteSecret(ctx context.Context, userID, secretID string) error {
	return s.secretRepository.DeleteSecret(ctx, userID, secretID)
}

func (s *Service) ListSecrets(ctx context.Context, userID string) ([]*models.Secret, error) {
	return s.secretRepository.ListSecrets(ctx, userID)
}

func (s *Service) SyncFromClient(ctx context.Context, userID string, clientSecrets []*models.Secret) ([]*models.Secret, error) {

	var updateInClients []*models.Secret

	serverSecrets, err := s.secretRepository.ListSecretsWithStatuses(ctx, userID)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	type secretsType struct {
		clientSecret *models.Secret
		serverSecret *models.Secret
	}

	lenSecrets := len(serverSecrets)
	if len(serverSecrets) > lenSecrets {
		lenSecrets = len(serverSecrets)
	}

	comparisonMap := make(map[string]*secretsType, lenSecrets)
	for _, secret := range serverSecrets {
		comparisonMap[secret.ID] = &secretsType{serverSecret: secret}
	}

	for _, secret := range clientSecrets {

		val, ok := comparisonMap[secret.ID]
		if !ok {
			comparisonMap[secret.ID] = &secretsType{clientSecret: secret}
		} else {
			val.clientSecret = secret
			comparisonMap[secret.ID] = val
		}
	}

	var register bool
	for _, scrt := range comparisonMap {

		register = true
		switch {
		case (scrt.serverSecret == nil && scrt.clientSecret.Deleted) ||
			(scrt.clientSecret == nil && scrt.serverSecret.Deleted) ||
			(scrt.clientSecret.LastModified == scrt.serverSecret.LastModified):

			register = false
			continue

		case scrt.clientSecret == nil:

			updateInClients = append(updateInClients, scrt.serverSecret)
			register = false

		case scrt.serverSecret == nil:

			err = s.secretRepository.SetSecret(ctx, scrt.clientSecret)
			if err != nil {
				return updateInClients, errorsutils.WrapError(err)
			}

		case scrt.clientSecret.LastModified.After(scrt.serverSecret.LastModified):

			if scrt.clientSecret.Deleted {
				err = s.secretRepository.DeleteSecret(ctx, scrt.clientSecret.UserID, scrt.clientSecret.ID)
				if err != nil {
					return updateInClients, errorsutils.WrapError(err)
				}

			} else {
				err = s.secretRepository.SetSecret(ctx, scrt.clientSecret)
				if err != nil {
					return updateInClients, errorsutils.WrapError(err)
				}
			}
		case scrt.serverSecret.LastModified.After(scrt.clientSecret.LastModified):
			updateInClients = append(updateInClients, scrt.serverSecret)
			register = false
		default:
			log.Println(ErrUnknownTypeOperation, "server secret:", scrt.serverSecret, "client secret:", scrt.clientSecret)
			return updateInClients, errorsutils.WrapError(ErrUnknownTypeOperation)
		}

		if register {
			updatedSecret, err := s.secretRepository.GetSecret(ctx, scrt.clientSecret.UserID, scrt.clientSecret.ID)
			if err != nil {
				return updateInClients, errorsutils.WrapError(err)
			}
			updateInClients = append(updateInClients, updatedSecret)
		}
	}

	return updateInClients, nil
}
