package db

import (
	"context"
	"fmt"

	sdkdomain "github.com/elmasy-com/columbus-sdk/domain"
	"github.com/elmasy-com/elnet/domain"
	"go.mongodb.org/mongo-driver/bson"
)

// Update updates the domains if new publixsuffix is added.
// used at the begining to update the domains if new publicsuffix rule is added to the list.
// If the stored domain is not valid after the new rules, create a list of full names, rewrite it with Write() where the new rules will apply, and removes it.
func Update() error {

	cursor, err := Domains.Find(context.TODO(), bson.D{})
	if err != nil {
		return fmt.Errorf("find({}) failed: %s", err)
	}

	var d sdkdomain.Domain

	for cursor.Next(context.TODO()) {

		err = cursor.Decode(&d)
		if err != nil {
			return fmt.Errorf("failed to decode: %s", err)
		}

		dom := domain.GetDomain(d.Domain)
		if dom == d.Domain {
			// Everything is OK.
			continue
		}

		fmt.Printf("%s is not a valid domain, new domain: %s, resolving...\n", d.Domain, dom)

		l := d.GetFull()

		for i := range l {
			if err = Insert(l[i]); err != nil {
				return fmt.Errorf("failed to write %s: %s", l[i], err)
			}
		}

		_, err = Domains.DeleteOne(context.TODO(), bson.D{{Key: "domain", Value: d.Domain}, {Key: "shard", Value: d.Shard}})
		if err != nil {
			return fmt.Errorf("failed to remove %s/%d: %s", d.Domain, d.Shard, err)
		}

	}

	return nil
}
