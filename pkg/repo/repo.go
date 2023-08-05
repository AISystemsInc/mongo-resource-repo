package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrFindOne    = fmt.Errorf("find one error")
	ErrFind       = fmt.Errorf("find error")
	ErrInsertOne  = fmt.Errorf("insert one error")
	ErrInsertMany = fmt.Errorf("insert many error")
	ErrUpdateOne  = fmt.Errorf("update one error")
	ErrUpdateByID = fmt.Errorf("update by ID error")
	ErrUpdateMany = fmt.Errorf("update many error")
	ErrDeleteOne  = fmt.Errorf("delete one error")
	ErrDeleteMany = fmt.Errorf("delete many error")
)

// Repository is a generic repository for a model.
type Repository[M Model, I any] struct {
	client         *mongo.Client
	databaseName   string
	collectionName string
}

// NewRepository creates a new repository for a model.
// The model must implement the Model interface.
// e.g. usersRepo := NewRepository[*User](client)
func NewRepository[M Model, I any](client *mongo.Client) *Repository[M, I] {
	var v M
	return &Repository[M, I]{
		client:         client,
		databaseName:   v.GetDatabaseName(),
		collectionName: v.GetCollectionName(),
	}
}

func (r *Repository[M, I]) FindOne(
	ctx context.Context,
	filter any,
	opts ...*options.FindOneOptions,
) (M, error) {
	result := r.client.Database(r.databaseName).Collection(r.collectionName).FindOne(
		ctx,
		filter,
		opts...,
	)

	var value M

	if result.Err() != nil {
		return value, fmt.Errorf("%w: %w", ErrFindOne, result.Err())
	}

	err := result.Decode(&value)
	if err != nil {
		return value, fmt.Errorf("%w: failed to decode result: %w", ErrFindOne, err)
	}

	return value, nil
}

func (r *Repository[M, I]) Find(
	ctx context.Context,
	filter any,
	opts ...*options.FindOptions,
) ([]M, error) {
	cursor, err := r.client.Database(r.databaseName).Collection(r.collectionName).Find(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFind, err)
	}

	var values []M
	err = cursor.All(ctx, &values)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to decode results: %w", ErrFind, err)
	}

	return values, nil
}

func (r *Repository[M, I]) InsertOne(
	ctx context.Context,
	document M,
	opts ...*options.InsertOneOptions,
) (*primitive.ObjectID, error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).InsertOne(
		ctx,
		document,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInsertOne, err)
	}

	var insertedID primitive.ObjectID

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		insertedID = oid
	} else {
		return nil, fmt.Errorf("%w: failed to convert inserted ID to primitive.ObjectID", ErrInsertOne)
	}

	return &insertedID, nil
}

func (r *Repository[M, I]) InsertMany(
	ctx context.Context,
	documents []M,
	opts ...*options.InsertManyOptions,
) ([]I, error) {
	var interfaceSlice = make([]any, len(documents))
	for i, d := range documents {
		interfaceSlice[i] = d
	}

	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).InsertMany(
		ctx,
		interfaceSlice,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInsertMany, err)
	}

	var insertedIDs = make([]I, len(result.InsertedIDs))
	for i, id := range result.InsertedIDs {
		if oid, ok := id.(I); ok {
			insertedIDs[i] = oid
		} else {
			var insertID I
			return nil, fmt.Errorf("%w: failed to convert inserted ID to type %T", ErrInsertMany, insertID)
		}
	}

	return insertedIDs, nil
}

func (r *Repository[M, I]) UpdateByID(
	ctx context.Context,
	id primitive.ObjectID,
	update any,
	opts ...*options.UpdateOptions,
) (*UpdateResult[I], error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateByID(
		ctx,
		id,
		update,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateByID, err)
	}

	var updatedID I
	if oid, ok := result.UpsertedID.(I); ok {
		updatedID = oid
	} else {
		return nil, fmt.Errorf("%w: failed to convert updated ID to type %T", ErrUpdateByID, updatedID)
	}

	return &UpdateResult[I]{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    updatedID,
	}, nil
}

func (r *Repository[M, I]) UpdateOne(
	ctx context.Context,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) (*UpdateResult[I], error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateOne(
		ctx,
		filter,
		update,
		opts...,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateOne, err)
	}

	var updatedID I
	if oid, ok := result.UpsertedID.(I); ok {
		updatedID = oid
	} else {
		return nil, fmt.Errorf("%w: failed to convert updated ID to type %T", ErrUpdateOne, updatedID)
	}

	return &UpdateResult[I]{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    updatedID,
	}, nil
}

func (r *Repository[M, I]) UpdateMany(
	ctx context.Context,
	filter any,
	update any,
	opts ...*options.UpdateOptions,
) (*UpdateResult[I], error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).UpdateMany(
		ctx,
		filter,
		update,
		opts...,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUpdateMany, err)
	}

	var updatedID I
	if oid, ok := result.UpsertedID.(I); ok {
		updatedID = oid
	} else {
		return nil, fmt.Errorf("%w: failed to convert updated ID to type %T", ErrUpdateMany, updatedID)
	}

	return &UpdateResult[I]{
		MatchedCount:  result.MatchedCount,
		ModifiedCount: result.ModifiedCount,
		UpsertedCount: result.UpsertedCount,
		UpsertedID:    updatedID,
	}, nil
}

func (r *Repository[M, I]) DeleteOne(
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).DeleteOne(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDeleteOne, err)
	}

	return result, nil
}

func (r *Repository[M, I]) DeleteMany(
	ctx context.Context,
	filter any,
	opts ...*options.DeleteOptions,
) (*mongo.DeleteResult, error) {
	result, err := r.client.Database(r.databaseName).Collection(r.collectionName).DeleteMany(
		ctx,
		filter,
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDeleteMany, err)
	}

	return result, nil
}

// UpdateResult is the same as the mongo.UpdateResult type, but with generic support.
type UpdateResult[I any] struct {
	MatchedCount  int64 // The number of documents matched by the filter.
	ModifiedCount int64 // The number of documents modified by the operation.
	UpsertedCount int64 // The number of documents upserted by the operation.
	UpsertedID    I     // The _id field of the upserted document, or nil if no upsert was done.
}
