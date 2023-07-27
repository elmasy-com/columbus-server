package db

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ScannerUpdateIndex updates the index for scanner.
// The name is converted to lowercase.
func ScannerUpdateIndex(name string, index int64) error {

	name = strings.ToLower(name)

	_, err := Scanner.UpdateOne(context.TODO(), bson.D{{Key: "name", Value: name}}, bson.D{{Key: "$set", Value: bson.D{{Key: "index", Value: index}}}}, options.Update().SetUpsert(true))

	return err
}

// ScannerGetIndex returns the index for scanner.
// The name is converted to lowercase.
func ScannerGetIndex(name string) (int64, error) {

	name = strings.ToLower(name)

	s := new(ScannerSchema)

	err := Scanner.FindOne(context.TODO(), bson.D{{Key: "name", Value: name}}).Decode(s)

	return s.Index, err
}

// ScannerGetIndexes returns every entry from the "scanner" database.
func ScannerGetIndexes() ([]ScannerSchema, error) {

	cursor, err := Scanner.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"name": 1}))
	if err != nil {
		return nil, fmt.Errorf("failed to find: %w", err)
	}
	defer cursor.Close(context.TODO())

	scs := make([]ScannerSchema, 0)

	for cursor.Next(context.TODO()) {

		sc := new(ScannerSchema)

		err = cursor.Decode(sc)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w", err)
		}

		err = sc.UpdateTotal()
		if err != nil {
			return nil, fmt.Errorf("failed to update total of %s: %w", sc.Name, err)
		}

		scs = append(scs, *sc)
	}

	return scs, cursor.Err()
}
