package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

// GetStat returns the total number of domains.
func GetStat() (int64, error) {

	return Domains.CountDocuments(context.TODO(), bson.M{})
}
