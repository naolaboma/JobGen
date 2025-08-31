package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) domain.IRefreshTokenRepository {
	return &RefreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

// CleanupExpiredTokens implements domain.IRefreshTokenRepository.
func (r *RefreshTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	now := time.Now()
	_, err := r.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lt": now}})
	return err
}

// DeleteAllTokensForUser implements domain.IRefreshTokenRepository.
func (r *RefreshTokenRepository) DeleteAllTokensForUser(ctx context.Context, userID string) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}

// FindByTokenID implements domain.IRefreshTokenRepository.
func (r *RefreshTokenRepository) FindByTokenID(ctx context.Context, tokenID string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	err := r.collection.FindOne(ctx, bson.M{"_id": tokenID}).Decode(&token)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &token, err
}

// RevokeToken implements domain.IRefreshTokenRepository.
func (r *RefreshTokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"token": tokenString})
	return err
}

// StoreToken implements domain.IRefreshTokenRepository.
func (r *RefreshTokenRepository) StoreToken(ctx context.Context, token *domain.RefreshToken) error {
	_, err := r.collection.InsertOne(ctx, token)
	return err
}
