package mongo

import (
	"context"
	"fmt"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewRepository(client *mongo.Client, db, coll string) *Repository {
	return &Repository{client: client, coll: client.Database(db).Collection(coll)}
}

func (r *Repository) CreateIndexes(ctx context.Context) error { return createIndexes(ctx, r.coll) }

func (r *Repository) Ping(ctx context.Context) error { return r.client.Ping(ctx, readpref.Primary()) }

// her sürücü için tek tek veritabanına gitmek yerine mongo.WriteModel kullanarak tüm listeyi tek bir seferde (batch) MongoDB'ye gönderir. Bu, performans için hayati önem taşır.
func (r *Repository) BulkUpsertLocations(ctx context.Context, items []domain.DriverLocation) (ports.UpsertResult, error) {
	models := make([]mongo.WriteModel, 0, len(items))
	for _, item := range items {
		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"driverId": item.DriverID}).
			SetUpdate(bson.M{"$set": bson.M{
				"driverId":  item.DriverID,
				"location":  bson.M{"type": item.Location.Type, "coordinates": item.Location.Coordinates},
				"updatedAt": item.UpdatedAt,
			}}).
			SetUpsert(true))
	}
	res, err := r.coll.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	if err != nil {
		return ports.UpsertResult{}, fmt.Errorf("bulk write: %w", err)
	}
	return ports.UpsertResult{Upserted: int64(res.UpsertedCount), Updated: int64(res.ModifiedCount)}, nil
}

// en yakın sürücüleri getirir
func (r *Repository) SearchNearest(ctx context.Context, lon, lat float64, radiusM int64, limit int64) ([]domain.NearestDriver, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$geoNear", Value: bson.M{
			"near":          bson.M{"type": "Point", "coordinates": []float64{lon, lat}},
			"distanceField": "distanceM",
			"spherical":     true,
			"maxDistance":   radiusM,
		}}},
		{{Key: "$limit", Value: limit}},
	}
	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var rows []struct {
		DriverID  string  `bson:"driverId"`
		DistanceM float64 `bson:"distanceM"`
		Location  struct {
			Type        string    `bson:"type"`
			Coordinates []float64 `bson:"coordinates"`
		} `bson:"location"`
	}
	if err := cur.All(ctx, &rows); err != nil {
		return nil, err
	}
	out := make([]domain.NearestDriver, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.NearestDriver{
			DriverID:  row.DriverID,
			DistanceM: row.DistanceM,
			Location:  domain.Point{Type: row.Location.Type, Coordinates: row.Location.Coordinates},
		})
	}
	return out, nil
}
