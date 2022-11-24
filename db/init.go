package db

import (
	"context"
	crand "crypto/rand"
	"fmt"

	"github.com/sethvargo/go-password/password"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init initalize the keygenerator, creates the indexes in domains and users collection and creates an "admin" user if no admin privileged user is exist.
func Init() error {

	var err error

	// Initialize API key generator
	keyGeneratorInput = &password.GeneratorInput{
		LowerLetters: password.LowerLetters,
		UpperLetters: password.UpperLetters,
		Digits:       password.Digits,
		Symbols:      "~@#%^&*()_+-={}[]:;<>?,./", // Some characters removed from the original, to be easily usable in Bash
		Reader:       crand.Reader,
	}
	keyGenerator, err = password.NewGenerator(keyGeneratorInput)
	if err != nil {
		return fmt.Errorf("failed to create key generator: %w", err)
	}

	// Create a unique compound index for domain+shard in domains.
	// MongoDB will ignore this block if the index already exist.
	_, err = Domains.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "domain", Value: 1}, {Key: "shard", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create unique domain+shard compound index: %s", err)
	}

	// Create a unique compound index for key+name in users.
	// MongoDB will ignore this block if the index already exist.
	_, err = Users.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "key", Value: 1}, {Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create unique key+name compound index: %s", err)
	}

	// Create a new admin user, if no admin privileged user is exist
	n, err := Users.CountDocuments(context.TODO(), bson.M{"admin": true})
	if err != nil {
		return fmt.Errorf("failed to count admin users: %s", err)
	} else if n == 0 {
		fmt.Printf("No admin user found, creating...\n")

		// TODO: Give the key to the user, posibly write to a file and removes it later
		_, err = UserCreate("admin", true)
		if err != nil {
			return fmt.Errorf("failed to create initial admin API key: %s", err)
		}
	}

	return nil
}
