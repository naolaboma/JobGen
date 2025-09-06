package usecases

import (
	domain "jobgen-backend/Domain"
	"regexp"
	"sort"
	"strings"
	"time"
)

// --- Scoring Logic ---
// severityWeights controls the initial deduction per suggestion type (before diminishing returns).
// Tuned to be less punitive.
var severityWeights = map[string]int{
	"quantification":    12,
	"weak_action_verbs": 6,
	"missing_keywords":  8,
}

// perTypeCaps limits the maximum total deduction per suggestion type.
var perTypeCaps = map[string]int{
	"quantification":    24,
	"weak_action_verbs": 12,
	"missing_keywords":  16,
}

// CalculateScore computes the CV score based on suggestions.
func CalculateScore(suggestions []domain.Suggestion) int {
	const overallMaxDeduction = 40 // don't penalize beyond this overall

	// Diminishing returns per type + per-type caps
	counts := make(map[string]int)
	appliedPerType := make(map[string]int)
	totalDeduction := 0

	for _, s := range suggestions {
		base, ok := severityWeights[s.Type]
		if !ok {
			continue // ignore unknown types
		}
		counts[s.Type]++
		n := counts[s.Type]
		// integer decay by powers of 2: base >> (n-1)
		deduction := base >> (n - 1)
		if deduction < 1 {
			deduction = 1
		}
		// apply per-type cap
		capLeft := perTypeCaps[s.Type] - appliedPerType[s.Type]
		if capLeft <= 0 {
			continue
		}
		if deduction > capLeft {
			deduction = capLeft
		}
		// apply overall cap
		if totalDeduction+deduction > overallMaxDeduction {
			deduction = overallMaxDeduction - totalDeduction
		}
		if deduction <= 0 {
			break
		}
		appliedPerType[s.Type] += deduction
		totalDeduction += deduction
		if totalDeduction >= overallMaxDeduction {
			break
		}
	}
	score := 100 - totalDeduction
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	return score
}

var (
	experienceHeaderRegex = regexp.MustCompile(`(?i)\b(experience|work\s*history|employment|professional\s*experience|work\s*experience)\b`)
	educationHeaderRegex  = regexp.MustCompile(`(?i)\b(education|academic\s*background|academics|qualifications)\b`)
	skillsHeaderRegex     = regexp.MustCompile(`(?i)\b(skills|technical\s*proficiencies|technical\s*skills|key\s*skills)\b`)
	monthShortRE = `(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Sept|Oct|Nov|Dec)`
	monthLongRE  = `(January|February|March|April|May|June|July|August|September|October|November|December)`
	dateTokenRE  = regexp.MustCompile(`(?i)\b((` + monthShortRE + `|` + monthLongRE + `)\s+\d{4}|\d{4}[-/.]\d{1,2}|\d{1,2}[-/.]\d{4}|\d{4})\b`)
)

