package Mongo

import "go.mongodb.org/mongo-driver/mongo"

type Collection struct {
	collection *mongo.Collection
}
