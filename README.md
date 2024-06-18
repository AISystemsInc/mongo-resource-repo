Go Mongo Repo
=============

A light wrapper around the MongoDB Go Driver that makes it trivial to interact with various model structs.

## Installation

```bash
go get -u github.com/AISystemsInc/mongo-resource-repo
```

## Usage

To use the repository to find documents in a MongoDB collection, you can follow this simple example. This example demonstrates how to find all documents in a collection that match a certain filter.

### Example: Finding Documents

Here's a quick example using a hypothetical `Person` model:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// define your model struct
type Person struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

// implement the GetDatabaseName method for your model
func (p *Person) GetDatabaseName() string {
	return "person_db"
}

// implement the GetCollectionName method for your model
func (p *Person) GetCollectionName() string {
	return "person_col"
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

    // create a new repository for your model, you dont need to use a pointer for the model, but it is typical.
    // the ID type is the type of the _id field in your model, in this case primitive.ObjectID
    personRepo := NewRepository[*Person, primitive.ObjectID](client)

    persons, err := personRepo.Find(ctx, bson.M{"name": "John Doe"})
    if err != nil {
        fmt.Printf("Error finding persons: %v", err)
        return
    }

    fmt.Printf("%+v\n", persons)
}
```
