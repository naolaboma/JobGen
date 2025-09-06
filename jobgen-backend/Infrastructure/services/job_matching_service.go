package services

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type JobMatchingService struct {
	jobRepo  domain.IJobRepository
	userRepo domain.IUserRepository
}

func NewJobMatchingService(jobRepo domain.IJobRepository, userRepo domain.IUserRepository) domain.IJobMatchingService {
	return &JobMatchingService{
		jobRepo:  jobRepo,
		userRepo: userRepo,
	}
}

func (j *JobMatchingService) CalculateMatchScore(job domain.Job, preferences domain.UserJobPreferences) float64 {
	var totalScore float64
	
	// Skills matching (70% weight)
	skillScore := j.calculateSkillScore(job.ExtractedSkills, preferences.Skills)
	totalScore += skillScore * 0.7
	
	// Experience matching (20% weight)
	expScore := j.calculateExperienceScore(job.Description, preferences.ExperienceYears)
	totalScore += expScore * 0.2
	
	// Location matching (10% weight)
	locationScore := j.calculateLocationScore(job.Location, preferences.Locations)
	totalScore += locationScore * 0.1
	
	// Ensure score is between 0 and 100
	totalScore = math.Min(100, math.Max(0, totalScore))
	
	return totalScore
}

func (j *JobMatchingService) calculateSkillScore(jobSkills, userSkills []string) float64 {
	if len(userSkills) == 0 {
		return 0
	}
	
	// Convert to lowercase for case-insensitive matching
	jobSkillsLower := make(map[string]bool)
	for _, skill := range jobSkills {
		jobSkillsLower[strings.ToLower(skill)] = true
	}
	
	matchedSkills := 0
	for _, userSkill := range userSkills {
		if jobSkillsLower[strings.ToLower(userSkill)] {
			matchedSkills++
		}
	}
	
	// Calculate percentage of user skills that match
	return float64(matchedSkills) / float64(len(userSkills)) * 100
}

func (j *JobMatchingService) calculateExperienceScore(jobDescription string, userExperience int) float64 {
	if userExperience == 0 {
		return 50 // Neutral score for entry level
	}
	
	// Extract experience requirements from job description
	requiredExp := j.extractExperienceRequirement(jobDescription)
	
	if requiredExp == 0 {
		return 75 // Good score if no specific requirement
	}
	
	// Calculate score based on how well user experience matches requirement
	diff := float64(userExperience - requiredExp)
	
	if diff >= 0 {
		// User has more or equal experience than required
		if diff <= 2 {
			return 100 // Perfect match
		} else if diff <= 5 {
			return 85 // Still very good
		} else {
			return 70 // Might be overqualified but still good
		}
	} else {
		// User has less experience than required
		deficit := -diff
		if deficit <= 1 {
			return 80 // Close enough
		} else if deficit <= 3 {
			return 60 // Some gap but manageable
		} else {
			return 30 // Significant gap
		}
	}
}

func (j *JobMatchingService) extractExperienceRequirement(description string) int {
	description = strings.ToLower(description)
	
	// Common patterns for experience requirements
	patterns := []string{
		`(\d+)\+?\s*years?\s*of?\s*experience`,
		`(\d+)\+?\s*years?\s*experience`,
		`minimum\s*(\d+)\s*years?`,
		`at least\s*(\d+)\s*years?`,
		`(\d+)\s*to\s*\d+\s*years?`,
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(description)
		if len(matches) > 1 {
			if exp, err := strconv.Atoi(matches[1]); err == nil {
				return exp
			}
		}
	}
	
	// Check for seniority levels
	if strings.Contains(description, "senior") || strings.Contains(description, "lead") {
		return 5
	} else if strings.Contains(description, "mid-level") || strings.Contains(description, "intermediate") {
		return 3
	} else if strings.Contains(description, "junior") || strings.Contains(description, "entry") {
		return 0
	}
	
	return 0 // No specific requirement found
}

func (j *JobMatchingService) calculateLocationScore(jobLocation string, preferredLocations []string) float64 {
	if len(preferredLocations) == 0 {
		return 100 // No preference means all locations are fine
	}
	
	jobLocationLower := strings.ToLower(jobLocation)
	
	// Remote work gets high score
	if strings.Contains(jobLocationLower, "remote") || 
	   strings.Contains(jobLocationLower, "anywhere") ||
	   strings.Contains(jobLocationLower, "worldwide") {
		return 100
	}
	
	// Check against preferred locations
	for _, prefLocation := range preferredLocations {
		prefLocationLower := strings.ToLower(prefLocation)
		
		// Exact match or partial match
		if strings.Contains(jobLocationLower, prefLocationLower) || 
		   strings.Contains(prefLocationLower, jobLocationLower) {
			return 100
		}
	}
	
	return 20 // Location doesn't match preferences
}

