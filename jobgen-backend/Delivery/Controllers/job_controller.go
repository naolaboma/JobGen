package controllers

import (
	"net/http"
	"strconv"
	"strings"

	domain "jobgen-backend/Domain"

	"github.com/gin-gonic/gin"
)

type JobController struct {
	jobUsecase domain.IJobUsecase
}

func NewJobController(jobUsecase domain.IJobUsecase) *JobController {
	return &JobController{
		jobUsecase: jobUsecase,
	}
}

// SearchJobsRequest represents the search request body
type SearchJobsRequest struct {
	Query       string   `json:"query,omitempty"`
	Skills      []string `json:"skills,omitempty"`
	Location    string   `json:"location,omitempty"`
	Sponsorship *bool    `json:"sponsorship,omitempty"`
	Source      string   `json:"source,omitempty"`
}

// @Summary Get all jobs with filtering and pagination
// @Description Retrieve jobs with optional filtering by query, skills, location, etc.
// @Tags Jobs
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Param query query string false "Search query for title, company, or description"
// @Param skills query string false "Comma-separated list of skills"
// @Param location query string false "Location filter"
// @Param sponsorship query bool false "Filter by sponsorship availability"
// @Param source query string false "Filter by job source"
// @Param sort_by query string false "Sort field" default(posted_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} StandardResponse "List of jobs"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs [get]
func (c *JobController) GetJobs(ctx *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Parse skills from comma-separated string
	var skills []string
	if skillsStr := ctx.Query("skills"); skillsStr != "" {
		skills = strings.Split(skillsStr, ",")
		// Trim whitespace from each skill
		for i, skill := range skills {
			skills[i] = strings.TrimSpace(skill)
		}
	}

	// Parse sponsorship parameter
	var sponsorship *bool
	if sponsorshipStr := ctx.Query("sponsorship"); sponsorshipStr != "" {
		if val, err := strconv.ParseBool(sponsorshipStr); err == nil {
			sponsorship = &val
		}
	}

	filter := domain.JobFilter{
		Query:       ctx.Query("query"),
		Skills:      skills,
		Location:    ctx.Query("location"),
		Sponsorship: sponsorship,
		Source:      ctx.Query("source"),
		Page:        page,
		Limit:       limit,
		SortBy:      ctx.DefaultQuery("sort_by", "posted_at"),
		SortOrder:   ctx.DefaultQuery("sort_order", "desc"),
	}

	result, err := c.jobUsecase.GetJobs(ctx, filter)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to retrieve jobs")
		return
	}

	paginatedData := &PaginatedResponse{
		Items:      result.Jobs,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}

	PaginatedSuccessResponse(ctx, http.StatusOK, "Jobs retrieved successfully", paginatedData)
}

// @Summary Get a specific job by ID
// @Description Retrieve detailed information about a specific job
// @Tags Jobs
// @Accept json
// @Produce json
// @Param id path string true "Job ID"
// @Success 200 {object} StandardResponse "Job details"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 404 {object} StandardResponse "Job not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/{id} [get]
func (c *JobController) GetJobByID(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if jobID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Job ID is required", nil)
		return
	}

	job, err := c.jobUsecase.GetJobByID(ctx, jobID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			NotFoundResponse(ctx, "Job not found")
		} else {
			InternalErrorResponse(ctx, "Failed to retrieve job")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Job retrieved successfully", gin.H{
		"job": job,
	})
}

// @Summary Search jobs with user context
// @Description Search and filter jobs with personalized matching if user is authenticated
// @Tags Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Param query query string false "Search query for title, company, or description"
// @Param skills query string false "Comma-separated list of skills"
// @Param location query string false "Location filter"
// @Param sponsorship query bool false "Filter by sponsorship availability"
// @Param source query string false "Filter by job source"
// @Param sort_by query string false "Sort field" default(posted_at)
// @Param sort_order query string false "Sort order" Enums(asc, desc) default(desc)
// @Success 200 {object} StandardResponse "Personalized job search results"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/search [get]
func (c *JobController) SearchJobs(ctx *gin.Context) {
	// Get user ID from context (might be empty for anonymous users)
	userID := ctx.GetString("user_id")

	// Parse query parameters (same as GetJobs)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	var skills []string
	if skillsStr := ctx.Query("skills"); skillsStr != "" {
		skills = strings.Split(skillsStr, ",")
		for i, skill := range skills {
			skills[i] = strings.TrimSpace(skill)
		}
	}

	var sponsorship *bool
	if sponsorshipStr := ctx.Query("sponsorship"); sponsorshipStr != "" {
		if val, err := strconv.ParseBool(sponsorshipStr); err == nil {
			sponsorship = &val
		}
	}

	filter := domain.JobFilter{
		Query:       ctx.Query("query"),
		Skills:      skills,
		Location:    ctx.Query("location"),
		Sponsorship: sponsorship,
		Source:      ctx.Query("source"),
		Page:        page,
		Limit:       limit,
		SortBy:      ctx.DefaultQuery("sort_by", "posted_at"),
		SortOrder:   ctx.DefaultQuery("sort_order", "desc"),
	}

	result, err := c.jobUsecase.SearchJobs(ctx, userID, filter)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to search jobs")
		return
	}

	paginatedData := &PaginatedResponse{
		Items:      result.Jobs,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}

	message := "Job search completed successfully"
	if userID != "" {
		message = "Personalized job search completed successfully"
	}

	PaginatedSuccessResponse(ctx, http.StatusOK, message, paginatedData)
}

