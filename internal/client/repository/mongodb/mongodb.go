package mongodb

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/s-turchinskiy/keeper/pkd/dbparse"
	"github.com/s-turchinskiy/keeper/pkd/mongo_generic_repository"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	collectionName = "secrets"
	entityName     = "secret"
	keyName        = "name"
)

type MongoDB struct {
	mongo_generic_repository.Repository[models.LocalSecret]
	client *mongo.Client
}

func NewMongoDBStorage(ctx context.Context, mongoDBURL string) (db *MongoDB, err error) {

	clientOptions := options.Client().ApplyURI(mongoDBURL)
	parsedStr, err := dbparse.ParsedConnectionString(mongoDBURL)
	if err != nil {
		return nil, err
	}

	clientOptions.SetAuth(options.Credential{
		AuthSource: parsedStr.DBName,
		Username:   parsedStr.Login,
		Password:   parsedStr.Password,
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongoDB due to error: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB due to error: %v", err)
	}

	return &MongoDB{
		client: client,
		Repository: *mongo_generic_repository.NewRepository[models.LocalSecret](
			client.Database(parsedStr.DBName).Collection(collectionName),
			entityName,
			keyName,
		),
	}, nil
}

func (m MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