func (j *JobMatchingService) GetMatchedJobs(ctx context.Context, userID string, limit int, offset int) ([]domain.Job, error) {
	// Get user preferences
	user, err := j.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	// Create user preferences from user profile
	preferences := domain.UserJobPreferences{
		Skills:          user.Skills,
		ExperienceYears: user.ExperienceYears,
		Locations:       []string{user.Location}, // Can be expanded to support multiple preferred locations
	}
	
	// Get jobs for matching (you might want to implement pagination here too)
	jobs, err := j.jobRepo.GetJobsForMatching(ctx, limit*2, offset) // Get more jobs to allow for filtering
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs for matching: %w", err)
	}
	
	// Calculate match scores and filter
	var matchedJobs []domain.Job
	for _, job := range jobs {
		score := j.CalculateMatchScore(job, preferences)
		
		// Only include jobs with score above threshold
		if score >= 30 { // 30% minimum match
			job.MatchScore = &score
			matchedJobs = append(matchedJobs, job)
		}
		
		// Stop when we have enough matches
		if len(matchedJobs) >= limit {
			break
		}
	}
	
	// Sort by match score (highest first)
	j.sortJobsByMatchScore(matchedJobs)
	
	return matchedJobs, nil
}

func (j *JobMatchingService) sortJobsByMatchScore(jobs []domain.Job) {
	// Simple bubble sort by match score (for small arrays)
	n := len(jobs)
	for i := 0; i < n-1; i++ {
		for k := 0; k < n-i-1; k++ {
			score1 := float64(0)
			score2 := float64(0)
			
			if jobs[k].MatchScore != nil {
				score1 = *jobs[k].MatchScore
			}
			if jobs[k+1].MatchScore != nil {
				score2 = *jobs[k+1].MatchScore
			}
			
			if score1 < score2 {
				jobs[k], jobs[k+1] = jobs[k+1], jobs[k]
			}
		}
	}
}

func (j *JobMatchingService) UpdateUserPreferences(ctx context.Context, userID string, preferences domain.UserJobPreferences) error {
	// Get current user
	user, err := j.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	
	// Update user profile with new preferences
	user.Skills = preferences.Skills
	user.ExperienceYears = preferences.ExperienceYears
	
	// Handle locations (for now, just take the first one as primary location)
	if len(preferences.Locations) > 0 {
		user.Location = preferences.Locations[0]
	}
	
	// Update user in repository
	if err := j.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}
	
	return nil
}

// GetJobRecommendations provides more advanced recommendations
func (j *JobMatchingService) GetJobRecommendations(ctx context.Context, userID string, limit int) ([]domain.Job, error) {
	// This could be expanded with ML algorithms in the future
	// For now, it's similar to GetMatchedJobs but with different scoring weights
	return j.GetMatchedJobs(ctx, userID, limit, 0)
}

// AnalyzeJobMarket provides insights about the job market based on user skills
func (j *JobMatchingService) AnalyzeJobMarket(ctx context.Context, skills []string) (map[string]interface{}, error) {
	// Get recent jobs for analysis
	filter := domain.JobFilter{
		Skills:    skills,
		Page:      1,
		Limit:     100,
		SortBy:    "posted_at",
		SortOrder: "desc",
	}
	
	jobs, total, err := j.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs for analysis: %w", err)
	}
	
	// Analyze the data
	analysis := map[string]interface{}{
		"total_matching_jobs": total,
		"analyzed_jobs":       len(jobs),
		"top_companies":       j.getTopCompanies(jobs, 10),
		"common_skills":       j.getMostCommonSkills(jobs, 15),
		"salary_insights":     j.analyzeSalaries(jobs),
		"location_insights":   j.analyzeLocations(jobs),
		"source_distribution": j.analyzeJobSources(jobs),
	}
	
	return analysis, nil
}

