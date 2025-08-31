package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type PasswordResetRepository struct {
	collection *mongo.Collection
}

func NewPasswordResetRepository(db *mongo.Database) domain.IPasswordResetRepository {
	return &PasswordResetRepository{
		collection: db.Collection("password_resets"),
	}
}

func (r *PasswordResetRepository) Store(ctx context.Context, token *domain.PasswordResetToken) error {
	// Generate new ObjectID if not set
	if token.ID == "" {
		token.ID = primitive.NewObjectID().Hex()
	}
	
	// Set creation time if not set
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}
	
	// Delete any existing unused tokens for this user
	r.collection.DeleteMany(ctx, bson.M{
		"user_id": token.UserID,
		"used":    false,
	})
	
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *PasswordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	// We need to find by comparing the hash, not exact match
	// Since we store bcrypt hash, we need to find all unused tokens and compare
	filter := bson.M{
		"used":       false,
		"expires_at": bson.M{"$gt": time.Now()},
	}
	
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	for cursor.Next(ctx) {
		var token domain.PasswordResetToken
		if err := cursor.Decode(&token); err != nil {
			continue
		}
		
		// Compare the provided token with stored hash
		if err := bcrypt.CompareHashAndPassword([]byte(token.TokenHash), []byte(tokenHash)); err == nil {
			return &token, nil
		}
	}
	
	return nil, domain.ErrInvalidResetToken
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"used": true,
		},
	}
	
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	
	if result.MatchedCount == 0 {
		return domain.ErrInvalidResetToken
	}
	
	return nil
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	filter := bson.M{"expires_at": bson.M{"$lt": now}}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
