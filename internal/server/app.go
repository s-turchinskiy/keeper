package server

import (
	"context"
	"github.com/s-turchinskiy/keeper/cmd/server/config"
	grpc2 "github.com/s-turchinskiy/keeper/internal/server/grpc"
	"github.com/s-turchinskiy/keeper/internal/server/repository/postgres"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"log"
	"net"

	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/s-turchinskiy/keeper/internal/server/token"
	"google.golang.org/grpc"
)

type App struct {
	grpcServer *grpc.Server
	grpcAddr   string
	db         repository.DBer
}

func NewApp(ctx context.Context, cfg *config.Config) (*App, error) {

	db, err := postgres.NewPostgresStorage(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	jwtManager := token.NewJWTManager(cfg.JWTSecret, cfg.JWTDuration)
	svc := service.NewService(jwtManager, postgres.NewUserRepository(db), postgres.NewSecretRepository(db))
	grpcServer := grpc2.NewGrpcServer(svc)

	return &App{
		grpcServer: grpcServer,
		db:         db,
		grpcAddr:   cfg.GrpcAddr,
	}, nil
}

func (a *App) Run() error {
	lis, err := net.Listen("tcp", a.grpcAddr)
	if err != nil {
		return err
	}

	log.Printf("Server starting on %s", a.grpcAddr)
	return a.grpcServer.Serve(lis)
}

func (a *App) Stop(ctx context.Context) {
	a.grpcServer.GracefulStop()
	a.db.Close(ctx)
}