// ParseTextToCVSections parses raw text into structured CV sections.
func ParseTextToCVSections(rawText string) (*domain.CV, error) {
	cv := &domain.CV{}
	normalized := normalizeText(rawText)
	lines := strings.Split(normalized, "\n")
	var currentSection string
	var sectionContent strings.Builder

	processSection := func() {
		content := strings.TrimSpace(sectionContent.String())
		if content == "" {
			return
		}
		switch currentSection {
		case "skills":
			cv.Skills = dedupeStrings(parseSkills(content, 200))
		case "experience":
			cv.Experiences = append(cv.Experiences, parseExperienceBlock(content)...)
		case "education":
			cv.Educations = append(cv.Educations, parseEducationBlock(content)...)
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

// normalizeText makes line endings consistent and collapses repeated whitespace.
func normalizeText(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	
	s = strings.ReplaceAll(s, "\t", " ")
	s = regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
	return s
}

// parseSkills splits a skills blob into a deduplicated, sanitized slice, preserving first-seen casing.
func parseSkills(content string, capLen int) []string {
	splitFn := func(r rune) bool {
		switch r {
		case ',', '\n', '|', ';', '•', '·':
			return true
		case '-', '–', '—':
			return false
		default:
			return false
		}
	}

	replacers := []struct{ old, new string }{
		{"|", ","}, {"•", ","}, {"·", ","}, {";", ","},
	}
	for _, r := range replacers {
		content = strings.ReplaceAll(content, r.old, r.new)
	}

	raw := strings.FieldsFunc(content, splitFn)
	seen := make(map[string]struct{})
	out := make([]string, 0, len(raw))
	for _, token := range raw {
		s := sanitizeSkill(token)
		if s == "" || len(s) > 50 { // sanity length
			continue
		}
		key := strings.ToLower(s)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, s)
		if capLen > 0 && len(out) >= capLen {
			break
		}
	}
	sort.Strings(out)
	return out
}

func sanitizeSkill(s string) string {
	s = strings.TrimSpace(s)
	// Strip leading common bullets and dashes
	s = strings.TrimLeft(s, "•-*–—· ")
	s = strings.TrimSpace(s)
	// Remove trailing punctuation
	s = strings.TrimRight(s, ",.;: ")
	// Remove surrounding parentheses
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") && len(s) > 2 {
		s = strings.TrimSuffix(strings.TrimPrefix(s, "("), ")")
	}
	// Collapse multiple spaces
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
	return s
}

// --- Experience/Education parsing ---
func parseExperienceBlock(content string) []domain.Experience {
	var out []domain.Experience
	// Split by blank lines to approximate entries
	blocks := regexp.MustCompile(`\n\s*\n+`).Split(strings.TrimSpace(content), -1)
	for _, b := range blocks {
		entry := strings.TrimSpace(b)
		if entry == "" {
			continue
		}
		lines := nonEmptyLines(entry)
		if len(lines) == 0 {
			continue
		}
		var title, company, location, desc string
		var start, end *time.Time
		header := lines[0]

		parts := regexp.MustCompile(`\s[-–—]| at `).Split(header, 2)
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			if looksLikeRole(left) && !looksLikeRole(right) {
				title = left
				company = right
			} else if looksLikeRole(right) {
				title = right
				company = left
			} else {
				company, title = left, right
			}
		} else {
			title = strings.TrimSpace(header)
		}
		// Scan for date tokens in first 2 lines
		joined := strings.Join(lines[:min(2, len(lines))], " ")
		startT, endT := extractDateRange(joined)
		if startT != nil {
			start = startT
		}
		if endT != nil {
			end = endT
		}
		
		for _, ln := range lines[:min(3, len(lines))] {
			loc := pickLocation(ln)
			if loc != "" {
				location = loc
				break
			}
		}
		// Description: remaining lines
		if len(lines) > 1 {
			desc = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		}
		out = append(out, domain.Experience{
			ID:          "exp-" + hashString(header),
			Title:       title,
			Company:     normalizeName(company),
			Location:    location,
			StartDate:   derefOrZero(start),
			EndDate:     end,
			Description: desc,
		})
	}
	return out
}

func parseEducationBlock(content string) []domain.Education {
	var out []domain.Education
	blocks := regexp.MustCompile(`\n\s*\n+`).Split(strings.TrimSpace(content), -1)
	for _, b := range blocks {
		entry := strings.TrimSpace(b)
		if entry == "" {
			continue
		}
		lines := nonEmptyLines(entry)
		if len(lines) == 0 {
			continue
		}
		header := lines[0]
		// Heuristic: Degree, Institution in first line(s)
		degree := pickDegree(header)
		institution := header
		if degree != "" {
			institution = strings.TrimSpace(strings.Replace(header, degree, "", 1))
		}
		// If degree empty, check second line
		if degree == "" && len(lines) > 1 {
			degree = pickDegree(lines[1])
			if degree != "" {
				institution = lines[0]
			}
		}
		// Graduation date: scan first 2-3 lines
		joined := strings.Join(lines[:min(3, len(lines))], " ")
		grad := extractSingleDate(joined)
		out = append(out, domain.Education{
			ID:             "edu-" + hashString(header),
			Degree:         strings.TrimSpace(degree),
			Institution:    normalizeName(strings.TrimSpace(institution)),
			Location:       pickLocation(joined),
			GraduationDate: derefOrZero(grad),
		})
	}
	return out
}

