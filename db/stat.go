package db

import (
	"context"
	"fmt"

	"github.com/elmasy-com/columbus-sdk/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// GetStat resturns the total number of domains (d), the total number of subdomains (s) and the error (if any).
func GetStat() (d int64, s int64, err error) {

	var dom domain.Domain

	cursor, err := Domains.Find(context.TODO(), bson.D{})
	if err != nil {
		return d, s, fmt.Errorf("find() failed: %w", err)
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {

		err = cursor.Decode(&dom)
		if err != nil {
			return d, s, fmt.Errorf("decode() failed: %w", err)
		}

		d += 1
		s += int64(len(dom.Subs))
	}

	if err := cursor.Err(); err != nil {
		return d, s, fmt.Errorf("cursor failed: %w", err)
	}

	return d, s, nil
}
