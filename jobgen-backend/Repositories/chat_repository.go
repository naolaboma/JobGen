package repositories

import (
	"context"
	"time"

	"jobgen-backend/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type chatRepository struct {
	db *mongo.Database
}

func NewChatRepository(db *mongo.Database) domain.IChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateSession(ctx context.Context, session *domain.ChatSession) error {
	session.ID = primitive.NewObjectID().Hex()
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	
	_, err := r.db.Collection("chat_sessions").InsertOne(ctx, session)
	return err
}

func (r *chatRepository) GetSession(ctx context.Context, sessionID, userID string) (*domain.ChatSession, error) {
	var session domain.ChatSession
	err := r.db.Collection("chat_sessions").FindOne(ctx, bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}).Decode(&session)
	
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *chatRepository) UpdateSession(ctx context.Context, session *domain.ChatSession) error {
	session.UpdatedAt = time.Now()
	_, err := r.db.Collection("chat_sessions").UpdateOne(ctx, 
		bson.M{"_id": session.ID}, 
		bson.M{"$set": session})
	return err
}

func (r *chatRepository) GetUserSessions(ctx context.Context, userID string, limit, offset int) ([]domain.ChatSession, error) {
	opts := options.Find().
		SetSort(bson.M{"updated_at": -1}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	
	cursor, err := r.db.Collection("chat_sessions").Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var sessions []domain.ChatSession
	if err = cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *chatRepository) SaveMessage(ctx context.Context, message *domain.ChatMessage) error {
	message.ID = primitive.NewObjectID().Hex()
	message.Timestamp = time.Now()
	
	_, err := r.db.Collection("chat_messages").InsertOne(ctx, message)
	return err
}

func (r *chatRepository) GetSessionMessages(ctx context.Context, sessionID, userID string, limit int) ([]domain.ChatMessage, error) {
	// Verify the session belongs to the user
	var session domain.ChatSession
	err := r.db.Collection("chat_sessions").FindOne(ctx, bson.M{
		"_id":     sessionID,
		"user_id": userID,
	}).Decode(&session)
	
	if err != nil {
		return nil, err
	}
	
	opts := options.Find().SetSort(bson.M{"timestamp": 1})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	
	cursor, err := r.db.Collection("chat_messages").Find(ctx, bson.M{"session_id": sessionID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var messages []domain.ChatMessage
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *chatRepository) DeleteSession(ctx context.Context, sessionID, userID string) error {
	// Delete session and all its messages in a transaction
	session, err := r.db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Delete the session
		_, err := r.db.Collection("chat_sessions").DeleteOne(sessCtx, bson.M{
			"_id":     sessionID,
			"user_id": userID,
		})
		if err != nil {
			return nil, err
		}
		
		// Delete all messages in the session
		_, err = r.db.Collection("chat_messages").DeleteMany(sessCtx, bson.M{
			"session_id": sessionID,
		})
		return nil, err
	})
	
	return err
}
