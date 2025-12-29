package main

import (
	"context"
	"github.com/s-turchinskiy/keeper/internal/server/config"
	"log"

	"github.com/joho/godotenv"
	"github.com/s-turchinskiy/keeper/internal/server"
)

func main() {

	err := godotenv.Load("./.env")
	if err != nil {
		_ = godotenv.Load("./cmd/server/.env")
	}

	cfg, err := config.LoadCfg(config.WithRedis(), config.WithJWT())
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}

	ctx := context.Background()
	app, err := server.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to create app:", err)
	}
	defer app.Stop(ctx)

	if err := app.Run(); err != nil {
		log.Fatal("Failed to run app:", err)
	}
}
