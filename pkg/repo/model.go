package repo

// Model is a generic model.
//
// example:
//
//	type User struct {
//		ID       primitive.ObjectID `bson:"_id"`
//		Username string             `bson:"username"`
//		Password string             `bson:"password"`
//	}
//
//	func (u *User) GetDatabaseName() string {
//		return "users_db"
//	}
//
//	func (u *User) GetCollectionName() string {
//		return "users_col"
//	}
type Model interface {
	GetDatabaseName() string
	GetCollectionName() string
}
