package functional_tests

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	clientmodels "github.com/s-turchinskiy/keeper/internal/client/models"
	grpcserver "github.com/s-turchinskiy/keeper/internal/server/grpc"
	servermodels "github.com/s-turchinskiy/keeper/internal/server/models"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	mockserverrepository "github.com/s-turchinskiy/keeper/internal/server/repository/mock"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/internal/server/token"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"testing"
	"time"
)

func FunctionalTestGRPC(t *testing.T, usersRepository repository.UserRepositorier, secretRepository repository.SecretRepositorier, debug bool) {

	lis = bufconn.Listen(bufSize)

	jwtManager := token.NewJWTManager("secret", time.Minute)
	srvc := service.NewService(
		jwtManager,
		usersRepository,
		secretRepository,
	)
	grpcServer := grpcserver.NewGrpcServer(srvc)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	if !debug {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	clientSecrets := []*clientmodels.RemoteSecret{
		{
			Name:         secretOnlyCreating,
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
		{
			Name:         secretForDeleting,
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
		{
			Name:         secretForUpdating,
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
	}

	if debug {
		functionalTestGRPC(ctx, t, clientSecrets, debug)
	} else {
		t.Run(fmt.Sprintf("test #%d: %s", 1, "Функциональный тест grpc на все методы"), func(t *testing.T) {
			functionalTestGRPC(ctx, t, clientSecrets, debug)
		})
	}

}

func functionalTestGRPC(ctx context.Context, t *testing.T, clientSecrets []*clientmodels.RemoteSecret, debug bool) {

	grpcClient, err := grpcclient.NewGRPCClient(ctx, "passthrough://bufnet", loginExistingUser, password, grpc.WithContextDialer(bufDialer))
	require.NoError(t, err)

	_, err = grpcClient.Register(ctx, loginNewUser, password)
	if !debug {
		require.NoError(t, err)
	}

	_, err = grpcClient.Register(ctx, loginExistingUser, password)
	require.Error(t, err)

	err = grpcClient.Login(ctx, loginExistingUser, password)
	require.NoError(t, err)

	for _, secret := range clientSecrets {
		err = grpcClient.SetSecret(ctx, secret)
		require.NoError(t, err)
	}

	err = grpcClient.DeleteSecret(ctx, secretForDeleting)
	require.NoError(t, err)

	remoteSecrets, err := grpcClient.ListSecrets(ctx)
	require.NoError(t, err)
	require.Len(t, remoteSecrets, 2)

	err = grpcClient.Close()
	require.NoError(t, err)
}

func UserMockRepository(ctrl *gomock.Controller) repository.UserRepositorier {

	passwordHashByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	passwordHash := string(passwordHashByte)

	newUser := &servermodels.User{
		ID:           "1",
		Login:        loginNewUser,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	existingUser := &servermodels.User{
		ID:           "2",
		Login:        loginExistingUser,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}

	userMockRepository := mockserverrepository.NewMockUserRepositorier(ctrl)

	userMockRepository.EXPECT().GetByLogin(gomock.Any(), loginNewUser).Return(nil, nil)
	userMockRepository.EXPECT().Create(gomock.Any(), loginNewUser, gomock.Any()).Return(newUser, nil)
	userMockRepository.EXPECT().GetByLogin(gomock.Any(), loginExistingUser).Return(existingUser, nil).MaxTimes(2)

	return userMockRepository
}

func SecretMockRepository(ctrl *gomock.Controller) repository.SecretRepositorier {

	secretMockRepository := mockserverrepository.NewMockSecretRepositorier(ctrl)

	secretMockRepository.EXPECT().CreateUpdate(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
	secretMockRepository.EXPECT().Delete(gomock.Any(), gomock.Any(), secretForDeleting).Return(nil)
	secretMockRepository.EXPECT().GetByID(gomock.Any(), gomock.Any(), secretForDeleting).Return(nil, nil)

	resList := []*servermodels.Secret{
		{
			ID: secretOnlyCreating,
		},
		{
			ID: secretForUpdating,
		},
	}
	secretMockRepository.EXPECT().GetAll(gomock.Any(), gomock.Any()).Return(resList, nil)

	return secretMockRepository
}
