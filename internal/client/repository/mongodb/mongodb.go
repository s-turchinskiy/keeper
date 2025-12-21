package mongodb

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/pkd/dbparse"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionName = "secrets"

type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
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
		client:     client,
		collection: client.Database(parsedStr.DBName).Collection(collectionName),
	}, nil
}

func (m MongoDB) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
