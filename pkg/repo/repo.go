package repo

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

// Repository is a generic repository for a model.
type Repository[T Model] struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
}

// NewRepository creates a new repository for a model.
// The model must implement the Model interface.
// e.g. usersRepo := NewRepository[*User](client)
func NewRepository[T Model](client *mongo.Client) *Repository[T] {
	var v T
	return &Repository[T]{
		client:         client,
		databaseName:   v.GetDatabaseName(),
		collectionName: v.GetCollectionName(),
	}
}

func (r *Repository[T]) FindOne(
	ctx context.Context,
	filter any,
	opts ...*options.FindOneOptions,
) (T, error) {
	result := r.client.Database(r.databaseName).Collection(r.collectionName).FindOne(
		ctx,
		filter,
		opts...,
	)

	var value T

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return value, ErrNotFound
	}

	err := result.Decode(&value)
	if err != nil {
		return value, fmt.Errorf("failed to decode result: %w", err)
	}

	return value, nil
}

func (r *Repository[T]) Find(
	ctx context.Context,
	filter any,
	opts ...*options.FindOptions,
) ([]T, error) {
	cursor, err := r.client.Database(r.databaseName).Collection(r.collectionName).Find(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %w", err)
	}

	var values []T
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, fmt.Errorf("failed to decode results: %w", err)
	}

	return values, nil
}

func (r *Repository[T]) InsertOne(
	ctx context.Context,
	document T,
	opts ...*options.InsertOneOptions,
) (*primitive.ObjectID, error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).InsertOne(
		ctx,
		document,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	var insertedID primitive.ObjectID

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		insertedID = oid
	} else {
		return nil, fmt.Errorf("failed to convert inserted ID to primitive.ObjectID")
	}

	return &insertedID, nil
}

func (r *Repository[T]) InsertMany(
	ctx context.Context,
	documents []T,
	opts ...*options.InsertManyOptions,
) error {
	var interfaceSlice = make([]any, len(documents))
	for i, d := range documents {
		interfaceSlice[i] = d
	}

	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).InsertMany(
		ctx,
		interfaceSlice,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to insert documents: %w", err)
	}

	return nil
}

func (r *Repository[T]) UpdateByID(
	ctx context.Context,
	id primitive.ObjectID,
	update any,
	opts ...*options.UpdateOptions,
) error {
	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateByID(
		ctx,
		id,
		update,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

func (r *Repository[T]) UpdateOne(
	ctx context.Context,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) error {
	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateOne(
		ctx,
		filter,
		update,
		opts...,
	)

	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

func (r *Repository[T]) UpdateMany(
	ctx context.Context,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) error {
	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateMany(
		ctx,
		filter,
		update,
		opts...,
	)

	if err != nil {
		return fmt.Errorf("failed to update documents: %w", err)
	}

	return nil
}

func (r *Repository[T]) DeleteOne(
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) error {
	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).DeleteOne(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	return nil
}

func (r *Repository[T]) DeleteMany(
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) error {
	_, err := r.client.Database(r.databaseName).Collection(r.collectionName).DeleteMany(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}

	return nil
}
