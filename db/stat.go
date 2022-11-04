package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// GetStat resturns the total number of domains (d), the total number of subdomains (s) and the error (if any).
func GetStat() (d int64, s int64, err error) {

	var dom Domain

	result, err := Domains.Find(context.TODO(), bson.D{})
	if err != nil {
		return d, s, fmt.Errorf("find() failed: %w", err)
	}
	defer result.Close(context.TODO())

	for result.Next(context.TODO()) {

		err = result.Decode(&dom)
		if err != nil {
			return d, s, fmt.Errorf("decode() failed: %w", err)
		}

		d += 1
		s += int64(len(dom.Subs))
	}

	if err := result.Err(); err != nil {
		return d, s, fmt.Errorf("cursor failed: %w", err)
	}

	return d, s, nil
}