// @Summary Get matched jobs for authenticated user
// @Description Get personalized job recommendations based on user profile
// @Tags Jobs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page (max 100)" default(10)
// @Success 200 {object} StandardResponse "Matched jobs for user"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/matched [get]
func (c *JobController) GetMatchedJobs(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		UnauthorizedResponse(ctx, "User authentication required")
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	// Calculate offset
	offset := (page - 1) * limit

	result, err := c.jobUsecase.GetMatchedJobs(ctx, userID, limit, offset)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to get matched jobs")
		return
	}

	paginatedData := &PaginatedResponse{
		Items:      result.Jobs,
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: result.TotalPages,
		HasNext:    result.HasNext,
		HasPrev:    result.HasPrev,
	}

	PaginatedSuccessResponse(ctx, http.StatusOK, "Matched jobs retrieved successfully", paginatedData)
}

// @Summary Get trending jobs
// @Description Get currently trending job listings
// @Tags Jobs
// @Accept json
// @Produce json
// @Param limit query int false "Number of jobs to return (max 50)" default(20)
// @Success 200 {object} StandardResponse "Trending jobs"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/trending [get]
func (c *JobController) GetTrendingJobs(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if limit > 50 {
		limit = 50
	}

	jobs, err := c.jobUsecase.GetTrendingJobs(ctx, limit)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to get trending jobs")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Trending jobs retrieved successfully", gin.H{
		"jobs": jobs,
		"count": len(jobs),
	})
}

// @Summary Get job statistics
// @Description Get statistics about jobs in the system
// @Tags Jobs
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse "Job statistics"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/stats [get]
func (c *JobController) GetJobStats(ctx *gin.Context) {
	stats, err := c.jobUsecase.GetJobStats(ctx)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to get job statistics")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Job statistics retrieved successfully", stats)
}

// @Summary Get supported job sources
// @Description Get list of all supported job scraping sources
// @Tags Jobs
// @Accept json
// @Produce json
// @Success 200 {object} StandardResponse "List of job sources"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/sources [get]
func (c *JobController) GetJobSources(ctx *gin.Context) {
	sources, err := c.jobUsecase.GetJobSources(ctx)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to get job sources")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Job sources retrieved successfully", gin.H{
		"sources": sources,
		"count":   len(sources),
	})
}

// @Summary Search jobs by skills
// @Description Quick search for jobs based on specific skills
// @Tags Jobs
// @Accept json
// @Produce json
// @Param skills query string true "Comma-separated list of skills"
// @Param limit query int false "Number of jobs to return (max 50)" default(20)
// @Success 200 {object} StandardResponse "Jobs matching skills"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /jobs/search-by-skills [get]
func (c *JobController) SearchJobsBySkills(ctx *gin.Context) {
	skillsStr := ctx.Query("skills")
	if skillsStr == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Skills parameter is required", nil)
		return
	}

	skills := strings.Split(skillsStr, ",")
	for i, skill := range skills {
		skills[i] = strings.TrimSpace(skill)
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if limit > 50 {
		limit = 50
	}

	jobs, err := c.jobUsecase.SearchJobsBySkills(ctx, skills, limit)
	if err != nil {
		InternalErrorResponse(ctx, "Failed to search jobs by skills")
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Jobs found for specified skills", gin.H{
		"jobs":   jobs,
		"count":  len(jobs),
		"skills": skills,
	})
}

// Admin endpoints

