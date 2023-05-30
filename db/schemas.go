package db

import "strings"

// Schema used in *notFound* collection.
type NotFoundSchema struct {
	Domain string `bson:"domain" json:"domain"`
}

// Schema used in *topList* collection.
type TopListSchema struct {
	Domain string `bson:"domain" json:"domain"`
	Count  int    `bson:"count" json:"count"`
}

// Schema used in the "domains" collection.
type DomainSchema struct {
	Domain string `bson:"domain" json:"domain"`
	TLD    string `bson:"tld" json:"tld"`
	Sub    string `bson:"sub" json:"sub"`
}

// Returns the full hostname (eg.: sub.domain.tld).
func (d *DomainSchema) String() string {

	if d.Sub == "" {
		return strings.Join([]string{d.Domain, d.TLD}, ".")
	} else {
		return strings.Join([]string{d.Sub, d.Domain, d.TLD}, ".")
	}
}

// Returns the domain and tld only (eg.: domain.tld)
func (d *DomainSchema) FullDomain() string {
	return strings.Join([]string{d.Domain, d.TLD}, ".")
}