// --- helpers ---
func dedupeStrings(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		k := strings.ToLower(strings.TrimSpace(s))
		if k == "" {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, s)
	}
	return out
}

func looksLikeRole(s string) bool {
	s = strings.ToLower(s)
	roleHints := []string{"engineer", "developer", "manager", "lead", "architect", "intern", "specialist", "analyst"}
	for _, h := range roleHints {
		if strings.Contains(s, h) {
			return true
		}
	}
	return false
}

func extractDateRange(s string) (*time.Time, *time.Time) {
	matches := dateTokenRE.FindAllString(s, -1)
	if len(matches) == 0 {
		return nil, nil
	}
	if len(matches) == 1 {
		t := parseAnyDate(matches[0])
		return t, nil
	}
	start := parseAnyDate(matches[0])
	end := parseAnyDate(matches[1])
	return start, end
}

func extractSingleDate(s string) *time.Time {
	m := dateTokenRE.FindString(s)
	if m == "" {
		return nil
	}
	return parseAnyDate(m)
}

func parseAnyDate(tok string) *time.Time {
	tok = strings.TrimSpace(tok)
	// Try formats in order
	fmts := []string{
		"Jan 2006", "January 2006", "2006-01", "2006/01", "01/2006", "2006",
	}
	for _, f := range fmts {
		if t, err := time.Parse(f, normalizeDateToken(tok)); err == nil {
			return &t
		}
	}
	// If a bare year like 2020
	if yearRE := regexp.MustCompile(`^\d{4}$`); yearRE.MatchString(tok) {
		t, _ := time.Parse("2006", tok)
		return &t
	}
	return nil
}

func normalizeDateToken(tok string) string {
	// Normalize Sept to Sep
	tok = strings.ReplaceAll(tok, "Sept", "Sep")
	tok = strings.ReplaceAll(tok, "sept", "Sep")
	// Collapse extra spaces
	tok = regexp.MustCompile(`\s{2,}`).ReplaceAllString(tok, " ")
	return tok
}

func pickLocation(s string) string {
	// Very light heuristic: look for ", City" or country/state words
	if i := strings.Index(s, ","); i > 0 {
		cand := strings.TrimSpace(s[i+1:])
		if len(cand) >= 2 && len(cand) <= 40 {
			return cand
		}
	}
	locHints := []string{"remote", "usa", "ethiopia", "uk", "canada", "germany", "netherlands"}
	for _, h := range locHints {
		if strings.Contains(strings.ToLower(s), h) {
			return h
		}
	}
	return ""
}

func pickDegree(s string) string {
	s = strings.ToLower(s)
	degrees := []string{"bsc", "msc", "bs", "ms", "phd", "bachelor", "master", "doctor", "diploma"}
	for _, d := range degrees {
		if strings.Contains(s, d) {
			return d
		}
	}
	return ""
}

func normalizeName(s string) string {
	// Basic normalization: collapse spaces and uppercase initials
	s = strings.TrimSpace(s)
	s = regexp.MustCompile(`\s{2,}`).ReplaceAllString(s, " ")
	return s
}

func nonEmptyLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := lines[:0]
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func derefOrZero(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}
func hashString(s string) string {
	// Fast non-cryptographic short hash for IDs
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return strings.ToLower(base36(uint64(h)))
}
func base36(x uint64) string {
	const digits = "0123456789abcdefghijklmnopqrstuvwxyz"
	if x == 0 {
		return "0"
	}
	b := make([]byte, 0, 8)
	for x > 0 {
		b = append([]byte{digits[x%36]}, b...)
		x /= 36
	}
	return string(b)
}
