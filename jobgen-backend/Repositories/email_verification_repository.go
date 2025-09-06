package repositories

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmailVerificationRepository struct {
	collection *mongo.Collection
}

func NewEmailVerificationRepository(db *mongo.Database) domain.IEmailVerificationRepository {
	repo := &EmailVerificationRepository{
		collection: db.Collection("email_verifications"),
	}
	repo.createIndexes()
	return repo
}

func (e *EmailVerificationRepository) createIndexes() {
	ctx := context.Background()
	// Compound index for email and purpose to quickly find active OTPs
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}, {Key: "purpose", Value: 1}, {Key: "used", Value: 1}, {Key: "expires_at", Value: 1}},
		Options: options.Index(),
	}
	e.collection.Indexes().CreateOne(ctx, indexModel)
}

func (e *EmailVerificationRepository) Store(ctx context.Context, verification *domain.EmailVerification) error {
	if verification.ID == "" {
		// Use timestamp-based string ID instead of ObjectID
		verification.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if verification.CreatedAt.IsZero() {
		verification.CreatedAt = time.Now()
	}

	// Delete any existing non-used verification for this email and purpose first
	if err := e.DeleteExisting(ctx, verification.Email, verification.Purpose); err != nil {
		return err
	}

	_, err := e.collection.InsertOne(ctx, verification)
	return err
}

func (e *EmailVerificationRepository) GetByEmailAndPurpose(ctx context.Context, email string, purpose domain.OTPPurpose) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	filter := bson.M{
		"email":      email,
		"purpose":    purpose,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()}, // Not expired
	}

	err := e.collection.FindOne(ctx, filter).Decode(&ev)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &ev, err
}

func (e *EmailVerificationRepository) GetByOTPAndEmailAndPurpose(ctx context.Context, otp, email string, purpose domain.OTPPurpose) (*domain.EmailVerification, error) {
	var ev domain.EmailVerification
	filter := bson.M{
		"otp":        otp,
		"email":      email,
		"purpose":    purpose,
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	err := e.collection.FindOne(ctx, filter).Decode(&ev)
	fmt.Println("GET by OTP", err)
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

func (e *EmailVerificationRepository) DeleteExisting(ctx context.Context, email string, purpose domain.OTPPurpose) error {
	filter := bson.M{
		"email":   email,
		"purpose": purpose,
		"used":    false,
	}
	_, err := e.collection.DeleteMany(ctx, filter)
	return err
}
