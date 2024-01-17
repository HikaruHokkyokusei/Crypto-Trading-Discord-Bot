package Mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	C *mongo.Collection
}

type UpsertOneResult struct {
	Error       error
	Success     bool
	WasUpdated  bool
	WasInserted bool
}

func (c Collection) UpsertOne(ctx context.Context, filter any, updated any) *UpsertOneResult {
	operation := bson.D{{"$set", updated}}
	opts := options.Update().SetUpsert(true)

	if res, err := c.C.UpdateOne(ctx, filter, operation, opts); err == nil {
		return &UpsertOneResult{
			Success:     true,
			WasInserted: res.UpsertedCount > 0,
			WasUpdated:  res.ModifiedCount > 0,
		}
	} else {
		return &UpsertOneResult{
			Success: false,
			Error:   err,
		}
	}
}
