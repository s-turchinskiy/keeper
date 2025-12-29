package service

import (
	"context"
	"errors"
	"github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
	"sync/atomic"
)

type secretsType struct {
	clientSecret *models.Secret
	serverSecret *models.Secret
}

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

	existing, _ := s.usersRepository.GetByLogin(ctx, login)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.usersRepository.Create(ctx, login, string(passwordHash))
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *Service) Login(ctx context.Context, login, password string) (string, *models.User, error) {
	user, err := s.usersRepository.GetByLogin(ctx, login)
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

func (s *Service) CreateSecret(ctx context.Context, secret *models.Secret) error {
	if len(secret.Data) > maxSecretSize {
		return ErrSecretTooLarge
	}

	return s.secretRepository.CreateUpdate(ctx, secret)
}

func (s *Service) GetSecret(ctx context.Context, userID, secretID string) (*models.Secret, error) {

	if s.redisClient != nil {
		secret, _ := s.redisClient.Get(ctx, userID, secretID)
		if secret != nil {
			return secret, nil
		}
	}

	secret, err := s.secretRepository.GetByID(ctx, userID, secretID)

	if s.redisClient != nil {
		_ = s.redisClient.Set(ctx, secret)
	}

	return secret, err
}

func (s *Service) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	if len(secret.Data) > maxSecretSize {
		return ErrSecretTooLarge
	}

	return s.secretRepository.CreateUpdate(ctx, secret)
}

func (s *Service) DeleteSecret(ctx context.Context, userID, secretID string) error {
	return s.secretRepository.Delete(ctx, userID, secretID)
}

func (s *Service) ListSecrets(ctx context.Context, userID string) ([]*models.Secret, error) {
	return s.secretRepository.GetAll(ctx, userID)
}

func (s *Service) SyncFromClient(ctx context.Context, userID string, clientSecrets []*models.Secret) ([]*models.Secret, error) {

	var updateInClients []*models.Secret

	serverSecrets, err := s.secretRepository.GetAllWithStatuses(ctx, userID)
	if err != nil {
		return nil, errorsutils.WrapError(err)
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

	grp, ctx := errgroup.WithContext(ctx)
	var mutex *sync.Mutex

	for _, scrt := range comparisonMap {
		scrt := scrt
		grp.Go(func() error {
			secretForClient, err := s.syncSecret(ctx, scrt)
			if err == nil && secretForClient != nil {
				mutex.Lock()
				defer mutex.Unlock()
				updateInClients = append(updateInClients, secretForClient)
			}
			return err
		})
	}

	if err = grp.Wait(); err != nil {
		return updateInClients, err
	}

	return updateInClients, nil
}

func (s *Service) syncSecret(ctx context.Context, scrt *secretsType) (*models.Secret, error) {

	switch {
	case (scrt.serverSecret == nil && scrt.clientSecret.Deleted) ||
		(scrt.clientSecret == nil && scrt.serverSecret.Deleted) ||
		scrt.clientSecret.LastModified.Equal(scrt.serverSecret.LastModified):

		return nil, nil

	case scrt.clientSecret == nil:

		return scrt.serverSecret, nil

	case scrt.serverSecret == nil:

		err := s.secretRepository.CreateUpdate(ctx, scrt.clientSecret)
		if err != nil {
			return nil, errorsutils.WrapError(err)
		}

	case scrt.clientSecret.LastModified.After(scrt.serverSecret.LastModified):

		if scrt.clientSecret.Deleted {
			err := s.secretRepository.Delete(ctx, scrt.clientSecret.UserID, scrt.clientSecret.ID)
			if err != nil {
				return nil, errorsutils.WrapError(err)
			}

		} else {
			err := s.secretRepository.CreateUpdate(ctx, scrt.clientSecret)
			if err != nil {
				return nil, errorsutils.WrapError(err)
			}
		}

	case scrt.serverSecret.LastModified.After(scrt.clientSecret.LastModified):

		return scrt.serverSecret, nil
	default:
		log.Println(ErrUnknownTypeOperation, "server secret:", scrt.serverSecret, "client secret:", scrt.clientSecret)
		return nil, errorsutils.WrapError(ErrUnknownTypeOperation)
	}

	updatedSecret, err := s.secretRepository.GetByID(ctx, scrt.clientSecret.UserID, scrt.clientSecret.ID)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	return updatedSecret, nil
}
