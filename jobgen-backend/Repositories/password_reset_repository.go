package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *PasswordResetRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*domain.PasswordResetToken, error) {
	var pr domain.PasswordResetToken
	err := r.collection.FindOne(ctx, bson.M{"token_hash": tokenHash, "used": false}).Decode(&pr)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &pr, err
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"used": true}})
	return err
}

func (r *PasswordResetRepository) DeleteExpired(ctx context.Context) error {
	now := time.Now()
	_, err := r.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": now}})
	return err
}
