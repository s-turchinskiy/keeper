package grpcclient

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn         *grpc.ClientConn
	authClient   proto.AuthServiceClient
	secretClient proto.SecretServiceClient

	stream grpc.ServerStreamingClient[proto.GetUpdatedSecretsResponse]

	serverAddress string

	connectionNumber uint64
	login            string
	password         string
	token            string
}

func NewGRPCClient(ctx context.Context, serverAddress string, login string, password string, extraOpts ...grpc.DialOption) (*GRPCClient, error) {

	grpcClient := &GRPCClient{
		serverAddress: serverAddress,
		login:         login,
		password:      password,
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	opts = append(opts, extraOpts...)

	conn, err := grpc.NewClient(serverAddress, opts...)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	grpcClient.conn = conn
	grpcClient.authClient = proto.NewAuthServiceClient(conn)
	grpcClient.secretClient = proto.NewSecretServiceClient(conn)

	grpcClient.connectionNumber, err = grpcClient.GetConnectionNumber(ctx)
	if err != nil {
		return nil, errorsutils.WrapError(err)
	}

	return grpcClient, nil
}