// @Summary Trigger job aggregation
// @Description Manually trigger job aggregation from all sources (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StandardResponse "Job aggregation started"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /admin/jobs/aggregate [post]
func (c *JobController) TriggerJobAggregation(ctx *gin.Context) {
	// This will be handled by admin middleware for authorization
	go func() {
		if err := c.jobUsecase.AggregateJobs(ctx); err != nil {
			// Log error in production
			println("Job aggregation failed:", err.Error())
		}
	}()

	SuccessResponse(ctx, http.StatusOK, "Job aggregation started in background", gin.H{
		"status": "started",
		"note":   "Aggregation is running in background and may take several minutes",
	})
}

// @Summary Create a new job
// @Description Create a new job listing (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateJobRequest true "Job details"
// @Success 201 {object} StandardResponse "Job created successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /admin/jobs [post]
func (c *JobController) CreateJob(ctx *gin.Context) {
	var req CreateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	job := &domain.Job{
		Title:                  req.Title,
		CompanyName:           req.CompanyName,
		Location:              req.Location,
		Description:           req.Description,
		FullDescriptionHTML:   req.FullDescriptionHTML,
		ApplyURL:              req.ApplyURL,
		Source:                req.Source,
		IsSponsorshipAvailable: req.IsSponsorshipAvailable,
		ExtractedSkills:       req.ExtractedSkills,
		Salary:                req.Salary,
		Tags:                  req.Tags,
	}

	if err := c.jobUsecase.CreateJob(ctx, job); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			ConflictResponse(ctx, "Job with this apply URL already exists")
		} else {
			InternalErrorResponse(ctx, "Failed to create job")
		}
		return
	}

	SuccessResponse(ctx, http.StatusCreated, "Job created successfully", gin.H{
		"job": job,
	})
}

// @Summary Update a job
// @Description Update an existing job listing (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Param request body UpdateJobRequest true "Job update details"
// @Success 200 {object} StandardResponse "Job updated successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Failure 404 {object} StandardResponse "Job not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /admin/jobs/{id} [put]
func (c *JobController) UpdateJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if jobID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Job ID is required", nil)
		return
	}

	var req UpdateJobRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ValidationErrorResponse(ctx, err)
		return
	}

	// Convert request to map for updates
	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.CompanyName != nil {
		updates["company_name"] = *req.CompanyName
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ExtractedSkills != nil {
		updates["extracted_skills"] = *req.ExtractedSkills
	}

	if err := c.jobUsecase.UpdateJob(ctx, jobID, updates); err != nil {
		if strings.Contains(err.Error(), "not found") {
			NotFoundResponse(ctx, "Job not found")
		} else {
			InternalErrorResponse(ctx, "Failed to update job")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Job updated successfully", nil)
}

// @Summary Delete a job
// @Description Delete a job listing (Admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Job ID"
// @Success 200 {object} StandardResponse "Job deleted successfully"
// @Failure 400 {object} StandardResponse "Bad request"
// @Failure 401 {object} StandardResponse "Unauthorized"
// @Failure 403 {object} StandardResponse "Forbidden"
// @Failure 404 {object} StandardResponse "Job not found"
// @Failure 500 {object} StandardResponse "Internal server error"
// @Router /admin/jobs/{id} [delete]
func (c *JobController) DeleteJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if jobID == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "VALIDATION_ERROR", "Job ID is required", nil)
		return
	}

	if err := c.jobUsecase.DeleteJob(ctx, jobID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			NotFoundResponse(ctx, "Job not found")
		} else {
			InternalErrorResponse(ctx, "Failed to delete job")
		}
		return
	}

	SuccessResponse(ctx, http.StatusOK, "Job deleted successfully", nil)
}

// Request structs for Swagger documentation
type CreateJobRequest struct {
	Title                  string   `json:"title" binding:"required"`
	CompanyName            string   `json:"company_name" binding:"required"`
	Location               string   `json:"location" binding:"required"`
	Description            string   `json:"description" binding:"required"`
	FullDescriptionHTML    string   `json:"full_description_html,omitempty"`
	ApplyURL               string   `json:"apply_url" binding:"required,url"`
	Source                 string   `json:"source,omitempty"`
	IsSponsorshipAvailable bool     `json:"is_sponsorship_available,omitempty"`
	ExtractedSkills        []string `json:"extracted_skills,omitempty"`
	Salary                 string   `json:"salary,omitempty"`
	Tags                   []string `json:"tags,omitempty"`
}

type UpdateJobRequest struct {
	Title           *string   `json:"title,omitempty"`
	CompanyName     *string   `json:"company_name,omitempty"`
	Location        *string   `json:"location,omitempty"`
	Description     *string   `json:"description,omitempty"`
	ExtractedSkills *[]string `json:"extracted_skills,omitempty"`
}
