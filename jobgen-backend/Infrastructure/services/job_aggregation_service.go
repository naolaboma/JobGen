// jobgen-backend/Infrastructure/services/job_aggregation_service.go
package services

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"jobgen-backend/Infrastructure/scrapers"
	"sync"
	"time"
)

type JobAggregationService struct {
	jobRepo  domain.IJobRepository
	scrapers map[string]domain.IJobScraper
	mu       sync.RWMutex
}

func NewJobAggregationService(jobRepo domain.IJobRepository) domain.IJobAggregationService {
	service := &JobAggregationService{
		jobRepo:  jobRepo,
		scrapers: make(map[string]domain.IJobScraper),
	}
	
	// Initialize scrapers
	service.initializeScrapers()
	
	return service
}

func (j *JobAggregationService) initializeScrapers() {
	// Initialize RemoteOK scraper
	remoteOKScraper, err := scrapers.NewRemoteOKScraper()
	if err != nil {
		fmt.Printf("Failed to initialize RemoteOK scraper: %v\n", err)
	} else {
		j.scrapers["RemoteOK"] = remoteOKScraper
	}
	
	// Initialize Colly-based scrapers
	j.scrapers["WeWorkRemotely"] = scrapers.NewWeWorkRemotelyScraper()
	j.scrapers["Remote.co"] = scrapers.NewRemoteCoScraper()
	j.scrapers["NoDesk"] = scrapers.NewNoDeskScraper()
}

func (j *JobAggregationService) AggregateFromAllSources(ctx context.Context) error {
	j.mu.RLock()
	scrapers := make(map[string]domain.IJobScraper)
	for name, scraper := range j.scrapers {
		scrapers[name] = scraper
	}
	j.mu.RUnlock()
	
	var wg sync.WaitGroup
	errors := make(chan error, len(scrapers))
	
	for name, scraper := range scrapers {
		wg.Add(1)
		go func(name string, scraper domain.IJobScraper) {
			defer wg.Done()
			
			fmt.Printf("Starting aggregation from %s\n", name)
			
			if err := j.aggregateFromScraper(ctx, scraper); err != nil {
				errors <- fmt.Errorf("failed to aggregate from %s: %w", name, err)
				return
			}
			
			fmt.Printf("Successfully aggregated jobs from %s\n", name)
		}(name, scraper)
		
		// Add delay between starting scrapers to be respectful
		time.Sleep(2 * time.Second)
	}
	
	wg.Wait()
	close(errors)
	
	// Collect any errors
	var aggregatedErrors []error
	for err := range errors {
		aggregatedErrors = append(aggregatedErrors, err)
	}
	
	if len(aggregatedErrors) > 0 {
		return fmt.Errorf("aggregation completed with %d errors: %v", len(aggregatedErrors), aggregatedErrors)
	}
	
	return nil
}

func (j *JobAggregationService) AggregateFromSource(ctx context.Context, sourceName string) error {
	j.mu.RLock()
	scraper, exists := j.scrapers[sourceName]
	j.mu.RUnlock()
	
	if !exists {
		return fmt.Errorf("scraper for source %s not found", sourceName)
	}
	
	return j.aggregateFromScraper(ctx, scraper)
}

func (j *JobAggregationService) aggregateFromScraper(ctx context.Context, scraper domain.IJobScraper) error {
	// Create a context with timeout for scraping
	scrapingCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()
	
	// Scrape jobs with a reasonable limit
	jobs, err := scraper.ScrapeJobs(scrapingCtx, 100) // Limit to 100 jobs per source
	if err != nil {
		return fmt.Errorf("failed to scrape jobs: %w", err)
	}
	
	if len(jobs) == 0 {
		fmt.Printf("No jobs found from %s\n", scraper.GetName())
		return nil
	}
	
	// Enhance jobs with skill extraction
	for i := range jobs {
		jobs[i].ExtractedSkills = j.enhanceSkills(jobs[i])
	}
	
	// Bulk upsert jobs to database
	if err := j.jobRepo.BulkUpsert(ctx, jobs); err != nil {
		return fmt.Errorf("failed to bulk upsert jobs: %w", err)
	}
	
	fmt.Printf("Successfully processed %d jobs from %s\n", len(jobs), scraper.GetName())
	
	return nil
}

func (j *JobAggregationService) enhanceSkills(job domain.Job) []string {
	// Combine existing skills with extracted ones from description and title
	skillsMap := make(map[string]bool)
	
	// Add existing skills
	for _, skill := range job.ExtractedSkills {
		skillsMap[skill] = true
	}
	
	// Add tags if available
	for _, tag := range job.Tags {
		skillsMap[tag] = true
	}
	
	// Extract additional skills from description and title using our utility function
	extractedSkills := scrapers.ExtractSkillsFromDescription(job.Description, job.Title)
	for _, skill := range extractedSkills {
		skillsMap[skill] = true
	}
	
	// Convert back to slice
	var enhancedSkills []string
	for skill := range skillsMap {
		enhancedSkills = append(enhancedSkills, skill)
	}
	
	return enhancedSkills
}

func (j *JobAggregationService) GetSupportedSources() []domain.JobScrapeSource {
	j.mu.RLock()
	defer j.mu.RUnlock()
	
	var sources []domain.JobScrapeSource
	
	for name, scraper := range j.scrapers {
		source := domain.JobScrapeSource{
			Name:      name,
			BaseURL:   scraper.GetBaseURL(),
			IsActive:  true,
			RateLimit: scraper.GetRateLimit(),
		}
		sources = append(sources, source)
	}
	
	return sources
}

// AddScraper allows adding new scrapers dynamically
func (j *JobAggregationService) AddScraper(name string, scraper domain.IJobScraper) {
	j.mu.Lock()
	defer j.mu.Unlock()
	
	j.scrapers[name] = scraper
}

// RemoveScraper allows removing scrapers
func (j *JobAggregationService) RemoveScraper(name string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	
	delete(j.scrapers, name)
}

// GetScraperStatus returns status information for all scrapers
func (j *JobAggregationService) GetScraperStatus() map[string]interface{} {
	j.mu.RLock()
	defer j.mu.RUnlock()
	
	status := make(map[string]interface{})
	
	for name, scraper := range j.scrapers {
		status[name] = map[string]interface{}{
			"name":       scraper.GetName(),
			"base_url":   scraper.GetBaseURL(),
			"rate_limit": scraper.GetRateLimit(),
			"active":     true,
		}
	}
	
	return status
}
