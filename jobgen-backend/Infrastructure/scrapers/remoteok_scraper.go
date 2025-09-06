package scrapers

import (
	"context"
	"encoding/json"
	"fmt"
	domain "jobgen-backend/Domain"
	"strings"
	"time"

	"github.com/go-api-libs/remote-ok-jobs/pkg/remoteokjobs"
)

type RemoteOKScraper struct {
	client    *remoteokjobs.Client
	rateLimit int
}

func NewRemoteOKScraper() (domain.IJobScraper, error) {
	client, err := remoteokjobs.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create RemoteOK client: %w", err)
	}

	return &RemoteOKScraper{
		client:    client,
		rateLimit: 10, // 10 requests per minute
	}, nil
}

func (r *RemoteOKScraper) GetName() string {
	return "RemoteOK"
}

func (r *RemoteOKScraper) GetBaseURL() string {
	return "https://remoteok.io"
}

func (r *RemoteOKScraper) GetRateLimit() int {
	return r.rateLimit
}

func (r *RemoteOKScraper) ScrapeJobs(ctx context.Context, maxJobs int) ([]domain.Job, error) {
	remoteJobs, err := r.client.GetJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch jobs from RemoteOK: %w", err)
	}

	var jobs []domain.Job
	limit := len(remoteJobs)
	if maxJobs > 0 && maxJobs < limit {
		limit = maxJobs
	}

	for i := 0; i < limit; i++ {
		remoteJob := remoteJobs[i]

		job, err := r.convertRemoteOKJob(remoteJob)
		if err != nil {
			fmt.Printf("Error converting RemoteOK job: %v\n", err)
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (r *RemoteOKScraper) convertRemoteOKJob(remoteJob remoteokjobs.Job) (domain.Job, error) {
	// Posted date
	postedAt := time.Now()
	if remoteJob.Date != nil {
		postedAt = *remoteJob.Date
	}

	// Extract skills/tags
	var skills []string
	var tags []string
	if len(remoteJob.Tags) > 0 {
		for _, tag := range remoteJob.Tags {
			tagStr := strings.ToLower(strings.TrimSpace(tag))
			if tagStr != "" {
				skills = append(skills, tagStr)
				tags = append(tags, tagStr)
			}
		}
	}

	// Description
	description := ""
	if remoteJob.Description != nil {
		description = *remoteJob.Description
	}

	// Apply URL
	applyURL := fmt.Sprintf("https://remoteok.io/remote-jobs/%s", safeDeref(remoteJob.ID))
	if remoteJob.URL != nil {
		applyURL = remoteJob.URL.String()
	}

	// Company info
	companyName := "Unknown Company"
	if remoteJob.Company != nil && *remoteJob.Company != "" {
		companyName = *remoteJob.Company
	}

	// Location
	location := "Remote"
	if remoteJob.Location != nil && *remoteJob.Location != "" {
		location = *remoteJob.Location
	}

	// Salary
	salary := ""
	if remoteJob.SalaryMin != nil && remoteJob.SalaryMax != nil {
		salary = fmt.Sprintf("$%d - $%d", *remoteJob.SalaryMin, *remoteJob.SalaryMax)
	} else if remoteJob.SalaryMin != nil {
		salary = fmt.Sprintf("From $%d", *remoteJob.SalaryMin)
	} else if remoteJob.SalaryMax != nil {
		salary = fmt.Sprintf("Up to $%d", *remoteJob.SalaryMax)
	}

	// Company logo
	companyLogo := ""
	if remoteJob.CompanyLogo != nil {
		companyLogo = remoteJob.CompanyLogo.String()
	}

	// Store original data
	originalData, _ := json.Marshal(remoteJob)

	job := domain.Job{
		Title:                 safeDeref(remoteJob.Position),
		CompanyName:           companyName,
		Location:              location,
		Description:           description,
		FullDescriptionHTML:   description,
		ApplyURL:              applyURL,
		Source:                "RemoteOK",
		PostedAt:              postedAt,
		IsSponsorshipAvailable: false,
		ExtractedSkills:       skills,
		RemoteOKID:            safeDeref(remoteJob.ID),
		Salary:                salary,
		Tags:                  tags,
		CompanyLogo:           companyLogo,
		OriginalData:          string(originalData),
	}

	return job, nil
}


// Helpers
func safeDeref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// --- Filtering helpers ---
func (r *RemoteOKScraper) GetJobsByTag(ctx context.Context, tag string, maxJobs int) ([]domain.Job, error) {
	allJobs, err := r.ScrapeJobs(ctx, 0)
	if err != nil {
		return nil, err
	}

	var filteredJobs []domain.Job
	tagLower := strings.ToLower(tag)

	for _, job := range allJobs {
		for _, jobTag := range job.Tags {
			if strings.ToLower(jobTag) == tagLower {
				filteredJobs = append(filteredJobs, job)
				break
			}
		}
		if maxJobs > 0 && len(filteredJobs) >= maxJobs {
			break
		}
	}

	return filteredJobs, nil
}

func (r *RemoteOKScraper) GetRecentJobs(ctx context.Context, days int, maxJobs int) ([]domain.Job, error) {
	allJobs, err := r.ScrapeJobs(ctx, 0)
	if err != nil {
		return nil, err
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	var recentJobs []domain.Job

	for _, job := range allJobs {
		if job.PostedAt.After(cutoffDate) {
			recentJobs = append(recentJobs, job)
			if maxJobs > 0 && len(recentJobs) >= maxJobs {
				break
			}
		}
	}

	return recentJobs, nil
}
