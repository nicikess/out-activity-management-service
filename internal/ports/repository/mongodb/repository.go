package mongodb

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nicikess/out-run-management-service/internal/domain"
	"github.com/nicikess/out-run-management-service/internal/ports/repository"
)

const (
	databaseName   = "run_management"
	collectionName = "runs"
)

type mongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewRepository creates a new MongoDB repository
func NewRepository(ctx context.Context, uri string) (repository.RunRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database(databaseName).Collection(collectionName)

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("user_status_idx"),
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "start_time", Value: -1},
			},
			Options: options.Index().SetName("user_starttime_idx"),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}

	return &mongoRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *mongoRepository) Create(ctx context.Context, run *domain.Run) error {
	_, err := r.collection.InsertOne(ctx, run)
	return err
}

func (r *mongoRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Run, error) {
	var run domain.Run
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&run)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRunNotFound
		}
		return nil, err
	}
	return &run, nil
}

func (r *mongoRepository) GetActiveByUserID(ctx context.Context, userID uuid.UUID) (*domain.Run, error) {
	var run domain.Run
	err := r.collection.FindOne(ctx, bson.M{
		"user_id": userID,
		"status":  domain.RunStatusActive,
	}).Decode(&run)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRunNotFound
		}
		return nil, err
	}
	return &run, nil
}

func (r *mongoRepository) Update(ctx context.Context, run *domain.Run) error {
	result, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": run.ID},
		run,
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrRunNotFound
	}

	return nil
}

func (r *mongoRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Run, error) {
	opts := options.Find().
		SetSort(bson.D{{Key: "start_time", Value: -1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var runs []*domain.Run
	if err = cursor.All(ctx, &runs); err != nil {
		return nil, err
	}

	return runs, nil
}

func (r *mongoRepository) Close(ctx context.Context) error {
	return r.client.Disconnect(ctx)
}
