package usecases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"jobgen-backend/Domain"
)

type chatUsecase struct {
	chatRepo  domain.IChatRepository
	aiService domain.IAIService
}

func NewChatUsecase(chatRepo domain.IChatRepository, aiService domain.IAIService) domain.IChatUsecase {
	return &chatUsecase{
		chatRepo:  chatRepo,
		aiService: aiService,
	}
}

func (u *chatUsecase) SendMessage(ctx context.Context, req *domain.ChatRequest) (*domain.ChatResponse, error) {
	// Get or create session
	var session *domain.ChatSession
	var err error
	
	if req.SessionID != "" {
		session, err = u.chatRepo.GetSession(ctx, req.SessionID, req.UserID)
		if err != nil {
			return nil, fmt.Errorf("session not found: %v", err)
		}
	} else {
		// Create new session
		session = &domain.ChatSession{
			UserID:    req.UserID,
			Title:     truncateString(req.Message, 50), // Use first message as temporary title
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = u.chatRepo.CreateSession(ctx, session)
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %v", err)
		}
	}
	
	// Get recent message history (last 10 messages)
	history, err := u.chatRepo.GetSessionMessages(ctx, session.ID, req.UserID, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get message history: %v", err)
	}
	
	// Save user message
	userMessage := &domain.ChatMessage{
		SessionID: session.ID,
		Role:      "user",
		Content:   req.Message,
	}
	err = u.chatRepo.SaveMessage(ctx, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %v", err)
	}
	
	// Generate AI response
	var aiResponse string
	if strings.Contains(strings.ToLower(req.Message), "analyze my cv") {
		// Special handling for CV analysis
		aiResponse, err = u.aiService.AnalyzeCV(ctx, extractCVText(req.Message))
	} else if strings.Contains(strings.ToLower(req.Message), "find job") || 
		strings.Contains(strings.ToLower(req.Message), "job search") {
		// Special handling for job search
		aiResponse, err = u.aiService.FindJobs(ctx, "", req.Message)
	} else {
		// General conversation
		aiResponse, err = u.aiService.GenerateResponse(ctx, req.Message, history)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %v", err)
	}
	
	// Save AI response
	aiMessage := &domain.ChatMessage{
		SessionID: session.ID,
		Role:      "assistant",
		Content:   aiResponse,
	}
	err = u.chatRepo.SaveMessage(ctx, aiMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to save AI message: %v", err)
	}
	
	// Update session
	session.MessageCount += 2 // User + AI messages
	session.UpdatedAt = time.Now()
	
	// If this is a new session and we have an AI response, generate a better title
	if session.Title == truncateString(req.Message, 50) {
		titlePrompt := fmt.Sprintf("Generate a short title (max 5 words) for a conversation that started with: %s", req.Message)
		title, err := u.aiService.GenerateResponse(ctx, titlePrompt, nil)
		if err == nil {
			session.Title = truncateString(title, 30)
		}
	}
	
	err = u.chatRepo.UpdateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %v", err)
	}
	
	// Get updated message history
	updatedHistory, err := u.chatRepo.GetSessionMessages(ctx, session.ID, req.UserID, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated message history: %v", err)
	}
	
	return &domain.ChatResponse{
		SessionID: session.ID,
		Message:   aiResponse,
		History:   updatedHistory,
	}, nil
}

func (u *chatUsecase) GetSessionHistory(ctx context.Context, sessionID, userID string) ([]domain.ChatMessage, error) {
	// Verify the session belongs to the user
	_, err := u.chatRepo.GetSession(ctx, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %v", err)
	}
	
	return u.chatRepo.GetSessionMessages(ctx, sessionID, userID, 0)
}

func (u *chatUsecase) GetUserSessions(ctx context.Context, userID string, limit, offset int) ([]domain.ChatSession, error) {
	return u.chatRepo.GetUserSessions(ctx, userID, limit, offset)
}

func (u *chatUsecase) DeleteSession(ctx context.Context, sessionID, userID string) error {
	return u.chatRepo.DeleteSession(ctx, sessionID, userID)
}

// Helper function to truncate strings
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// Helper function to extract CV text from a message
func extractCVText(message string) string {
	// This is a simple implementation - in a real app, you might have a separate CV upload feature
	// For now, we'll just return the message itself, assuming it contains CV text
	return message
}
