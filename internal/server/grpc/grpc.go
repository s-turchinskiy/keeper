package grpc

import (
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/models/proto"
	"google.golang.org/grpc"
)

func NewGrpcServer(service *service.Service) *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor(),
			AuthInterceptor(service),
		),
	)

	proto.RegisterAuthServiceServer(grpcServer, NewAuthHandler(service))
	proto.RegisterSecretServiceServer(grpcServer, NewSecretHandler(service))

	return grpcServer

}
