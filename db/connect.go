package db

import (
	"context"
	"fmt"
	"os"

	"github.com/elmasy-com/columbus-server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connect connects to the database.
func Connect() error {

	var err error

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %s", err)
	}

	Domains = Client.Database("columbus").Collection("domains")
	Users = Client.Database("columbus").Collection("users")

	return nil
}

// Disconnect gracefully disconnect from the database.
func Disconnect() {

	err := Client.Disconnect(context.TODO())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to disconnect from MongoDB: %s\n", err)
	} else {
		fmt.Printf("MongoDB disconnected!\n")
	}
}
