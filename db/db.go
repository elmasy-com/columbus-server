package db

import (
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Key   string `bson:"key" json:"key"`
	Name  string `bson:"name" json:"name"`
	Admin bool   `bson:"admin" json:"admin"`
}

type Domain struct {
	Domain string   `bson:"domain" json:"domain"`
	Shard  int      `bson:"shard" json:"shard"`
	Subs   []string `bson:"subs" json:"subs"`
}

var (
	Client  *mongo.Client
	Domains *mongo.Collection
	Users   *mongo.Collection
)

// GetFull resturns the hostnames as a slice of string.
// If Subs is empty returns nil (theoretically impossible).
func (d *Domain) GetFull() []string {

	var list []string

	for i := range d.Subs {
		if d.Subs[i] == "" {
			list = append(list, d.Domain)
		} else {
			list = append(list, strings.Join([]string{d.Subs[i], d.Domain}, "."))
		}
	}

	return list
}
