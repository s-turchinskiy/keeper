package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/s-turchinskiy/keeper/internal/client/config"
	"github.com/s-turchinskiy/keeper/internal/client/repository/mongodb"
	"github.com/s-turchinskiy/keeper/internal/functional_tests"
	serverconfig "github.com/s-turchinskiy/keeper/internal/server/config"
	"github.com/s-turchinskiy/keeper/internal/server/repository/postgres"
	"log"
	"testing"
)

func main() {

	ctx := context.Background()
	_ = godotenv.Load("./cmd/server/.env")
	serverCfg, err := serverconfig.LoadCfg()
	if err != nil {
		log.Fatal(err)
	}
	db, err := postgres.NewPostgresStorage(ctx, serverCfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	secretRepository := postgres.NewSecretRepository(db)
	err = secretRepository.TruncateAllTabs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	functional_tests.FunctionalTestGRPC(
		&testing.T{},
		postgres.NewUserRepository(db),
		secretRepository,
		true)

	err = secretRepository.TruncateAllTabs(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load("./cmd/client/.env")
	if err != nil {
		log.Fatal(err)
	}
	clientCfg, err := config.LoadCfg(config.WithDB())
	if err != nil {
		log.Fatal(err)
	}

	mongoRepository, err := mongodb.NewMongoDBStorage(ctx, clientCfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	err = mongoRepository.DeleteAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	functional_tests.FunctionalTestAll(
		&testing.T{},
		postgres.NewUserRepository(db),
		secretRepository,
		mongoRepository,
		true)
}
