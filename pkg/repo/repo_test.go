package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type FindOneModel struct {
	ID primitive.ObjectID `bson:"_id"`
}

func (f *FindOneModel) GetDatabaseName() string {
	return "find_one_model_db"
}

func (f *FindOneModel) GetCollectionName() string {
	return "find_one_model_col"
}

func TestRepository_FindOne(t *testing.T) {
	type args struct {
		ctx    context.Context
		filter any
		opts   []*options.FindOneOptions
	}
	type testCase[T Model] struct {
		name    string
		r       Repository[T]
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[*FindOneModel]{
		{
			name: "should return a model",
			r: NewRepository[*FindOneModel](nil),
			args: args{
				ctx:    context.Background(),
			}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.FindOne(tt.args.ctx, tt.args.filter, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}
