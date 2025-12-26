package main

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/crypto"
	"github.com/s-turchinskiy/keeper/internal/client/grpcclient"
	"github.com/s-turchinskiy/keeper/internal/client/repository/mongodb"
	"github.com/s-turchinskiy/keeper/internal/client/service"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/s-turchinskiy/keeper/internal/client/config"
	"github.com/s-turchinskiy/keeper/internal/client/models"
)

func main() {
	//testGRPCClient()
	testService()
}

func testService() {

	_ = godotenv.Load("./cmd/client/.env")
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	/*ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()*/

	grpcClient, err := grpcclient.NewGRPCClient(ctx, cfg.ServerAddress, cfg.Login, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := grpcClient.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	cryptor := crypto.NewCryptor(cfg.Password, cfg.Login)

	repository, err := mongodb.NewMongoDBStorage(ctx, cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	srvc := service.NewService(ctx, cryptor, repository, grpcClient)

	fmt.Println("=== Auth ===")

	if err = srvc.Register(ctx, cfg.Password, cfg.Login); err != nil {
		log.Printf("Register failed: %v", err)
	}

	if err := srvc.Login(ctx, cfg.Password, cfg.Login); err != nil {
		log.Fatal("Login failed:", err)
	}
	fmt.Println("Login successfully")

	fmt.Println("\n=== Secrets ===")

	type secretsType struct {
		base models.BaseSecret
		data models.TextData
	}

	secrets := []secretsType{
		{
			base: models.BaseSecret{
				Type: models.SecretTypeText,
				Name: "Secret only creating",
			},

			data: models.TextData{
				Content: "data",
			},
		},
		{
			base: models.BaseSecret{
				Type: models.SecretTypeText,
				Name: "Secret for deleting",
			},

			data: models.TextData{
				Content: "data",
			},
		},
		{
			base: models.BaseSecret{
				Type: models.SecretTypeText,
				Name: "Secret for updating",
			},

			data: models.TextData{
				Content: "data",
			},
		},
	}
	for _, secret := range secrets {
		if _, err = srvc.CreateSecret(ctx, secret.base, secret.data); err != nil {
			fmt.Printf("CreateUpdate failed secret %s: %v\n", secret.base.Name, err)
		}
	}

	secretForUpdating, err := srvc.ReadSecret(ctx, "Secret for updating")
	if err != nil {
		log.Fatal("ReadSecret failed:", err)
	}

	secretForUpdating.Metadata = "updated"
	err = srvc.UpdateSecret(ctx, secretForUpdating)
	if err != nil {
		log.Fatal("Update failed:", err)
	}

	err = srvc.DeleteSecret(ctx, "Secret for deleting")
	if err != nil {
		log.Fatal("Delete failed:", err)
	}

	err = srvc.DeleteSecret(ctx, "Secret for deleting")
	if err == nil {
		log.Fatal("Delete failed: err must != nil")
	}

	remoteSecrets, err := srvc.ListLocalSecrets(ctx)
	if err != nil {
		log.Fatal("GetAll failed:", err)
	}

	fmt.Printf("Found %d secrets:\n", len(remoteSecrets))
	for i, secret := range remoteSecrets {
		fmt.Printf("%d. %s (last modified: %s) = %v\n", i+1, secret.Name, secret.LastModified, string(secret.Data))
	}

	//TODO: возможно надо подождать пока стрим отработает
	//TODO: но вообще нет, сначала отработывает для 1 и 4, аотом только для 4, потом ни для какого. всегда так и непонятно почему, где-то непонятная ошибка
	time.Sleep(1 * time.Second)

	_ = srvc.Close(ctx)

}

// nolint func testGRPCClient()
func testGRPCClient() {

	_ = godotenv.Load("./cmd/client/.env")
	cfg, err := config.LoadCfg()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	grpcClient, err := grpcclient.NewGRPCClient(ctx, cfg.ServerAddress, cfg.Login, cfg.Password)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := grpcClient.Close(); err != nil {
			log.Printf("Error closing client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := grpcClient.Connect(ctx); err != nil {
		log.Fatal("Connect failed:", err)
	}

	fmt.Println("=== Auth ===")
	if _, err := grpcClient.Register(ctx, cfg.Password, cfg.Login); err != nil {
		log.Printf("Register failed: %v", err)
	}

	if err := grpcClient.Login(ctx, cfg.Password, cfg.Login); err != nil {
		log.Fatal("Login failed:", err)
	}
	fmt.Println("Login successfully")

	fmt.Println("\n=== Secrets ===")

	secrets := []*models.RemoteSecret{
		{
			Name:         "Secret1",
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
		{
			Name:         "Secret2",
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
		{
			Name:         "Secret3",
			LastModified: time.Now(),
			Hash:         "hash",
			Data:         []byte("data"),
		},
	}

	for _, secret := range secrets {
		if err := grpcClient.SetSecret(ctx, secret); err != nil {
			log.Fatal("CreateUpdate failed:", err)
		}
	}

	if err = grpcClient.DeleteSecret(ctx, "Client1"); err != nil {
		log.Fatal("Delete failed:", err)
	}

	remoteSecrets, err := grpcClient.ListSecrets(ctx)
	if err != nil {
		log.Fatal("GetAll failed:", err)
	}

	fmt.Printf("Found %d secrets:\n", len(remoteSecrets))
	for i, secret := range remoteSecrets {
		fmt.Printf("%d. %s (hash: %s)\n", i+1, secret.Name, secret.Hash)
	}

}
