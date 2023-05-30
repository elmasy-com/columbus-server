package db

import (
	"context"

	"github.com/elmasy-com/columbus-server/fault"
	"github.com/elmasy-com/elnet/domain"
	"github.com/elmasy-com/elnet/valid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert inserts the given domain d to the *domains* database.
// Checks if d is valid, do a Clean() and then splits into sub|domain|tld parts.
//
// Returns true if d is new and inserted into the database.
// If domain is invalid, returns fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns ault.ErrGetPartsFailed.
func Insert(d string) (bool, error) {

	if !valid.Domain(d) {
		return false, fault.ErrInvalidDomain
	}

	d = domain.Clean(d)

	p := domain.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return false, fault.ErrGetPartsFailed
	}

	doc := bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "sub", Value: p.Sub}}

	// UpdateOne will insert the document with $setOnInsert + upsert or do nothing
	res, err := Domains.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc}, options.Update().SetUpsert(true))

	return res.UpsertedCount != 0, err
}

// InsertNotFound inserts the given domain d to the *notFound* database.
// Checks if d is valid, do a Clean() and removes the subdomain from d.
//
// Returns true if d is new and inserted into the database.
// If domain is invalid or failed to remove the subdomain, returns fault.ErrInvalidDomain.
func InsertNotFound(d string) (bool, error) {

	if !valid.Domain(d) {
		return false, fault.ErrInvalidDomain
	}

	d = domain.Clean(d)

	v := domain.GetDomain(d)
	if v == "" {
		return false, fault.ErrInvalidDomain
	}

	doc := bson.M{"domain": v}

	// UpdateOne will insert the document with $setOnInsert + upsert or do nothing
	res, err := NotFound.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc}, options.Update().SetUpsert(true))

	return res.UpsertedCount != 0, err
}

// InsertTopList inserts the given domain d to the *topList* database or increase the counter if exists.
// Checks if d is valid, do a Clean() and removes the subdomain from d.
//
// Returns true if d is new and inserted into the database.
// If domain is invalid or failed to remove the subdomain, returns fault.ErrInvalidDomain.
func InsertTopList(d string) (bool, error) {

	if !valid.Domain(d) {
		return false, fault.ErrInvalidDomain
	}

	d = domain.Clean(d)

	v := domain.GetDomain(d)
	if v == "" {
		return false, fault.ErrInvalidDomain
	}

	doc := bson.M{"domain": v}

	// UpdateOne will insert the document with $setOnInsert + $inc + upsert or do nothing
	res, err := TopList.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc, "$inc": bson.M{"count": 1}}, options.Update().SetUpsert(true))

	return res.UpsertedCount != 0, err
}
