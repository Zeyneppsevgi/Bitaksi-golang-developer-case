package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createIndexes(ctx context.Context, coll *mongo.Collection) error {
	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "location", Value: "2dsphere"}}, Options: options.Index().SetName("location_2dsphere")},
		{Keys: bson.D{{Key: "driverId", Value: 1}}, Options: options.Index().SetName("driverId_unique").SetUnique(true)},
	})
	return err
}
