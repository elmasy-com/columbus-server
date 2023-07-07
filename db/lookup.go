package db

import (
	"context"
	"fmt"
	"time"

	"github.com/elmasy-com/columbus-server/fault"
	"github.com/elmasy-com/elnet/dns"
	"github.com/elmasy-com/slices"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lookup query the DB and returns a list subdomains.
// days specify, that the returned subdomain must had a valid record in the previous n days.
//
// If d has a subdomain, removes it before the query.
// If days is -1, every subdomain returned.
//
// If d is invalid return fault.ErrInvalidDomain.
// If failed to get parts of d (eg.: d is a TLD), returns ault.ErrGetPartsFailed.
func Lookup(d string, days int) ([]string, error) {

	if !dns.IsValid(d) {
		return nil, fault.ErrInvalidDomain
	}

	d = dns.Clean(d)

	p := dns.GetParts(d)
	if p == nil || p.Domain == "" || p.TLD == "" {
		return nil, fault.ErrGetPartsFailed
	}

	var doc primitive.D
	if days > 0 {
		doc = bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}, {Key: "records.time", Value: bson.D{{Key: "$gt", Value: time.Now().AddDate(0, 0, -1*days).Unix()}}}}
	} else {
		doc = bson.D{{Key: "domain", Value: p.Domain}, {Key: "tld", Value: p.TLD}}
	}

	// Use Find() to find every shard of the domain
	cursor, err := Domains.Find(context.TODO(), doc)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var subs []string

	for cursor.Next(context.TODO()) {

		r := new(FastDomainSchema)

		err = cursor.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		subs = append(subs, r.Sub)
	}

	if err := cursor.Err(); err != nil {
		return subs, fmt.Errorf("cursor failed: %w", err)
	}

	return subs, nil
}

// TLD query the DB and returns a list of TLDs for the given domain d.
//
// Domain d must be a valid Second Level Domain (eg.: "example").
//
// NOTE: This function not validate adn Clean() d!
func TLD(d string) ([]string, error) {

	// Use Find() to find every shard of the domain
	cursor, err := Domains.Find(context.TODO(), bson.M{"domain": d})
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var tlds []string

	for cursor.Next(context.TODO()) {

		var r DomainSchema

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		tlds = slices.AppendUnique(tlds, r.TLD)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor failed: %w", err)
	}

	return tlds, nil
}

// Starts query the DB and returns a list of Second Level Domains (eg.: example) that starts with d.
//
// Domain d must be a valid Second Level Domain (eg.: "example").
// This function validate with IsValidSLD() and Clean().
//
// Returns fault.ErrInvalidDomain is d is not a valid Second Level Domain.
func Starts(d string) ([]string, error) {

	if !dns.IsValidSLD(d) {
		return nil, fault.ErrInvalidDomain
	}

	d = dns.Clean(d)

	doc := bson.M{"domain": bson.M{"$regex": fmt.Sprintf("^%s", d)}}

	// Use Find() to find every shard of the domain
	cursor, err := Domains.Find(context.TODO(), doc)
	if err != nil {
		return nil, fmt.Errorf("failed to find: %s", err)
	}
	defer cursor.Close(context.TODO())

	var domains []string

	for cursor.Next(context.TODO()) {

		var r DomainSchema

		err = cursor.Decode(&r)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %s", err)
		}

		domains = slices.AppendUnique(domains, r.Domain)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor failed: %w", err)
	}

	return domains, nil
}
