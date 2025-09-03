package usecases

import (
	"jobgen-backend/Domain"
	"regexp"
	"strings"
)

// --- Scoring Logic ---
var severityWeights = map[string]int{
	"quantification":    20,
	"weak_action_verbs": 10,
	"missing_keywords":  15,
}

func CalculateScore(suggestions []Domain.Suggestion) int {
	score := 100
	for _, suggestion := range suggestions {
		if deduction, ok := severityWeights[suggestion.Type]; ok {
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
	expHeaderRegex    = regexp.MustCompile(`(?i)\b(experience|work history|employment)\b`)
	eduHeaderRegex    = regexp.MustCompile(`(?i)\b(education|academic background)\b`)
	skillsHeaderRegex = regexp.MustCompile(`(?i)\b(skills|technical proficiencies)\b`)
)

func ParseTextToCVSections(rawText string) (*Domain.CV, error) {
	// This is a simplified heuristic parser
	// A production system might use more advanced NLP techniques
	cv := &Domain.CV{}
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
			skills := strings.FieldsFunc(content, func(r rune) { return r == ',' || r == '\n' || r == '•' })
			skillMap := make(map[string]bool)
			for _, skill := range skills {
				s := strings.TrimSpace(skill)
				if s != "" && len(s) < 50 { // Basic sanity check
					skillMap[s] = true
				}
			}
			for s := range skillMap {
				cv.Skills = append(cv.Skills, s)
			}
		}
		sectionContent.Reset()
	}

	for _, line := range lines {
		if expHeaderRegex.MatchString(line) {
			processSection()
			currentSection = "experience"
			continue
		}
		if eduHeaderRegex.MatchString(line) {
			processSection()
			currentSection = "education"
			continue
		}
		if skillsHeaderRegex.MatchString(line) {
			processSection()
			currentSection = "skills"
			continue
		}

		if currentSection != "" {
			sectionContent.WriteString(line + "\n")
		}
	}
	processSection() // Process the last section

	return cv, nil
}
