package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

type DataSource struct {
	Client            *mongo.Client
	AccountCollection *mongo.Collection
}

func GetDataSource(db string, collection string, ctx context.Context) *DataSource {

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("database_uri")))
	if err != nil {
		log.Fatal(err)
	}

	var dataSource DataSource
	dataSource.Client = client
	dataSource.AccountCollection = dataSource.Client.Database(db).Collection(collection)

	err = dataSource.Client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &dataSource
}
