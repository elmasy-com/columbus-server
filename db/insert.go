package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/elmasy-com/columbus-sdk/fault"
	"github.com/elmasy-com/elnet/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert insert the given domain d to the database.
// Firstly, checks if d is valid. Then split into sub|domain parts.
// Sharding means, if the document is reached the 16MB limit increase the "shard" field by one.
//
// If domain is invalid, returns fault.ErrInvalidDomain.
func Insert(d string) error {

	if !domain.IsValid(d) {
		return fault.ErrInvalidDomain
	}

	d = strings.ToLower(d)

	dom := domain.GetDomain(d)
	if dom == "" {
		return fault.ErrInvalidDomain
	}

	// TODO: What if sub is empty? For now, add the empty string to the array.
	sub := domain.GetSub(d)

	shard := 0

	/*
	 * Always iterate over every shard, because $addToSet iterate over every shard's every subs and append it only if not subdomain exist.
	 * If sub exist, do nothing.
	 * If sub not exist, add it to the last shard.
	 * This method is slow!
	 */

	for {

		filter := bson.D{{Key: "domain", Value: dom}, {Key: "shard", Value: shard}}
		update := bson.D{{Key: "$addToSet", Value: bson.M{"subs": sub}}}
		opts := options.Update().SetUpsert(true)

		_, err := Domains.UpdateOne(context.TODO(), filter, update, opts)
		if err == nil {
			return nil
		}

		switch {
		case strings.Contains(err.Error(), "Resulting document after update is larger than 16777216"):
			// Increase shard number by one.
			// So, if document with (domain == example.com && shard == 0) is full, update the (document == example.com && shard == 1).
			shard++
		default:
			return fmt.Errorf("failed to update %s: %s", d, err)
		}
	}
}
