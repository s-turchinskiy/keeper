package mongodb

import (
	"context"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/client/models"
	"github.com/s-turchinskiy/keeper/internal/client/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func (m MongoDB) CreateSecret(ctx context.Context, secret *models.LocalSecret) (string, error) {

	if _, err := m.GetSecret(ctx, secret.Name); err == nil {
		return "", repository.ErrSecretAlreadyExist
	}

	result, err := m.collection.InsertOne(ctx, secret)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert ObjectiveID")
	}
	return oid.Hex(), nil

}

func (m MongoDB) GetSecret(ctx context.Context, secretID string) (secret *models.LocalSecret, err error) {

	filter := bson.D{{Key: "name", Value: secretID}}
	result := m.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	err = result.Decode(&secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (m MongoDB) UpdateSecret(ctx context.Context, secret *models.LocalSecret) (bool, error) {

	secret.LastModified = time.Now()
	filter := bson.D{{Key: "name", Value: secret.Name}}
	result, err := m.collection.ReplaceOne(ctx, filter, secret)

	if err != nil {
		return false, err
	}

	return result.ModifiedCount == 1, nil

}

func (m MongoDB) DeleteSecret(ctx context.Context, secretID string) (bool, error) {
	filter := bson.D{{Key: "name", Value: secretID}}
	result, err := m.collection.DeleteOne(ctx, filter)

	if err != nil {
		return false, err
	}

	return result.DeletedCount == 1, nil
}

func (m MongoDB) ListSecrets(ctx context.Context) (secrets []*models.LocalSecret, err error) {

	cursor, err := m.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var secret *models.LocalSecret
		if err := cursor.Decode(&secret); err != nil {
			fmt.Println("Error decoding document:", err)
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	return secrets, nil

}
