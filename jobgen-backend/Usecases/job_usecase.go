package usecases

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"time"
)

type jobUsecase struct {
	jobRepo              domain.IJobRepository
	userRepo             domain.IUserRepository
	jobAggregationSvc    domain.IJobAggregationService
	jobMatchingSvc       domain.IJobMatchingService
	contextTimeout       time.Duration
}

func NewJobUsecase(
	jobRepo domain.IJobRepository,
	userRepo domain.IUserRepository,
	jobAggregationSvc domain.IJobAggregationService,
	jobMatchingSvc domain.IJobMatchingService,
	timeout time.Duration,
) domain.IJobUsecase {
	return &jobUsecase{
		jobRepo:           jobRepo,
		userRepo:          userRepo,
		jobAggregationSvc: jobAggregationSvc,
		jobMatchingSvc:    jobMatchingSvc,
		contextTimeout:    timeout,
	}
}

func (j *jobUsecase) GetJobs(ctx context.Context, filter domain.JobFilter) (*domain.PaginatedJobsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 10
	}
	if filter.SortBy == "" {
		filter.SortBy = "posted_at"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// Get jobs from repository
	jobs, total, err := j.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get jobs: %w", err)
	}

	// Calculate pagination info
	totalPages := int((total + int64(filter.Limit) - 1) / int64(filter.Limit))

	response := &domain.PaginatedJobsResponse{
		Jobs:       jobs,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    filter.Page < totalPages,
		HasPrev:    filter.Page > 1,
	}

	return response, nil
}

func (j *jobUsecase) GetJobByID(ctx context.Context, id string) (*domain.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	if id == "" {
		return nil, fmt.Errorf("job ID is required")
	}

	job, err := j.jobRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

func (j *jobUsecase) SearchJobs(ctx context.Context, userID string, filter domain.JobFilter) (*domain.PaginatedJobsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	// If user is provided and no specific skills filter, use user's skills
	if userID != "" && len(filter.Skills) == 0 {
		user, err := j.userRepo.GetByID(ctx, userID)
		if err == nil && len(user.Skills) > 0 {
			filter.Skills = user.Skills
		}
	}

	// Get filtered jobs
	response, err := j.GetJobs(ctx, filter)
	if err != nil {
		return nil, err
	}

	// If user ID is provided, calculate match scores for jobs
	if userID != "" {
		user, err := j.userRepo.GetByID(ctx, userID)
		if err == nil {
			preferences := domain.UserJobPreferences{
				Skills:          user.Skills,
				ExperienceYears: user.ExperienceYears,
				Locations:       []string{user.Location},
			}

			for i := range response.Jobs {
				score := j.jobMatchingSvc.CalculateMatchScore(response.Jobs[i], preferences)
				response.Jobs[i].MatchScore = &score
			}

			// Sort by match score if no specific sort order was requested
			if filter.SortBy == "posted_at" || filter.SortBy == "" {
				j.sortJobsByMatchScore(response.Jobs)
			}
		}
	}

	return response, nil
}

func (j *jobUsecase) GetMatchedJobs(ctx context.Context, userID string, limit int, offset int) (*domain.PaginatedJobsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	if userID == "" {
		return nil, fmt.Errorf("user ID is required for matched jobs")
	}

	// Set defaults
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Get matched jobs from matching service
	matchedJobs, err := j.jobMatchingSvc.GetMatchedJobs(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get matched jobs: %w", err)
	}

	// Calculate pagination info
	page := (offset / limit) + 1
	totalPages := 1 // We don't have total count from matching service
	hasNext := len(matchedJobs) == limit
	hasPrev := page > 1

	response := &domain.PaginatedJobsResponse{
		Jobs:       matchedJobs,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(matchedJobs)), // Approximate
		TotalPages: totalPages,              // Approximate
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return response, nil
}

func (j *jobUsecase) AggregateJobs(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute) // Longer timeout for aggregation
	defer cancel()

	err := j.jobAggregationSvc.AggregateFromAllSources(ctx)
	if err != nil {
		return fmt.Errorf("failed to aggregate jobs: %w", err)
	}

	return nil
}

func (j *jobUsecase) GetJobSources(ctx context.Context) ([]domain.JobScrapeSource, error) {
	sources := j.jobAggregationSvc.GetSupportedSources()
	return sources, nil
}

// Additional methods for job management

