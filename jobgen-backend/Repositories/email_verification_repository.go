package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailVerificationRepository struct {
	collection *mongo.Collection
}

func NewEmailVerificationRepository(db *mongo.Database) domain.IEmailVerificationRepository {
	return &EmailVerificationRepository{
		collection: db.Collection("email_verifications"),
	}
}

func (e *EmailVerificationRepository) Store(ctx context.Context, verification *domain.EmailVerification) error {
	// Generate new ObjectID if not set
	if verification.ID == "" {
		verification.ID = primitive.NewObjectID().Hex()
	}
	
	// Set creation time if not set
	if verification.CreatedAt.IsZero() {
		verification.CreatedAt = time.Now()
	}
	
	// Delete any existing non-used verification for this email first
	e.collection.DeleteMany(ctx, bson.M{
		"email": verification.Email,
		"used":  false,
	})
	
	_, err := e.collection.InsertOne(ctx, verification)
	return err
}

func (e *EmailVerificationRepository) GetByEmail(ctx context.Context, email string) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	filter := bson.M{
		"email": email,
		"used":  false,
		"expires_at": bson.M{"$gt": time.Now()}, // Not expired
	}
	
	err := e.collection.FindOne(ctx, filter).Decode(&ev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ev, err
}

func (e *EmailVerificationRepository) GetByOTP(ctx context.Context, otp string, email string) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	filter := bson.M{
		"otp":   otp,
		"email": email,
		"used":  false,
		"expires_at": bson.M{"$gt": time.Now()}, // Not expired
	}
	
	err := e.collection.FindOne(ctx, filter).Decode(&ev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ev, err
}

func (e *EmailVerificationRepository) MarkUsed(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"used": true,
		},
	}
	
	result, err := e.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return domain.ErrInvalidOTP
	}
	
	return nil
}

func (e *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	filter := bson.M{"expires_at": bson.M{"$lt": now}}
	_, err := e.collection.DeleteMany(ctx, filter)
	return err
}
