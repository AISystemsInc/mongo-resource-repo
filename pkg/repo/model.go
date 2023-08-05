package repo

type Model interface {
	GetDatabaseName() string
	GetCollectionName() string
}
