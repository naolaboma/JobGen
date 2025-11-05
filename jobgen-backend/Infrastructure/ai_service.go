package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	domain "jobgen-backend/Domain"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"
)

type aiService struct {
	client      *genai.Client
	model       *genai.GenerativeModel
	rateLimiter *rate.Limiter // may be nil when disabled
}

func applyDefaultGenConfig(m *genai.GenerativeModel) {
	// genai.GenerationConfig is a value type with pointer fields
	t := float32(0.7)
	tp := float32(0.95)
	tk := int32(40)
	max := int32(2048)
	m.GenerationConfig = genai.GenerationConfig{
		Temperature:     &t,
		TopP:            &tp,
		TopK:            &tk,
		MaxOutputTokens: &max,
	}
}

func isModelNotFound(err error) bool {
	if err == nil {
		return false
	}
	es := strings.ToLower(err.Error())
	return strings.Contains(es, "404") ||
		strings.Contains(es, "not found") ||
		strings.Contains(es, "is not supported for")
}

func (s *aiService) tryModels(ctx context.Context, candidates []string, run func(*genai.GenerativeModel) (string, error)) (string, error) {
	for _, name := range candidates {
		mn := normalizeModel(name)
		m := s.client.GenerativeModel(mn)
		applyDefaultGenConfig(m)
		res, err := run(m)
		if err == nil {
			// adopt working model
			s.model = m
			return res, nil
		}
		if !isModelNotFound(err) {
			return "", err
		}
	}
	return "", fmt.Errorf("no available model among candidates")
}

func NewAIService() (domain.IAIService, error) {
	apiKey := Env.GeminiAPIKey
	modelName := normalizeModel(Env.GeminiModel)

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	// Log selected models and key length for diagnostics
	log.Printf("Gemini init: sdkModel=%s restModel=%s apiKeyLen=%d", modelName, normalizeModelV1Beta(Env.GeminiModel), len(apiKey))

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel(modelName)
	applyDefaultGenConfig(model)

	// Set up configurable rate limiting
	var limiter *rate.Limiter
	rpm := Env.GeminiRPM
	if rpm > 0 {
		// average rpm with burst = rpm
		perReq := time.Minute / time.Duration(rpm)
		limiter = rate.NewLimiter(rate.Every(perReq), rpm)
	}

	return &aiService{
		client:      client,
		model:       model,
		rateLimiter: limiter,
	}, nil
}

