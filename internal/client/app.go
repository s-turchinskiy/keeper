package client

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/client/cmds"
	"github.com/s-turchinskiy/keeper/internal/client/crypto"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	"github.com/s-turchinskiy/keeper/internal/client/repository/mongodb"
	"github.com/s-turchinskiy/keeper/internal/client/service"
	"log"
	"time"

	"github.com/s-turchinskiy/keeper/internal/client/config"
)

type CmdEr interface {
	Run() error
}

type App struct {
	cfg     *config.Config
	cmd     CmdEr
	Service service.Servicer
}

func NewApp() (*App, error) {

	cfg, err := config.LoadCfg(config.WithDB())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	cryptor := crypto.NewCryptor(cfg.Password, cfg.Login)

	repository, err := mongodb.NewMongoDBStorage(ctx, cfg.DBURL)
	if err != nil {
		return nil, err
	}

	serverPassword := cryptor.GenerateServerPassword()
	grpcClient, err := grpcclient.NewGRPCClient(ctx, cfg.ServerAddress, cfg.Login, serverPassword)
	if err != nil {
		return nil, err
	}
	srvc := service.NewService(ctx, repository, grpcClient, service.WithCrypto(cryptor))
	app := &App{
		cfg:     cfg,
		cmd:     cmds.New(srvc),
		Service: srvc,
	}

	return app, nil
}

func (a *App) Run() error {
	return a.cmd.Run()
}

func (a *App) Close() {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if a.Service != nil {
		if err := a.Service.Close(ctx); err != nil {
			log.Printf("Error closing service: %v", err)
		}
	}
}
