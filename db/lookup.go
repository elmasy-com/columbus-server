package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/elmasy-com/columbus-sdk/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// Lookup query the DB and returns a list subdomains.
func Lookup(d string) ([]string, error) {

	d = strings.ToLower(d)

	// Use Find() to find every shard of the domain

	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": d})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var subs []string

	for cursor.Next(context.TODO()) {

		var r domain.Domain

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		subs = append(subs, r.Subs...)
	}

	if err := cursor.Err(); err != nil {
		return subs, fmt.Errorf("cursor failed: %w", err)
	}

	return subs, nil
}
