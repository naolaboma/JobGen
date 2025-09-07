package domain

import (
	"context"
	"time"
)

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	SessionID string    `json:"session_id" bson:"session_id"`
	Role      string    `json:"role" bson:"role"` // "user" or "assistant"
	Content   string    `json:"content" bson:"content"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

// ChatSession represents a conversation session
type ChatSession struct {
	ID           string       `json:"id" bson:"_id,omitempty"`
	UserID       string       `json:"user_id" bson:"user_id"`
	CreatedAt    time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" bson:"updated_at"`
	MessageCount int          `json:"message_count" bson:"message_count"`
	Title        string       `json:"title" bson:"title"` // First message or generated title
}

// ChatRequest represents a request to the AI chatbot
type ChatRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Message   string `json:"message" binding:"required"`
	UserID    string `json:"-"` // Set from auth context
}

// ChatResponse represents a response from the AI chatbot
type ChatResponse struct {
	SessionID string        `json:"session_id"`
	Message   string        `json:"message"`
	History   []ChatMessage `json:"history,omitempty"`
}

// IChatRepository defines the interface for chat storage
type IChatRepository interface {
	CreateSession(ctx context.Context, session *ChatSession) error
	GetSession(ctx context.Context, sessionID, userID string) (*ChatSession, error)
	UpdateSession(ctx context.Context, session *ChatSession) error
	GetUserSessions(ctx context.Context, userID string, limit, offset int) ([]ChatSession, error)
	SaveMessage(ctx context.Context, message *ChatMessage) error
	GetSessionMessages(ctx context.Context, sessionID, userID string, limit int) ([]ChatMessage, error)
	DeleteSession(ctx context.Context, sessionID, userID string) error
}

// IAIService defines the interface for AI interactions
type IAIService interface {
	GenerateResponse(ctx context.Context, prompt string, history []ChatMessage) (string, error)
	AnalyzeCV(ctx context.Context, cvText string) (string, error)
	FindJobs(ctx context.Context, userProfile, query string) (string, error)
}
