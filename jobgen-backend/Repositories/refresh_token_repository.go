package repositories

import (
	"context"
	domain "jobgen-backend/Domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) domain.IRefreshTokenRepository {
	repo := &RefreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
	
	// Create index for token_id
	repo.createIndexes()
	
	return repo
}

func (r *RefreshTokenRepository) createIndexes() {
	ctx := context.Background()
	
	// Create index for token_id for faster lookups
	r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "token_id", Value: 1}},
	})
	
	// Create index for user_id for faster cleanup
	r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	
	// Create TTL index for automatic cleanup of expired tokens
	r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "expires_at", Value: 1}},
	})
}

func (r *RefreshTokenRepository) StoreToken(ctx context.Context, token *domain.RefreshToken) error {
	// Generate new ObjectID if not set
	if token.ID == "" {
		token.ID = primitive.NewObjectID().Hex()
	}
	
	// Set timestamps
	now := time.Now()
	if token.CreatedAt.IsZero() {
		token.CreatedAt = now
	}
	if token.UpdatedAt.IsZero() {
		token.UpdatedAt = now
	}
	
	_, err := r.collection.InsertOne(ctx, token)
	return err
}

func (r *RefreshTokenRepository) FindByTokenID(ctx context.Context, tokenID string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	filter := bson.M{
		"token_id": tokenID,
		"revoked_at": bson.M{"$exists": false}, // Not revoked
		"expires_at": bson.M{"$gt": time.Now()}, // Not expired
	}
	
	err := r.collection.FindOne(ctx, filter).Decode(&token)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &token, err
}

func (r *RefreshTokenRepository) RevokeToken(ctx context.Context, tokenString string) error {
	filter := bson.M{"token": tokenString}
	update := bson.M{
		"$set": bson.M{
			"revoked_at":  time.Now(),
			"updated_at": time.Now(),
		},
	}
	
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *RefreshTokenRepository) DeleteAllTokensForUser(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *RefreshTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	now := time.Now()
	filter := bson.M{
		"$or": []bson.M{
			{"expires_at": bson.M{"$lt": now}},
			{"revoked_at": bson.M{"$exists": true}},
		},
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}
