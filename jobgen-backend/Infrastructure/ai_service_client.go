package infrastructure

import (
	"jobgen-backend/Domain"
)

// This is a MOCK implementation of the AIService interface.
// In a real scenario, this would make an HTTP call to the AI service.
type mockAIServiceClient struct{}

func NewAIServiceClient() domain.AIService {
	return &mockAIServiceClient{}
}

func (s *mockAIServiceClient) AnalyzeCV(rawText string) ([]domain.Suggestion, error) {
	suggestions := []domain.Suggestion{
		{ID: "s1", Type: "quantification", Content: "Quantify achievements in your last role.", Applied: false},
		{ID: "s2", Type: "weak_action_verbs", Content: "Replace 'Managed' with a stronger verb like 'Orchestrated'.", Applied: false},
	}
	return suggestions, nil
}
