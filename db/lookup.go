package db

import (
	"context"
	"fmt"

	"github.com/elmasy-com/columbus-sdk/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// Lookup query the DB and returns a list subdomains.
// If full is true, return the full hostname, not just the subs.
func Lookup(d string, full bool) ([]string, error) {

	// Use Find() to find every shard of the domain

	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": d})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var r domain.Domain
	var subs []string

	for cursor.Next(context.TODO()) {

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		if full {
			subs = append(subs, r.GetFull()...)
		} else {
			subs = append(subs, r.Subs...)
		}
	}

	if err := cursor.Err(); err != nil {
		return subs, fmt.Errorf("cursor failed: %w", err)
	}

	return subs, nil
}
