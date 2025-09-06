package repositories

import (
	"context"
	"jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCVRepository struct {
	collection *mongo.Collection
}

// NewCVRepository creates a new CV repository with MongoDB and ensures indexes are set.
func NewCVRepository(db *mongo.Database) (domain.CVRepository, error) {
	collection := db.Collection("cvs")

	// Create indexes for performance and searching
	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userid", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
		{
			Keys:    bson.D{{Key: "rawtext", Value: "text"}, {Key: "skills", Value: "text"}},
			Options: options.Index().SetName("TextSearchIndex"),
		},
	})
	if err != nil {
		return nil, err
	}

	return &mongoCVRepository{collection: collection}, nil
}

func (r *mongoCVRepository) Create(cv *domain.CV) error {
	_, err := r.collection.InsertOne(context.Background(), cv)
	return err
}

func (r *mongoCVRepository) GetByID(id string) (*domain.CV, error) {
	var cv domain.CV
	err := r.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&cv)
	return &cv, err
}

func (r *mongoCVRepository) UpdateStatus(id string, status domain.JobStatus, procError ...string) error {
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedat": time.Now().UTC(),
		},
	}
	if len(procError) > 0 && procError[0] != "" {
		update["$set"].(bson.M)["processingerror"] = procError[0]
	}

	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}

func (r *mongoCVRepository) UpdateWithResults(id string, results *domain.CV) error {
	update := bson.M{
		"$set": bson.M{
			"status":         domain.StatusCompleted,
			"rawtext":        results.RawText,
			"profilesummary": results.ProfileSummary,
			"experiences":    results.Experiences,
			"educations":     results.Educations,
			"skills":         results.Skills,
			"suggestions":    results.Suggestions,
			"score":          results.Score,
			"updatedat":      time.Now().UTC(),
		},
	}
	_, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}
