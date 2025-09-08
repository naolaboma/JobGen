package domain

import "context"

type IChatUsecase interface {
    SendMessage(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    GetSessionHistory(ctx context.Context, sessionID, userID string) ([]ChatMessage, error)
    GetUserSessions(ctx context.Context, userID string, limit, offset int) ([]ChatSession, error)
    DeleteSession(ctx context.Context, sessionID, userID string) error
}
