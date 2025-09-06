package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	domain "jobgen-backend/Domain"
)

func NewAIServiceClient() domain.AIService { // satisfies domain.AIService
	// Use only Gemini. No local fallback.
	return &geminiAIServiceClient{
		apiKey:     Env.GeminiAPIKey,
		model:      Env.GeminiModel,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

// ---------------- Gemini AI client ----------------

type geminiAIServiceClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

type geminiContent struct {
	Role  string `json:"role,omitempty"`
	Parts []struct {
		Text string `json:"text"`
	} `json:"parts"`
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

func (c *geminiAIServiceClient) AnalyzeCV(rawText string) ([]domain.Suggestion, error) {
	if c.apiKey == "" {
		return nil, errors.New("gemini api key missing")
	}
	model := c.model
	if model == "" {
		model = "gemini-1.5-flash"
	}
	endpoint := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:generateContent?key=%s", model, c.apiKey)

	instruction := "You are an assistant that extracts CV improvement suggestions. Output ONLY a compact JSON array with objects: {id: string, type: one of [\"quantification\", \"weak_action_verbs\", \"missing_keywords\"], content: string, applied: false}. No markdown, no extra text."

	if len(rawText) > 8000 {
		rawText = rawText[:8000]
	}
	payload := geminiRequest{
		Contents: []geminiContent{
			{Role: "user", Parts: []struct {
				Text string `json:"text"`
			}{{Text: instruction + "\nCV:\n" + rawText}}},
		},
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		
		b, _ := io.ReadAll(resp.Body)
		// try to decode standard error shape
		var er geminiResponse
		_ = json.Unmarshal(b, &er)
		if er.Error != nil {
			return nil, fmt.Errorf("gemini http %d: %s", resp.StatusCode, er.Error.Message)
		}
		msg := strings.TrimSpace(string(b))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		if len(msg) > 300 {
			msg = msg[:300]
		}
		return nil, fmt.Errorf("gemini http %d: %s", resp.StatusCode, msg)
	}

	var gr geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}
	if gr.Error != nil {
		return nil, fmt.Errorf("gemini error: %s", gr.Error.Message)
	}
	if len(gr.Candidates) == 0 || len(gr.Candidates[0].Content.Parts) == 0 {
		return nil, errors.New("gemini returned no candidates")
	}
	text := gr.Candidates[0].Content.Parts[0].Text
	// Extract JSON from potential formatting and decode
	jsonStr := extractJSON(text)

	var suggestions []domain.Suggestion
	if err := json.Unmarshal([]byte(jsonStr), &suggestions); err != nil || len(suggestions) == 0 {
		return nil, fmt.Errorf("gemini invalid json response")
	}
	// Ensure Applied defaults to false
	for i := range suggestions {
		suggestions[i].Applied = false
	}
	return suggestions, nil
}

// extractJSON tries to isolate the first JSON array in text.
func extractJSON(s string) string {
	s = strings.TrimSpace(s)

	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	
	start := strings.Index(s, "[")
	end := strings.LastIndex(s, "]")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