func (j *jobUsecase) CreateJob(ctx context.Context, job *domain.Job) error {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	if job.Title == "" || job.CompanyName == "" || job.ApplyURL == "" {
		return fmt.Errorf("title, company name, and apply URL are required")
	}

	// Check if job already exists
	existingJob, err := j.jobRepo.GetByApplyURL(ctx, job.ApplyURL)
	if err != nil {
		return fmt.Errorf("failed to check for existing job: %w", err)
	}
	if existingJob != nil {
		return fmt.Errorf("job with this apply URL already exists")
	}

	// Set default values
	if job.Source == "" {
		job.Source = "Manual"
	}
	if job.PostedAt.IsZero() {
		job.PostedAt = time.Now()
	}

	return j.jobRepo.Create(ctx, job)
}

func (j *jobUsecase) UpdateJob(ctx context.Context, id string, updates map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	// Get existing job
	job, err := j.jobRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Apply updates (simplified version - in production you'd want proper validation)
	if title, ok := updates["title"].(string); ok && title != "" {
		job.Title = title
	}
	if companyName, ok := updates["company_name"].(string); ok && companyName != "" {
		job.CompanyName = companyName
	}
	if description, ok := updates["description"].(string); ok {
		job.Description = description
	}
	if location, ok := updates["location"].(string); ok {
		job.Location = location
	}
	if skills, ok := updates["extracted_skills"].([]string); ok {
		job.ExtractedSkills = skills
	}

	return j.jobRepo.Update(ctx, job)
}

func (j *jobUsecase) DeleteJob(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	return j.jobRepo.Delete(ctx, id)
}

func (j *jobUsecase) GetJobStats(ctx context.Context) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	// Get total job count
	filter := domain.JobFilter{
		Page:  1,
		Limit: 1,
	}
	
	_, total, err := j.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get job count: %w", err)
	}

	// Get jobs by source
	sources := j.jobAggregationSvc.GetSupportedSources()
	sourceStats := make(map[string]int64)
	
	for _, source := range sources {
		sourceFilter := domain.JobFilter{
			Source: source.Name,
			Page:   1,
			Limit:  1,
		}
		_, sourceTotal, err := j.jobRepo.List(ctx, sourceFilter)
		if err == nil {
			sourceStats[source.Name] = sourceTotal
		}
	}

	// Get recent jobs count (last 7 days)
	recentFilter := domain.JobFilter{
		Page:      1,
		Limit:     1,
		SortBy:    "posted_at",
		SortOrder: "desc",
	}
	
	recentJobs, _, err := j.jobRepo.List(ctx, recentFilter)
	recentCount := int64(0)
	if err == nil && len(recentJobs) > 0 {
		weekAgo := time.Now().AddDate(0, 0, -7)
		for _, job := range recentJobs {
			if job.PostedAt.After(weekAgo) {
				recentCount++
			}
		}
	}

	stats := map[string]interface{}{
		"total_jobs":          total,
		"jobs_by_source":      sourceStats,
		"recent_jobs_7_days":  recentCount,
		"supported_sources":   len(sources),
		"last_updated":        time.Now(),
	}

	return stats, nil
}

// Helper method to sort jobs by match score
func (j *jobUsecase) sortJobsByMatchScore(jobs []domain.Job) {
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

// SearchJobsBySkills provides a simplified skill-based search
func (j *jobUsecase) SearchJobsBySkills(ctx context.Context, skills []string, limit int) ([]domain.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	filter := domain.JobFilter{
		Skills:    skills,
		Page:      1,
		Limit:     limit,
		SortBy:    "posted_at",
		SortOrder: "desc",
	}

	jobs, _, err := j.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search jobs by skills: %w", err)
	}

	return jobs, nil
}

// GetTrendingJobs returns jobs that are trending based on various factors
func (j *jobUsecase) GetTrendingJobs(ctx context.Context, limit int) ([]domain.Job, error) {
	ctx, cancel := context.WithTimeout(ctx, j.contextTimeout)
	defer cancel()

	// For now, trending jobs are just recent jobs with good match potential
	// In the future, this could incorporate view counts, application rates, etc.
	filter := domain.JobFilter{
		Page:      1,
		Limit:     limit,
		SortBy:    "posted_at",
		SortOrder: "desc",
	}

	jobs, _, err := j.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending jobs: %w", err)
	}

	return jobs, nil
}
