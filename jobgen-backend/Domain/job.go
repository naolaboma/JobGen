package domain

import (
	"context"
	"time"
)

// Job represents a single job listing document in the 'jobs' collection.
type Job struct {
	ID                     string    `json:"id" bson:"_id,omitempty"`
	Title                  string    `json:"title" bson:"title"`
	CompanyName            string    `json:"company_name" bson:"company_name"`
	Location               string    `json:"location" bson:"location"`
	Description            string    `json:"description" bson:"description"`
	FullDescriptionHTML    string    `json:"-" bson:"full_description_html"` // Store raw HTML, exclude from general API responses
	ApplyURL               string    `json:"apply_url" bson:"apply_url"`
	Source                 string    `json:"source" bson:"source"` // e.g., "RemoteOK", "Indeed"
	PostedAt               time.Time `json:"posted_at" bson:"posted_at"`
	IsSponsorshipAvailable bool      `json:"is_sponsorship_available" bson:"is_sponsorship_available"`
	ExtractedSkills        []string  `json:"extracted_skills" bson:"extracted_skills"` // Skills extracted from the description
	CreatedAt              time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" bson:"updated_at"`
	MatchScore             *float64  `json:"match_score,omitempty" bson:"-"` // Not stored in DB, calculated at runtime
	// RemoteOK specific fields
	RemoteOKID    string   `json:"remote_ok_id,omitempty" bson:"remote_ok_id,omitempty"`
	Salary        string   `json:"salary,omitempty" bson:"salary,omitempty"`
	Tags          []string `json:"tags,omitempty" bson:"tags,omitempty"`
	CompanyLogo   string   `json:"company_logo,omitempty" bson:"company_logo,omitempty"`
	OriginalData  string   `json:"-" bson:"original_data,omitempty"` // Store original JSON for reference
}

// JobFilter represents search and filter criteria for jobs
type JobFilter struct {
	Query       string   `json:"query,omitempty"`
	Skills      []string `json:"skills,omitempty"`
	Location    string   `json:"location,omitempty"`
	Sponsorship *bool    `json:"sponsorship,omitempty"`
	Source      string   `json:"source,omitempty"`
	Page        int      `json:"page"`
	Limit       int      `json:"limit"`
	SortBy      string   `json:"sort_by"`
	SortOrder   string   `json:"sort_order"`
}

// PaginatedJobsResponse represents paginated job results
type PaginatedJobsResponse struct {
	Jobs       []Job `json:"jobs"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// JobScrapeSource represents different job scraping sources
type JobScrapeSource struct {
	Name        string `json:"name"`
	BaseURL     string `json:"base_url"`
	IsActive    bool   `json:"is_active"`
	RateLimit   int    `json:"rate_limit"` // requests per minute
	LastScraped *time.Time `json:"last_scraped,omitempty"`
}

// UserJobPreferences represents user job matching preferences
type UserJobPreferences struct {
	Skills          []string `json:"skills"`
	ExperienceYears int      `json:"experience_years"`
	PreferredSalary string   `json:"preferred_salary,omitempty"`
	Locations       []string `json:"locations,omitempty"`
}

// Repository interfaces
type IJobRepository interface {
	Create(ctx context.Context, job *Job) error
	GetByID(ctx context.Context, id string) (*Job, error)
	GetByApplyURL(ctx context.Context, applyURL string) (*Job, error)
	Update(ctx context.Context, job *Job) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter JobFilter) ([]Job, int64, error)
	BulkUpsert(ctx context.Context, jobs []Job) error
	GetJobsForMatching(ctx context.Context, limit int, offset int) ([]Job, error)
}

// Scraper interfaces
type IJobScraper interface {
	GetName() string
	GetBaseURL() string
	ScrapeJobs(ctx context.Context, maxJobs int) ([]Job, error)
	GetRateLimit() int // requests per minute
}

// Job aggregation service
type IJobAggregationService interface {
	AggregateFromAllSources(ctx context.Context) error
	AggregateFromSource(ctx context.Context, sourceName string) error
	GetSupportedSources() []JobScrapeSource
}

// Job matching service
type IJobMatchingService interface {
	CalculateMatchScore(job Job, preferences UserJobPreferences) float64
	GetMatchedJobs(ctx context.Context, userID string, limit int, offset int) ([]Job, error)
	UpdateUserPreferences(ctx context.Context, userID string, preferences UserJobPreferences) error
}

// Use case interfaces
type IJobUsecase interface {
	GetJobs(ctx context.Context, filter JobFilter) (*PaginatedJobsResponse, error)
	GetJobByID(ctx context.Context, id string) (*Job, error)
	SearchJobs(ctx context.Context, userID string, filter JobFilter) (*PaginatedJobsResponse, error)
	GetMatchedJobs(ctx context.Context, userID string, limit int, offset int) (*PaginatedJobsResponse, error)
	AggregateJobs(ctx context.Context) error
	GetJobSources(ctx context.Context) ([]JobScrapeSource, error)

	// Newly added
	CreateJob(ctx context.Context, job *Job) error
	UpdateJob(ctx context.Context, id string, updates map[string]interface{}) error
	DeleteJob(ctx context.Context, id string) error
	GetJobStats(ctx context.Context) (map[string]interface{}, error)
	GetTrendingJobs(ctx context.Context, limit int) ([]Job, error)
	SearchJobsBySkills(ctx context.Context, skills []string, limit int) ([]Job, error)
}

// Skill extraction interface
type ISkillExtractor interface {
	ExtractSkills(jobDescription string) []string
	ExtractFromTitle(jobTitle string) []string
}
