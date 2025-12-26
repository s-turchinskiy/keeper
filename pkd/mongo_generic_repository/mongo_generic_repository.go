// Package mongo_generic_repository Generic Repository Pattern
// взято отсюда и доработано https://readmedium.com/implementing-a-generic-repository-pattern-for-efficient-data-access-in-golang-applications-ae022f1d0f84
package mongo_generic_repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/s-turchinskiy/keeper/internal/utils/errorsutils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository[T any] struct {
	collection   *mongo.Collection
	entityName   string
	keyFieldName string
}

func NewRepository[T any](collection *mongo.Collection, entityName, keyFieldName string) *Repository[T] {

	return &Repository[T]{
		collection:   collection,
		entityName:   entityName,
		keyFieldName: keyFieldName,
	}
}

func (r *Repository[T]) Create(ctx context.Context, doc *T) (*T, error) {
	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return doc, nil
}

func (r *Repository[T]) GetAll(ctx context.Context) ([]*T, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %v", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println(errorsutils.WrapError(err))
		}
	}(cursor, ctx)

	var results []*T
	_ = cursor.All(ctx, &results)

	return results, nil
}

func (r *Repository[T]) GetByID(ctx context.Context, id string) (*T, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrIdInvalid
	}

	var result T
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	return &result, err
}

func (r *Repository[T]) GetByKey(ctx context.Context, keyValue string) (*T, error) {

	filter := bson.D{{Key: r.keyFieldName, Value: keyValue}}
	result := r.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return nil, result.Err()
	}

	var doc T
	err := result.Decode(&doc)

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *Repository[T]) Update(ctx context.Context, id string, updateDoc *T) (*T, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrIdInvalid
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": updateDoc}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		return nil, NewEntityNotFoundError(r.entityName, id)
	}

	return updateDoc, nil
}

func (r *Repository[T]) UpdateByKey(ctx context.Context, keyValue string, updateDoc *T) (*T, error) {

	filter := bson.D{{Key: r.keyFieldName, Value: keyValue}}
	result, err := r.collection.ReplaceOne(ctx, filter, updateDoc)

	if err != nil {
		return nil, err
	}

	if result.ModifiedCount != 1 {
		return nil, NewEntityNotFoundError(r.entityName, keyValue)
	}

	return updateDoc, nil

}

func (r *Repository[T]) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrIdInvalid
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil || result.DeletedCount == 0 {
		return NewEntityNotFoundError(r.entityName, id)
	}

	return nil
}

func (r *Repository[T]) DeleteByKey(ctx context.Context, keyValue string) error {
	filter := bson.D{{Key: r.keyFieldName, Value: keyValue}}
	result, err := r.collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if result.DeletedCount == 1 {
		return NewEntityNotFoundError(r.entityName, keyValue)
	}

	return nil
}

func (r *Repository[T]) CountByID(ctx context.Context, id string) (int64, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, ErrIdInvalid
	}
	count, err := r.collection.CountDocuments(ctx, bson.M{"_id": objectID})
	if err != nil {
		return 0, err
	}

	return count, nil
}
