package usecases

import (
	domain "jobgen-backend/Domain"
	"regexp"
	"strings"
)

// --- Scoring Logic ---
var severityWeights = map[string]int{
	"quantification":    20,
	"weak_action_verbs": 10,
	"missing_keywords":  15,
}

// CalculateScore computes the CV score based on suggestions.
func CalculateScore(suggestions []domain.Suggestion) int {
	score := 100
	for _, suggestion := range suggestions {
		if deduction, exists := severityWeights[suggestion.Type]; exists {
			score -= deduction
		}
	}
	if score < 0 {
		return 0
	}
	return score
}

// --- Parsing Pipeline Logic ---
var (
	experienceHeaderRegex = regexp.MustCompile(`(?i)\b(experience|work history|employment)\b`)
	educationHeaderRegex  = regexp.MustCompile(`(?i)\b(education|academic background)\b`)
	skillsHeaderRegex     = regexp.MustCompile(`(?i)\b(skills|technical proficiencies)\b`)
)

// ParseTextToCVSections parses raw text into structured CV sections.
func ParseTextToCVSections(rawText string) (*domain.CV, error) {
	cv := &domain.CV{}
	lines := strings.Split(rawText, "\n")
	var currentSection string
	var sectionContent strings.Builder

	processSection := func() {
		content := strings.TrimSpace(sectionContent.String())
		if content == "" {
			return
		}
		switch currentSection {
		case "skills":
			skills := strings.FieldsFunc(content, func(r rune) bool {
				return r == ',' || r == '\n' || r == '•'
			})
			skillSet := make(map[string]bool)
			for _, skill := range skills {
				sanitizedSkill := strings.TrimSpace(skill)
				if sanitizedSkill != "" && len(sanitizedSkill) < 50 { // Basic sanity check
					skillSet[sanitizedSkill] = true
				}
			}
			for skill := range skillSet {
				cv.Skills = append(cv.Skills, skill)
			}
		}
		sectionContent.Reset()
	}

	for _, line := range lines {
		switch {
		case experienceHeaderRegex.MatchString(line):
			processSection()
			currentSection = "experience"
		case educationHeaderRegex.MatchString(line):
			processSection()
			currentSection = "education"
		case skillsHeaderRegex.MatchString(line):
			processSection()
			currentSection = "skills"
		default:
			if currentSection != "" {
				sectionContent.WriteString(line + "\n")
			}
		}
	}
	processSection() // Process the last section

	return cv, nil
}
