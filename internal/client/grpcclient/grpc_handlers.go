package grpcclient

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (c *GRPCClient) Connect(ctx context.Context) error {
	conn, err := grpc.NewClient(c.serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.authClient = proto.NewAuthServiceClient(conn)
	c.secretClient = proto.NewSecretServiceClient(conn)

	return nil
}

func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *GRPCClient) GetConnectionNumber(ctx context.Context) (uint64, error) {

	resp, err := c.authClient.GetConnectionNumber(ctx, &proto.GetConnectionNumberRequest{})
	if err != nil {
		return 0, err
	}

	return resp.ConnectionNumber, nil
}

func (c *GRPCClient) Login(ctx context.Context, login, password string) error {

	c.token = ""
	c.login = login
	c.password = password

	req := &proto.LoginRequest{
		Login:    c.login,
		Password: c.password,
	}

	resp, err := c.authClient.Login(c.withConnNumber(ctx), req)
	if err != nil {
		return err
	}

	c.token = resp.GetToken()

	err = c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		c.stream, err = c.secretClient.GetUpdatedSecrets(authCtx, &proto.GetUpdatedSecretsRequest{})
		return errorsutils.WrapError(err)
	})

	if err != nil {
		return errorsutils.WrapError(err)
	}

	return nil
}

func (c *GRPCClient) Register(ctx context.Context, login, password string) (string, error) {

	c.token = ""
	c.login = login
	c.password = password

	req := &proto.RegisterRequest{
		Login:    c.login,
		Password: c.password,
	}

	resp, err := c.authClient.Register(c.withConnNumber(ctx), req)
	if err != nil {
		return "", err
	}
	return resp.GetUserId(), nil
}

func (c *GRPCClient) SetSecret(ctx context.Context, secret *models.RemoteSecret) error {
	req := &proto.SetSecretRequest{
		Secret: models.ConvertRemoteSecretToProtoSecret(secret),
	}

	return c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		_, err := c.secretClient.SetSecret(authCtx, req)
		return err
	})
}

func (c *GRPCClient) GetSecret(ctx context.Context, secretID string) (*models.RemoteSecret, error) {
	var resp *proto.GetSecretResponse

	req := &proto.GetSecretRequest{
		SecretId: secretID,
	}

	err := c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		var err error
		resp, err = c.secretClient.GetSecret(authCtx, req)
		return err
	})
	if err != nil {
		return nil, err
	}

	secretResp := resp.GetSecret()

	return models.ConvertProtoSecretToRemoteSecret(secretResp), nil
}

func (c *GRPCClient) DeleteSecret(ctx context.Context, secretID string) error {

	req := &proto.DeleteSecretRequest{
		SecretId: secretID,
	}

	return c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		_, err := c.secretClient.DeleteSecret(authCtx, req)
		return err
	})
}

func (c *GRPCClient) UpdateSecret(ctx context.Context, secret *models.RemoteSecret) error {
	req := &proto.UpdateSecretRequest{
		Secret: models.ConvertRemoteSecretToProtoSecret(secret),
	}

	return c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		_, err := c.secretClient.UpdateSecret(authCtx, req)
		return err
	})
}

func (c *GRPCClient) ListSecrets(ctx context.Context) ([]*models.RemoteSecret, error) {
	var resp *proto.ListSecretsResponse
	err := c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		var err error
		resp, err = c.secretClient.ListSecrets(authCtx, &proto.ListSecretsRequest{})
		return err
	})
	if err != nil {
		return nil, err
	}

	secretsResp := resp.GetSecrets()
	secrets := make([]*models.RemoteSecret, len(secretsResp))
	for i, secretResp := range secretsResp {
		secrets[i] = &models.RemoteSecret{
			Name:         secretResp.GetId(),
			LastModified: secretResp.GetLastModified().AsTime(),
			Hash:         secretResp.GetHash(),
		}
	}

	return secrets, nil
}

func (c *GRPCClient) SyncSecretsFromClient(ctx context.Context, secrets []*models.RemoteSecret) error {

	protoSecrets := make([]*proto.Secret, len(secrets))
	for _, remoteSecret := range secrets {
		protoSecrets = append(protoSecrets, models.ConvertRemoteSecretToProtoSecret(remoteSecret))
	}

	req := &proto.SyncSecretsFromClientRequest{
		Secrets: protoSecrets,
	}

	return c.withAuthRetry(c.withConnNumber(ctx), func(authCtx context.Context) error {
		_, err := c.secretClient.SyncSecretsFromClient(authCtx, req)
		return err
	})
}

func (c *GRPCClient) GetStream() grpc.ServerStreamingClient[proto.GetUpdatedSecretsResponse] {
	return c.stream
}
