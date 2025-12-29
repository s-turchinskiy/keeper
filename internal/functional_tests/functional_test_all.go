package functional_tests

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/crypto"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	clientmodels "github.com/s-turchinskiy/keeper/internal/client/models"
	clientrepository "github.com/s-turchinskiy/keeper/internal/client/repository"
	"github.com/s-turchinskiy/keeper/internal/client/service"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type secretsType struct {
	base clientmodels.BaseSecret
	data clientmodels.TextData
}

func FunctionalTestAll(
	t *testing.T,
	usersRepository repository.UserRepositorier,
	secretRepository repository.SecretRepositorier,
	mongoRepository clientrepository.Repositorier,
	debug bool) {

	runGRPCServer(usersRepository, secretRepository)

	ctx := context.Background()
	if !debug {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	grpcClient, err := grpcclient.NewGRPCClient(ctx, "passthrough://bufnet", loginExistingUser, password, grpc.WithContextDialer(bufDialer))
	require.NoError(t, err)

	cryptor := crypto.NewCryptor(password, loginExistingUser)
	srvc := service.NewService(ctx, mongoRepository, grpcClient, service.WithCrypto(cryptor))

	clientSecrets := []secretsType{
		{
			base: clientmodels.BaseSecret{
				Type: clientmodels.SecretTypeText,
				Name: "Secret only creating",
			},

			data: clientmodels.TextData{
				Content: "data",
			},
		},
		{
			base: clientmodels.BaseSecret{
				Type: clientmodels.SecretTypeText,
				Name: "Secret for deleting",
			},

			data: clientmodels.TextData{
				Content: "data",
			},
		},
		{
			base: clientmodels.BaseSecret{
				Type: clientmodels.SecretTypeText,
				Name: "Secret for updating",
			},

			data: clientmodels.TextData{
				Content: "data",
			},
		},
	}

	if debug {
		functionalTestAll(ctx, t, srvc, clientSecrets, debug)
	} else {
		t.Run(fmt.Sprintf("test #%d: %s", 1, "Функциональный тест grpc на все методы"), func(t *testing.T) {
			functionalTestAll(ctx, t, srvc, clientSecrets, debug)
		})
	}

}

func functionalTestAll(ctx context.Context, t *testing.T, srvc *service.Service, clientSecrets []secretsType, debug bool) {

	err := srvc.Register(ctx, loginNewUser, password)
	if !debug {
		require.NoError(t, err)
	}

	err = srvc.Register(ctx, loginExistingUser, password)
	require.Error(t, err)

	err = srvc.Login(ctx, loginExistingUser, password)
	require.NoError(t, err)

	for _, secret := range clientSecrets {
		_, err = srvc.CreateSecret(ctx, secret.base, secret.data)
		require.NoError(t, err)
	}

	scrtForUpdating, err := srvc.ReadSecret(ctx, secretForUpdating)
	require.NoError(t, err)

	scrtForUpdating.Metadata = "updated"
	err = srvc.UpdateSecret(ctx, scrtForUpdating)
	require.NoError(t, err)

	err = srvc.DeleteSecret(ctx, secretForDeleting)
	require.NoError(t, err)

	err = srvc.DeleteSecret(ctx, secretForDeleting)
	require.Error(t, err)

	remoteSecrets, err := srvc.ListLocalSecrets(ctx)
	require.NoError(t, err)
	require.Len(t, remoteSecrets, 2)

	err = srvc.Close(ctx)
	require.NoError(t, err)
}
