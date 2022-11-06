package db

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Client  *mongo.Client
	Domains *mongo.Collection
	Users   *mongo.Collection
)
