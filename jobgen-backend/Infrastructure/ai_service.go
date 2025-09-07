package infrastructure

import (
	"context"
	"fmt"
	"time"

	"jobgen-backend/Domain"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
)

type aiService struct {
	client      *genai.Client
	model       *genai.GenerativeModel
	rateLimiter *rate.Limiter
}

func NewAIService() (domain.IAIService, error) {
	apiKey := Env.GeminiAPIKey
	modelName := Env.GeminiModel
	
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}
	
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	
	model := client.GenerativeModel(modelName)
	
	// Set up rate limiting (3 requests per minute)
	rateLimiter := rate.NewLimiter(rate.Every(time.Minute), 3)
	
	return &aiService{
		client:      client,
		model:       model,
		rateLimiter: rateLimiter,
	}, nil
}

func (s *aiService) GenerateResponse(ctx context.Context, prompt string, history []domain.ChatMessage) (string, error) {
	// Apply rate limiting
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %v", err)
	}
	
	// Start a chat session
	cs := s.model.StartChat()
	
	// Convert history to Gemini's format
	var historyParts []*genai.Content
	for _, msg := range history {
		role := "user"
		if msg.Role == "assistant" {
			role = "model"
		}
		
		historyParts = append(historyParts, &genai.Content{
			Parts: []genai.Part{genai.Text(msg.Content)},
			Role:  role,
		})
	}
	
	cs.History = historyParts
	
	// Generate response
	resp, err := cs.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	// Extract response text
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
	}
	
	return "I'm sorry, I didn't get a response. Please try again.", nil
}

func (s *aiService) AnalyzeCV(ctx context.Context, cvText string) (string, error) {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %v", err)
	}
	
	prompt := fmt.Sprintf(`You are JobGen, an AI career assistant specializing in helping African professionals find remote tech jobs. 
	Analyze the following CV text and provide specific, actionable suggestions for improvement. Focus on:
	1. Adding quantifiable metrics to achievements
	2. Using strong action verbs
	3. Including relevant technical keywords
	4. Improving structure for Applicant Tracking Systems (ATS)
	5. Tailoring for international remote tech jobs
	
	CV Text: %s
	
	Provide your response in a helpful, professional tone with clear bullet points.`, cvText)
	
	result, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]), nil
	}
	
	return "I couldn't analyze your CV at this time. Please try again later.", nil
}

func (s *aiService) FindJobs(ctx context.Context, userProfile, query string) (string, error) {
	if err := s.rateLimiter.Wait(ctx); err != nil {
		return "", fmt.Errorf("rate limit exceeded: %v", err)
	}
	
	prompt := fmt.Sprintf(`You are JobGen, an AI career assistant specializing in helping African professionals find remote tech jobs. 
	Based on the user profile and query below, provide personalized job search advice and suggestions for remote tech jobs that might be a good fit.
	
	User Profile: %s
	User Query: %s
	
	Provide your response with specific advice, potential job roles to explore, and tips for applying to remote positions. 
	Focus on opportunities that might be open to African candidates.`, userProfile, query)
	
	result, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}
	
	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		response := fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0])
		
		// Add a disclaimer since we're not connecting to real job APIs yet
		disclaimer := "\n\n*Note: This is AI-generated advice based on your profile. For actual job listings, we recommend checking dedicated job platforms like LinkedIn, RemoteOK, and Indeed.*"
		return response + disclaimer, nil
	}
	
	return "I couldn't search for jobs at this time. Please try again later.", nil
}
