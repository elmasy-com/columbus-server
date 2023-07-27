package db

import (
	"fmt"
	"strings"

	"github.com/elmasy-com/elnet/ctlog"
)

// Schema used in *notFound* collection.
type NotFoundSchema struct {
	Domain string `bson:"domain" json:"domain"`
}

// Schema used in *topList* collection.
type TopListSchema struct {
	Domain string `bson:"domain" json:"domain"`
	Count  int    `bson:"count" json:"count"`
}

// Schema used in Lookup() to ignore the Records field.
type FastDomainSchema struct {
	Domain string `bson:"domain" json:"domain"`
	TLD    string `bson:"tld" json:"tld"`
	Sub    string `bson:"sub" json:"sub"`
}

// Returns the full hostname (eg.: sub.domain.tld).
func (d *FastDomainSchema) String() string {

	if d.Sub == "" {
		return strings.Join([]string{d.Domain, d.TLD}, ".")
	} else {
		return strings.Join([]string{d.Sub, d.Domain, d.TLD}, ".")
	}
}

// Schema used to store a record in DomainSchema
type RecordSchema struct {
	Type  uint16 `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
	Time  int64  `bson:"time" json:"time"`
}

// Schema used in the "domains" collection.
type DomainSchema struct {
	Domain  string         `bson:"domain" json:"domain"`
	TLD     string         `bson:"tld" json:"tld"`
	Sub     string         `bson:"sub" json:"sub"`
	Updated int64          `bson:"updated" json:"updated"`
	Records []RecordSchema `bson:"records,omitempty" json:"records,omitempty"`
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

// Schema used in "scanner" collection
type ScannerSchema struct {
	Name  string `bson:"name" json:"name"`
	Index int64  `bson:"index" json:"index"`
	Total int64  `bson:"-" json:"total"`
}

func (s *ScannerSchema) UpdateTotal() error {

	l := ctlog.LogByName(s.Name)
	if l == nil {
		return fmt.Errorf("invalid name")
	}

	total, err := ctlog.Size(l.URI)
	if err != nil {
		return fmt.Errorf("failed to get size: %w", err)
	}

	s.Total = total

	return nil
}
