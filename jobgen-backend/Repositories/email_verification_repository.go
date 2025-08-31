package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	_, err := e.collection.InsertOne(ctx, verification)
	return err
}

func (e *EmailVerificationRepository) GetByEmail(ctx context.Context, email string) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	err := e.collection.FindOne(ctx, bson.M{"email": email, "used": false}).Decode(&ev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ev, err
}

func (e *EmailVerificationRepository) GetByOTP(ctx context.Context, otp string, email string) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	err := e.collection.FindOne(ctx, bson.M{"otp": otp, "email": email, "used": false}).Decode(&ev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ev, err
}

func (e *EmailVerificationRepository) MarkUsed(ctx context.Context, id string) error {
	_, err := e.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"used": true}})
	return err
}

func (e *EmailVerificationRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	_, err := e.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": now}})
	return err
}