// Minimal REST schema for v1beta fallback
type glContent struct {
	Role  string `json:"role,omitempty"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type glRequest struct {
	Contents []glContent `json:"contents"`
}

type glResponse struct {
	Candidates []struct {
		Content glContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func restGenerateText(ctx context.Context, apiKey, model, prompt string, history []domain.ChatMessage) (string, error) {
	if strings.TrimSpace(model) == "" {
		model = "gemini-1.5-flash"
	}
	base := normalizeModel(model)
	candidates := []string{base}
	if !strings.HasSuffix(base, "-latest") {
		candidates = append(candidates, base+"-latest")
	}

	var lastErr error
	for _, restModel := range candidates {
		endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", restModel, apiKey)

		var contents []glContent
		for _, msg := range history {
			role := "user"
			if msg.Role == "assistant" {
				role = "model"
			}
			contents = append(contents, glContent{Role: role, Parts: []struct {
				Text string `json:"text"`
			}{{Text: msg.Content}}})
		}
		contents = append(contents, glContent{Role: "user", Parts: []struct {
			Text string `json:"text"`
		}{{Text: prompt}}})

		payload := glRequest{Contents: contents}
		b, _ := json.Marshal(payload)

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(string(b)))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			bodyBytes, _ := io.ReadAll(resp.Body)
			var er glResponse
			_ = json.Unmarshal(bodyBytes, &er)
			msg := strings.TrimSpace(string(bodyBytes))
			if er.Error != nil && er.Error.Message != "" {
				msg = er.Error.Message
			}
			log.Printf("Gemini REST v1beta error %d (%s): %s", resp.StatusCode, restModel, msg)
			// Continue to next candidate if model not found/unsupported
			low := strings.ToLower(msg)
			if resp.StatusCode == 404 || strings.Contains(low, "not found") || strings.Contains(low, "not supported") {
				lastErr = fmt.Errorf("gemini http %d: %s", resp.StatusCode, msg)
				continue
			}
			return "", fmt.Errorf("gemini http %d: %s", resp.StatusCode, msg)
		}

		var gr glResponse
		if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
			lastErr = err
			continue
		}
		if gr.Error != nil {
			log.Printf("Gemini REST v1beta error (%s): %s", restModel, gr.Error.Message)
			lastErr = fmt.Errorf("gemini error: %s", gr.Error.Message)
			continue
		}
		if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
			lastErr = fmt.Errorf("no candidates")
			continue
		}
		return gr.Candidates[0].Content.Parts[0].Text, nil
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no rest model candidate succeeded")
	}
	return "", lastErr
}

func (s *aiService) GenerateResponse(ctx context.Context, prompt string, history []domain.ChatMessage) (string, error) {
	// Apply rate limiting
	if s.rateLimiter != nil {
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return "", fmt.Errorf("rate limit exceeded: %v", err)
		}
	}

	run := func(m *genai.GenerativeModel) (string, error) {
		cs := m.StartChat()
		// Convert history
		var historyParts []*genai.Content
		for _, msg := range history {
			role := "user"
			if msg.Role == "assistant" {
				role = "model"
			}
			historyParts = append(historyParts, &genai.Content{Parts: []genai.Part{genai.Text(msg.Content)}, Role: role})
		}
		cs.History = historyParts
		resp, err := cs.SendMessage(ctx, genai.Text(prompt))
		if err != nil {
			return "", err
		}
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			return fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0]), nil
		}
		return "I'm sorry, I didn't get a response. Please try again.", nil
	}

	// First try current model
	ans, err := run(s.model)
	if err == nil || !isModelNotFound(err) {
		return ans, err
	}
	// Fallback models via SDK
	candidates := []string{Env.GeminiModel, "gemini-1.5-flash", "gemini-1.5-flash-8b", "gemini-1.5-pro"}
	ans, err = s.tryModels(ctx, candidates, run)
	if err == nil {
		return ans, nil
	}
	// As final fallback, call REST v1beta
	text, restErr := restGenerateText(ctx, Env.GeminiAPIKey, Env.GeminiModel, prompt, history)
	if restErr == nil {
		return text, nil
	}
	return "", fmt.Errorf("sdk fallback error: %v; rest error: %v", err, restErr)
}

func (s *aiService) AnalyzeCV(ctx context.Context, cvText string) (string, error) {
	if s.rateLimiter != nil {
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return "", fmt.Errorf("rate limit exceeded: %v", err)
		}
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

	run := func(m *genai.GenerativeModel) (string, error) {
		result, err := m.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return "", err
		}
		if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
			return fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]), nil
		}
		return "I couldn't analyze your CV at this time. Please try again later.", nil
	}

	ans, err := run(s.model)
	if err == nil || !isModelNotFound(err) {
		return ans, err
	}
	candidates := []string{Env.GeminiModel, "gemini-1.5-flash", "gemini-1.5-flash-8b", "gemini-1.5-pro"}
	ans, err = s.tryModels(ctx, candidates, run)
	if err == nil {
		return ans, nil
	}
	// REST fallback - pass the prompt, not raw cvText
	text, restErr := restGenerateText(ctx, Env.GeminiAPIKey, Env.GeminiModel, prompt, nil)
	if restErr == nil {
		return text, nil
	}
	return "", fmt.Errorf("sdk fallback error: %v; rest error: %v", err, restErr)
}

func (s *aiService) FindJobs(ctx context.Context, userProfile, query string) (string, error) {
	if s.rateLimiter != nil {
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return "", fmt.Errorf("rate limit exceeded: %v", err)
		}
	}

	prompt := fmt.Sprintf(`You are JobGen, an AI career assistant specializing in helping African professionals find remote tech jobs. 
	Based on the user profile and query below, provide personalized job search advice and suggestions for remote tech jobs that might be a good fit.
	
	User Profile: %s
	User Query: %s
	
	Provide your response with specific advice, potential job roles to explore, and tips for applying to remote positions. 
	Focus on opportunities that might be open to African candidates.`, userProfile, query)

	run := func(m *genai.GenerativeModel) (string, error) {
		result, err := m.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return "", err
		}
		if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
			response := fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0])
			disclaimer := "\n\n*Note: This is AI-generated advice based on your profile. For actual job listings, we recommend checking dedicated job platforms like LinkedIn, RemoteOK, and Indeed.*"
			return response + disclaimer, nil
		}
		return "I couldn't search for jobs at this time. Please try again later.", nil
	}

	ans, err := run(s.model)
	if err == nil || !isModelNotFound(err) {
		return ans, err
	}
	candidates := []string{Env.GeminiModel, "gemini-1.5-flash", "gemini-1.5-flash-8b", "gemini-1.5-pro"}
	ans, err = s.tryModels(ctx, candidates, run)
	if err == nil {
		return ans, nil
	}
	// REST fallback - pass the prompt, not the raw query
	text, restErr := restGenerateText(ctx, Env.GeminiAPIKey, Env.GeminiModel, prompt, nil)
	if restErr == nil {
		return text, nil
	}
	return "", fmt.Errorf("sdk fallback error: %v; rest error: %v", err, restErr)
}

func (s *aiService) ImproveCV(ctx context.Context, cv *domain.CV, userQuery string, history []domain.ChatMessage) (string, []domain.Suggestion, error) {
	if s.rateLimiter != nil {
		if err := s.rateLimiter.Wait(ctx); err != nil {
			return "", nil, fmt.Errorf("rate limit exceeded: %v", err)
		}
	}

	// Check if CV is nil
	if cv == nil {
		return "Please provide a CV to analyze and improve.", nil, nil
	}

	// Build context from conversation history
	conversationContext := s.buildConversationContext(history)

	// Convert CV to text representation
	cvText := s.formatCVForAI(cv)

	// Build specialized prompt based on the user query
	prompt := s.buildCVImprovementPrompt(cvText, userQuery, conversationContext)

	run := func(m *genai.GenerativeModel) (string, error) {
		result, err := m.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return "", err
		}
		if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
			return fmt.Sprintf("%v", result.Candidates[0].Content.Parts[0]), nil
		}
		return "I couldn't analyze your CV at this time. Please try again later.", nil
	}

	response, err := run(s.model)
	if err != nil && isModelNotFound(err) {
		candidates := []string{Env.GeminiModel, "gemini-1.5-flash", "gemini-1.5-flash-8b", "gemini-1.5-pro"}
		response, err = s.tryModels(ctx, candidates, run)
	}
	if err != nil {
		return "", nil, err
	}

	// Parse response to extract suggestions
	suggestions := s.extractSuggestionsFromResponse(response)

	return response, suggestions, nil
}

func (s *aiService) buildConversationContext(history []domain.ChatMessage) string {
	if len(history) == 0 {
		return "No previous conversation"
	}

	var context strings.Builder
	for _, msg := range history {
		role := "User"
		if msg.Role == "assistant" {
			role = "Assistant"
		}
		context.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}
	return context.String()
}

func (s *aiService) formatCVForAI(cv *domain.CV) string {
	var sb strings.Builder

	sb.WriteString("=== CV ANALYSIS REQUEST ===\n\n")
	sb.WriteString(fmt.Sprintf("Profile Summary: %s\n", cv.ProfileSummary))
	sb.WriteString(fmt.Sprintf("Current Score: %d/100\n\n", cv.Score))

	sb.WriteString("EXPERIENCES:\n")
	for i, exp := range cv.Experiences {
		sb.WriteString(fmt.Sprintf("%d. %s at %s\n", i+1, exp.Title, exp.Company))
		sb.WriteString(fmt.Sprintf("   Location: %s\n", exp.Location))

		// Handle nil EndDate (current job)
		endDateStr := "Present"
		if exp.EndDate != nil {
			endDateStr = exp.EndDate.Format("2006-01-02")
		}
		sb.WriteString(fmt.Sprintf("   Period: %s to %s\n", exp.StartDate.Format("2006-01-02"), endDateStr))
		sb.WriteString(fmt.Sprintf("   Description: %s\n\n", exp.Description))
	}

	sb.WriteString("EDUCATION:\n")
	for i, edu := range cv.Educations {
		sb.WriteString(fmt.Sprintf("%d. %s from %s\n", i+1, edu.Degree, edu.Institution))
		sb.WriteString(fmt.Sprintf("   Location: %s\n", edu.Location))
		sb.WriteString(fmt.Sprintf("   Graduation: %s\n\n", edu.GraduationDate.Format("2006")))
	}

	sb.WriteString("SKILLS:\n")
	for i, skill := range cv.Skills {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, skill))
	}

	sb.WriteString("\nEXISTING SUGGESTIONS:\n")
	for i, suggestion := range cv.Suggestions {
		status := "PENDING"
		if suggestion.Applied {
			status = "APPLIED"
		}
		sb.WriteString(fmt.Sprintf("%d. [%s] %s: %s\n",
			i+1, status, suggestion.Type, suggestion.Content))
	}

	return sb.String()
}

func (s *aiService) buildCVImprovementPrompt(cvText, userQuery, conversationContext string) string {
	prompt := `You are JobGen, an AI career assistant specializing in helping African professionals improve their CVs for remote tech jobs.

Below is the user's CV information:
%s
----------------------------------------
CONVERSATION CONTEXT:
%s
----------------------------------------
USER'S CURRENT REQUEST:
%s

Please provide specific, actionable suggestions to improve this CV. Focus on:
1. Adding quantifiable metrics to achievements
2. Using strong action verbs
3. Including relevant technical keywords
4. Improving structure for Applicant Tracking Systems (ATS)
5. Tailoring for international remote tech jobs

Format your response with clear sections and bullet points. Provide specific examples whenever possible.`

	return fmt.Sprintf(prompt, cvText, conversationContext, userQuery)
}

func (s *aiService) extractSuggestionsFromResponse(response string) []domain.Suggestion {
	var suggestions []domain.Suggestion

	// Simple parsing logic
	lines := strings.Split(response, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Extract bullet points as suggestions
		if (strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "• ") || strings.HasPrefix(line, "* ")) && len(line) > 2 {
			content := strings.TrimPrefix(line, "- ")
			content = strings.TrimPrefix(content, "• ")
			content = strings.TrimPrefix(content, "* ")
			content = strings.TrimSpace(content)

			if content != "" {
				suggestions = append(suggestions, domain.Suggestion{
					Type:    "cv_improvement",
					Content: content,
					Applied: false,
				})
			}
		}
	}

	// If no structured suggestions found, use the whole response as a suggestion
	if len(suggestions) == 0 {
		suggestions = append(suggestions, domain.Suggestion{
			Type:    "cv_improvement",
			Content: response,
			Applied: false,
		})
	}

	return suggestions
}
