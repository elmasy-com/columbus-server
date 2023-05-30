package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client

	Domains  *mongo.Collection // The main collection to store the entries
	NotFound *mongo.Collection // Store domains that not found by Lookup
	TopList  *mongo.Collection // Store and count successful lookups
)

// Connect connects to the database using the standard Connection URI.
func Connect(uri string) error {

	var err error

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	err = Client.Ping(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	Domains = Client.Database("columbus").Collection("domains")
	NotFound = Client.Database("columbus").Collection("notFound")
	TopList = Client.Database("columbus").Collection("topList")

	return nil
}

// Disconnect gracefully disconnect from the database.
func Disconnect() error {
	return Client.Disconnect(context.Background())
}