func (j *JobMatchingService) getTopCompanies(jobs []domain.Job, limit int) []map[string]interface{} {
	companyCount := make(map[string]int)
	
	for _, job := range jobs {
		if job.CompanyName != "" && job.CompanyName != "Unknown Company" {
			companyCount[job.CompanyName]++
		}
	}
	
	// Convert to slice and sort
	type companyInfo struct {
		name  string
		count int
	}
	
	var companies []companyInfo
	for name, count := range companyCount {
		companies = append(companies, companyInfo{name: name, count: count})
	}
	
	// Simple sort by count
	for i := 0; i < len(companies)-1; i++ {
		for k := 0; k < len(companies)-i-1; k++ {
			if companies[k].count < companies[k+1].count {
				companies[k], companies[k+1] = companies[k+1], companies[k]
			}
		}
	}
	
	// Convert to result format
	var result []map[string]interface{}
	maxItems := limit
	if len(companies) < maxItems {
		maxItems = len(companies)
	}
	
	for i := 0; i < maxItems; i++ {
		result = append(result, map[string]interface{}{
			"company": companies[i].name,
			"count":   companies[i].count,
		})
	}
	
	return result
}

func (j *JobMatchingService) getMostCommonSkills(jobs []domain.Job, limit int) []map[string]interface{} {
	skillCount := make(map[string]int)
	
	for _, job := range jobs {
		for _, skill := range job.ExtractedSkills {
			if skill != "" {
				skillCount[strings.ToLower(skill)]++
			}
		}
	}
	
	// Convert to slice and sort
	type skillInfo struct {
		name  string
		count int
	}
	
	var skills []skillInfo
	for name, count := range skillCount {
		skills = append(skills, skillInfo{name: name, count: count})
	}
	
	// Sort by count
	for i := 0; i < len(skills)-1; i++ {
		for k := 0; k < len(skills)-i-1; k++ {
			if skills[k].count < skills[k+1].count {
				skills[k], skills[k+1] = skills[k+1], skills[k]
			}
		}
	}
	
	// Convert to result format
	var result []map[string]interface{}
	maxItems := limit
	if len(skills) < maxItems {
		maxItems = len(skills)
	}
	
	for i := 0; i < maxItems; i++ {
		result = append(result, map[string]interface{}{
			"skill": skills[i].name,
			"count": skills[i].count,
		})
	}
	
	return result
}

func (j *JobMatchingService) analyzeSalaries(jobs []domain.Job) map[string]interface{} {
	var salariesWithData []string
	
	for _, job := range jobs {
		if job.Salary != "" {
			salariesWithData = append(salariesWithData, job.Salary)
		}
	}
	
	return map[string]interface{}{
		"total_jobs_with_salary": len(salariesWithData),
		"percentage_with_salary": float64(len(salariesWithData)) / float64(len(jobs)) * 100,
		"sample_salaries":       salariesWithData[:min(5, len(salariesWithData))],
	}
}

func (j *JobMatchingService) analyzeLocations(jobs []domain.Job) map[string]interface{} {
	locationCount := make(map[string]int)
	
	for _, job := range jobs {
		location := strings.ToLower(job.Location)
		locationCount[location]++
	}
	
	// Find top locations
	type locationInfo struct {
		name  string
		count int
	}
	
	var locations []locationInfo
	for name, count := range locationCount {
		locations = append(locations, locationInfo{name: name, count: count})
	}
	
	// Sort by count
	for i := 0; i < len(locations)-1; i++ {
		for k := 0; k < len(locations)-i-1; k++ {
			if locations[k].count < locations[k+1].count {
				locations[k], locations[k+1] = locations[k+1], locations[k]
			}
		}
	}
	
	// Convert to result format
	var topLocations []map[string]interface{}
	maxItems := min(10, len(locations))
	
	for i := 0; i < maxItems; i++ {
		topLocations = append(topLocations, map[string]interface{}{
			"location": locations[i].name,
			"count":    locations[i].count,
		})
	}
	
	return map[string]interface{}{
		"top_locations":    topLocations,
		"remote_job_count": locationCount["remote"] + locationCount["anywhere"] + locationCount["worldwide"],
	}
}

func (j *JobMatchingService) analyzeJobSources(jobs []domain.Job) map[string]interface{} {
	sourceCount := make(map[string]int)
	
	for _, job := range jobs {
		sourceCount[job.Source]++
	}
	
	return map[string]interface{}{
		"sources": sourceCount,
		"total":   len(jobs),
	}
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
