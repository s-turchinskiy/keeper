package server

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/s-turchinskiy/keeper/cmd/server/config"
	grpc2 "github.com/s-turchinskiy/keeper/internal/server/grpc"
	"github.com/s-turchinskiy/keeper/internal/server/repository"
	"github.com/s-turchinskiy/keeper/internal/server/repository/postgres"
	"github.com/s-turchinskiy/keeper/internal/server/service"
	"github.com/s-turchinskiy/keeper/internal/server/token"
	"google.golang.org/grpc"
	"log"
	"net"
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
	srvc := service.NewService(
		jwtManager,
		postgres.NewUserRepository(db),
		postgres.NewSecretRepository(db),
		service.WithRedis(redisClient(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB), cfg.RedisExpiration),
	)
	grpcServer := grpc2.NewGrpcServer(srvc)

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

func redisClient(addr, password string, db int) *redis.Client {

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

}
