package db

import (
	"context"
	"fmt"

	"github.com/elmasy-com/columbus-server/fault"
	"github.com/elmasy-com/elnet/dns"
	"github.com/elmasy-com/elnet/valid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Insert inserts the given domain d to the *domains* database.
// Checks if d is valid, do a Clean() and then splits into sub|domain|tld parts.
// If a new domain found, do a RecordsUpdate() after the insert.
//
// Returns true if d is new and inserted into the database.
// If domain is invalid, returns fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns ault.ErrGetPartsFailed.
func Insert(d string) (bool, error) {

	if !valid.Domain(d) {
		return false, fault.ErrInvalidDomain
	}

	d = dns.Clean(d)

	p := dns.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return false, fault.ErrGetPartsFailed
	}

	doc := bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "sub", Value: p.Sub}}

	// UpdateOne will insert the document with $setOnInsert + upsert or do nothing
	res, err := Domains.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc}, options.Update().SetUpsert(true))
	if err != nil {
		return false, fmt.Errorf("failed to update: %w", err)
	}

	if res.UpsertedCount != 0 {
		err = RecordsUpdate(&DomainSchema{Domain: p.Domain, TLD: p.TLD, Sub: p.Sub})
		if err != nil {
			return false, fmt.Errorf("failed to update records for %s: %w", d, err)
		}
	}

	return res.UpsertedCount != 0, nil
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

	d = dns.Clean(d)

	v := dns.GetDomain(d)
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

	d = dns.Clean(d)

	v := dns.GetDomain(d)
	if v == "" {
		return false, fault.ErrInvalidDomain
	}

	doc := bson.M{"domain": v}

	// UpdateOne will insert the document with $setOnInsert + $inc + upsert or do nothing
	res, err := TopList.UpdateOne(context.TODO(), doc, bson.M{"$setOnInsert": doc, "$inc": bson.M{"count": 1}}, options.Update().SetUpsert(true))

	return res.UpsertedCount != 0, err
}
