package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactRepository struct {
	collection *mongo.Collection
}

func NewContactRepository(db *mongo.Database) domain.IContactRepository {
	return &ContactRepository{
		collection: db.Collection("contact_submissions"),
	}
}

func (r *ContactRepository) Create(ctx context.Context, contact *domain.Contact) error {
	contact.ID = primitive.NewObjectID().Hex()
	contact.CreatedAt = time.Now()
	contact.Replied = false // Default to not replied

	_, err := r.collection.InsertOne(ctx, contact)
	return err
}
