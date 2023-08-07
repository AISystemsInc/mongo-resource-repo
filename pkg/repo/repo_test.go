package repo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
	"time"
)

type FindOneModel struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (f *FindOneModel) GetDatabaseName() string {
	return "find_one_model_db"
}

func (f *FindOneModel) GetCollectionName() string {
	return "find_one_model_col"
}

func TestRepository_FindOne(t *testing.T) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:57018"))
	if err != nil {
		t.Errorf("error connecting to mongo: %v", err)
		return
	}

	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			t.Errorf("error disconnecting from mongo: %v", err)
		}
	}()

	var empty = &FindOneModel{}
	var findOneID = primitive.NewObjectID()

	type args struct {
		ctx    context.Context
		filter any
		opts   []*options.FindOneOptions
	}
	type testCase[T Model] struct {
		name      string
		r         *Repository[T, primitive.ObjectID]
		args      args
		bootstrap func() error
		tearDown  func() error
		want      T
		wantErr   error
	}
	tests := []testCase[*FindOneModel]{
		{
			name: "should return a model",
			r:    NewRepository[*FindOneModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx: context.TODO(),
				filter: bson.M{
					"_id": findOneID,
				},
			},
			bootstrap: func() error {
				_, err := mongoClient.Database("find_one_model_db").Collection("find_one_model_col").InsertOne(context.TODO(), &FindOneModel{
					ID: findOneID,
				})
				return err
			},
			tearDown: func() error {
				err := mongoClient.Database("find_one_model_db").Collection("find_one_model_col").Drop(context.Background())
				return err
			},
			want: &FindOneModel{
				ID: findOneID,
			},
			wantErr: nil,
		},
		{
			name: "should return an error if the model is not found",
			r:    NewRepository[*FindOneModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx: context.TODO(),
				filter: bson.M{
					"_id": findOneID,
				},
			},
			bootstrap: func() error {
				return nil
			},
			tearDown: func() error {
				err := mongoClient.Database("find_one_model_db").Collection("find_one_model_col").Drop(context.Background())
				return err
			},
			want:    nil,
			wantErr: ErrFindOne,
		},
		{
			name: "should return an error if there was a decode error",
			r:    NewRepository[*FindOneModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx:    context.TODO(),
				filter: bson.M{},
			},
			bootstrap: func() error {
				_, err := mongoClient.Database("find_one_model_db").Collection("find_one_model_col").InsertOne(context.TODO(), bson.M{
					"_id":  "not an object id",
					"name": true,
				})
				return err
			},
			tearDown: func() error {
				err := mongoClient.Database("find_one_model_db").Collection("find_one_model_col").Drop(context.Background())
				return err
			},
			want:    empty,
			wantErr: ErrFindOne,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bootstrap()
			if err != nil {
				t.Errorf("error bootstrapping test: %v", err)
				return
			}

			defer func() {
				err := tt.tearDown()
				if err != nil {
					t.Errorf("error tearing down test: %v", err)
				}
			}()

			got, err := tt.r.FindOne(tt.args.ctx, tt.args.filter, tt.args.opts...)
			if tt.wantErr != nil && (err == nil || !errors.Is(err, tt.wantErr)) {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && err != nil {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type FindModel struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

func (f *FindModel) GetDatabaseName() string {
	return "find_model_db"
}

func (f *FindModel) GetCollectionName() string {
	return "find_model_col"
}

func TestRepository_Find(t *testing.T) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:57018"))
	if err != nil {
		t.Errorf("error connecting to mongo: %v", err)
		return
	}

	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			t.Errorf("error disconnecting from mongo: %v", err)
		}
	}()

	var empty []*FindModel

	var firstID = primitive.NewObjectID()
	var secondID = primitive.NewObjectID()

	type args struct {
		ctx    context.Context
		filter any
		opts   []*options.FindOptions
	}
	type testCase[M Model, I any] struct {
		name      string
		r         *Repository[M, I]
		args      args
		bootstrap func() error
		tearDown  func() error
		want      []M
		wantErr   error
	}
	tests := []testCase[*FindModel, primitive.ObjectID]{
		{
			name: "should return a slice of models",
			r:    NewRepository[*FindModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx:    context.TODO(),
				filter: bson.M{},
			},
			bootstrap: func() error {
				_, err := mongoClient.Database("find_model_db").Collection("find_model_col").InsertMany(context.TODO(), []interface{}{
					&FindModel{
						ID:   firstID,
						Name: "model 1",
					},
					&FindModel{
						ID:   secondID,
						Name: "model 2",
					},
				})
				return err
			},
			tearDown: func() error {
				err := mongoClient.Database("find_model_db").Collection("find_model_col").Drop(context.Background())
				return err
			},
			want: []*FindModel{
				{
					ID:   firstID,
					Name: "model 1",
				},
				{
					ID:   secondID,
					Name: "model 2",
				},
			},
			wantErr: nil,
		},
		{
			name: "should return an empty slice if no models were found",
			r:    NewRepository[*FindModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx:    context.TODO(),
				filter: bson.M{},
			},
			bootstrap: func() error {
				return nil
			},
			tearDown: func() error {
				err := mongoClient.Database("find_model_db").Collection("find_model_col").Drop(context.Background())
				return err
			},
			want:    empty,
			wantErr: nil,
		},
		{
			name: "should return an error if an object failed to decode",
			r:    NewRepository[*FindModel, primitive.ObjectID](mongoClient),
			args: args{
				ctx:    context.TODO(),
				filter: bson.M{},
			},
			bootstrap: func() error {
				var items []interface{}
				items = append(items, bson.M{
					"_id":  firstID,
					"name": "model 1",
				})
				items = append(items, bson.M{
					"_id":  secondID,
					"name": true,
				})

				_, err := mongoClient.Database("find_model_db").Collection("find_model_col").InsertMany(context.TODO(), items)
				return err
			},
			tearDown: func() error {
				err := mongoClient.Database("find_model_db").Collection("find_model_col").Drop(context.Background())
				return err
			},
			want:    empty,
			wantErr: ErrFind,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.bootstrap()
			if err != nil {
				t.Errorf("error bootstrapping test: %v", err)
				return
			}

			defer func() {
				err := tt.tearDown()
				if err != nil {
					t.Errorf("error tearing down test: %v", err)
				}
			}()

			got, err := tt.r.Find(tt.args.ctx, tt.args.filter, tt.args.opts...)
			if tt.wantErr != nil && (err == nil || !errors.Is(err, tt.wantErr)) {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && err != nil {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type FindStreamModel struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (f *FindStreamModel) GetDatabaseName() string {
	return "find_stream_model_db"
}

func (f *FindStreamModel) GetCollectionName() string {
	return "find_stream_model_col"
}

func TestRepository_FindStream(t *testing.T) {
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:57018"))
	if err != nil {
		t.Errorf("error connecting to mongo: %v", err)
		return
	}

	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			t.Errorf("error disconnecting from mongo: %v", err)
		}
	}()

	var insertDocuments = func(num int) ([]FindStreamModel, error) {
		var models []FindStreamModel
		for i := 0; i < num; i++ {
			models = append(models, FindStreamModel{
				ID: primitive.NewObjectID(),
			})
		}

		var documents []interface{}
		for _, model := range models {
			documents = append(documents, model)
		}

		_, err := mongoClient.Database("find_stream_model_db").Collection("find_stream_model_col").InsertMany(context.Background(), documents)
		return models, err
	}

	var tearDown = func() error {
		return mongoClient.Database("find_stream_model_db").Collection("find_stream_model_col").Drop(context.Background())
	}

	t.Run("should return a stream of models", func(t *testing.T) {
		inserted, err := insertDocuments(10)
		if err != nil {
			t.Errorf("error inserting documents: %v", err)
			return
		}

		defer tearDown()

		var repository = NewRepository[*FindStreamModel, primitive.ObjectID](mongoClient)

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		stream, errChan, cancelChan, err := repository.FindStream(ctx, bson.M{})
		if err != nil {
			t.Errorf("error creating stream: %v", err)
			return
		}

		defer func() {
			close(cancelChan)
		}()

		go func() {
			for err := range errChan {
				t.Errorf("error streaming: %v", err)
				return
			}
		}()

		var models []*FindStreamModel
		for model := range stream {
			models = append(models, model)
		}

		if len(models) != len(inserted) {
			t.Errorf("expected %d models but got %d", len(inserted), len(models))
			return
		}

		for _, model := range models {
			var found bool
			for _, insertedModel := range inserted {
				if model.ID == insertedModel.ID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected model with id %s to be found", model.ID)
				return
			}
		}
	})
}
