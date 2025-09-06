package scrapers

import (
	"context"
	"fmt"
	domain "jobgen-backend/Domain"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

// BaseScraper provides common functionality for Colly-based scrapers
type BaseScraper struct {
	name      string
	baseURL   string
	rateLimit int
	// collectorOptions stored so we can create fresh collectors per run
	collectorOptions []colly.CollectorOption
}

// NewBaseScraper stores options for later fresh-collector creation
func NewBaseScraper(name, baseURL string, rateLimit int) *BaseScraper {
	opts := []colly.CollectorOption{
		colly.UserAgent("JobGenBot/1.0 (+https://jobgen.io/bot)"),
		colly.Debugger(&debug.LogDebugger{}),
	}

	return &BaseScraper{
		name:             name,
		baseURL:          baseURL,
		rateLimit:        rateLimit,
		collectorOptions: opts,
	}
}

// createCollector returns a NEW collector instance configured the same way
// (this ensures a fresh in-memory visited map each run)
func (b *BaseScraper) createCollector() *colly.Collector {
	c := colly.NewCollector(b.collectorOptions...)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       3 * time.Second,
	})
	c.SetRequestTimeout(30 * time.Second)
	return c
}

// WeWorkRemotelyScraper scrapes jobs from WeWorkRemotely
type WeWorkRemotelyScraper struct {
	*BaseScraper
}

func NewWeWorkRemotelyScraper() domain.IJobScraper {
	base := NewBaseScraper("WeWorkRemotely", "https://weworkremotely.com", 10)
	return &WeWorkRemotelyScraper{BaseScraper: base}
}

func (w *WeWorkRemotelyScraper) GetName() string {
	return w.name
}

func (w *WeWorkRemotelyScraper) GetBaseURL() string {
	return w.baseURL
}

func (w *WeWorkRemotelyScraper) GetRateLimit() int {
	return w.rateLimit
}
func (w *WeWorkRemotelyScraper) ScrapeJobs(ctx context.Context, maxJobs int) ([]domain.Job, error) {
	var jobs []domain.Job
	jobCount := 0

	c := w.createCollector()  // <--- fresh collector per-run

	c.OnError(func(r *colly.Response, err error) {
		// Ignore already-visited as non-fatal (this can happen in... odd cases).
		if strings.Contains(err.Error(), "already visited") {
			fmt.Printf("%s: already visited %s\n", w.name, r.Request.URL.String())
			return
		}

		fmt.Printf("WeWorkRemotely scraping error: %s\n", err.Error())
	})

	// rest remains the same...
	c.OnHTML(".jobs li", func(e *colly.HTMLElement) {
		if maxJobs > 0 && jobCount >= maxJobs {
			return
		}
		job := w.extractJobFromElement(e)
		if job != nil {
			jobs = append(jobs, *job)
			jobCount++
		}
	})

	if err := c.Visit("https://weworkremotely.com/remote-jobs"); err != nil {
    if strings.Contains(err.Error(), "already visited") {
        fmt.Println("Visited already (ignored): https://weworkremotely.com/remote-jobs")
    } else {
        return nil, fmt.Errorf("failed to scrape WeWorkRemotely: %w", err)
    }
}


	c.Wait()
	return jobs, nil
}


func (w *WeWorkRemotelyScraper) extractJobFromElement(e *colly.HTMLElement) *domain.Job {
	// Extract job title
	title := strings.TrimSpace(e.ChildText(".title"))
	if title == "" {
		return nil
	}
	
	// Extract company name
	companyName := strings.TrimSpace(e.ChildText(".company"))
	if companyName == "" {
		companyName = "Unknown Company"
	}
	
	// Extract job URL
	jobURL := e.ChildAttr("a", "href")
	if jobURL == "" {
		return nil
	}
	
	// Build full URL
	fullURL, err := url.JoinPath("https://weworkremotely.com", jobURL)
	if err != nil {
		return nil
	}
	
	// Extract location (usually "Anywhere")
	location := strings.TrimSpace(e.ChildText(".region"))
	if location == "" {
		location = "Remote"
	}
	
	// Extract category/tags for skills
	category := strings.TrimSpace(e.ChildText(".category"))
	var skills []string
	if category != "" {
		skills = append(skills, category)
	}
	
	// Create job description placeholder (we'll need to scrape individual job pages for full description)
	description := fmt.Sprintf("Remote job at %s in %s", companyName, category)
	
	job := &domain.Job{
		Title:                  title,
		CompanyName:           companyName,
		Location:              location,
		Description:           description,
		ApplyURL:              fullURL,
		Source:                "WeWorkRemotely",
		PostedAt:              time.Now(), // WeWorkRemotely doesn't show posting date on listing
		IsSponsorshipAvailable: false,
		ExtractedSkills:       skills,
	}
	
	return job
}

// RemoteCoScraper scrapes jobs from Remote.co
type RemoteCoScraper struct {
	*BaseScraper
}

func NewRemoteCoScraper() domain.IJobScraper {
	base := NewBaseScraper("Remote.co", "https://remote.co", 8)
	return &RemoteCoScraper{BaseScraper: base}
}

func (r *RemoteCoScraper) GetName() string {
	return r.name
}

func (r *RemoteCoScraper) GetBaseURL() string {
	return r.baseURL
}

func (r *RemoteCoScraper) GetRateLimit() int {
	return r.rateLimit
}

func (r *RemoteCoScraper) ScrapeJobs(ctx context.Context, maxJobs int) ([]domain.Job, error) {
	var jobs []domain.Job
	jobCount := 0
	
	c := r.createCollector()
	
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Remote.co scraping error: %s\n", err.Error())
	})
	
	// Extract job listings from Remote.co
	c.OnHTML(".job_board_table tbody tr", func(e *colly.HTMLElement) {
		if maxJobs > 0 && jobCount >= maxJobs {
			return
		}
		
		job := r.extractJobFromRow(e)
		if job != nil {
			jobs = append(jobs, *job)
			jobCount++
		}
	})
	
	err := c.Visit("https://remote.co/remote-jobs/")
	if err != nil {
		return nil, fmt.Errorf("failed to scrape Remote.co: %w", err)
	}
	
	c.Wait()
	
	return jobs, nil
}

func (r *RemoteCoScraper) extractJobFromRow(e *colly.HTMLElement) *domain.Job {
	// Extract job title and URL
	titleElement := e.DOM.Find("td.job_title a").First()
	if titleElement.Length() == 0 {
		return nil
	}
	
	title := strings.TrimSpace(titleElement.Text())
	if title == "" {
		return nil
	}
	
	jobURL, exists := titleElement.Attr("href")
	if !exists || jobURL == "" {
		return nil
	}
	
	// Build full URL
	fullURL, err := url.JoinPath("https://remote.co", jobURL)
	if err != nil {
		return nil
	}
	
	// Extract company name
	companyName := strings.TrimSpace(e.ChildText("td.company"))
	if companyName == "" {
		companyName = "Unknown Company"
	}
	
	// Extract category
	category := strings.TrimSpace(e.ChildText("td.category"))
	var skills []string
	if category != "" {
		skills = append(skills, category)
	}
	
	// Extract location (usually remote)
	location := strings.TrimSpace(e.ChildText("td.location"))
	if location == "" {
		location = "Remote"
	}
	
	// Extract posting date
	postedAt := time.Now()
	dateStr := strings.TrimSpace(e.ChildText("td.date"))
	if dateStr != "" {
		// Try to parse date (Remote.co uses relative dates like "2 days ago")
		postedAt = r.parseRelativeDate(dateStr)
	}
	
	description := fmt.Sprintf("Remote %s position at %s", category, companyName)
	
	job := &domain.Job{
		Title:                  title,
		CompanyName:           companyName,
		Location:              location,
		Description:           description,
		ApplyURL:              fullURL,
		Source:                "Remote.co",
		PostedAt:              postedAt,
		IsSponsorshipAvailable: false,
		ExtractedSkills:       skills,
	}
	
	return job
}

func (r *RemoteCoScraper) parseRelativeDate(dateStr string) time.Time {
	now := time.Now()
	dateStr = strings.ToLower(dateStr)
	
	// Parse relative dates like "2 days ago", "1 week ago"
	if strings.Contains(dateStr, "day") {
		re := regexp.MustCompile(`(\d+)\s+day`)
		matches := re.FindStringSubmatch(dateStr)
		if len(matches) > 1 {
			days, err := strconv.Atoi(matches[1])
			if err == nil {
				return now.AddDate(0, 0, -days)
			}
		}
	} else if strings.Contains(dateStr, "week") {
		re := regexp.MustCompile(`(\d+)\s+week`)
		matches := re.FindStringSubmatch(dateStr)
		if len(matches) > 1 {
			weeks, err := strconv.Atoi(matches[1])
			if err == nil {
				return now.AddDate(0, 0, -weeks*7)
			}
		}
	} else if strings.Contains(dateStr, "month") {
		re := regexp.MustCompile(`(\d+)\s+month`)
		matches := re.FindStringSubmatch(dateStr)
		if len(matches) > 1 {
			months, err := strconv.Atoi(matches[1])
			if err == nil {
				return now.AddDate(0, -months, 0)
			}
		}
	}
	
	return now
}

// NoDesk RemoteScraper scrapes jobs from NoDesk
type NoDeskScraper struct {
	*BaseScraper
}

func NewNoDeskScraper() domain.IJobScraper {
	base := NewBaseScraper("NoDesk", "https://nodesk.co", 6)
	return &NoDeskScraper{BaseScraper: base}
}

func (n *NoDeskScraper) GetName() string {
	return n.name
}

func (n *NoDeskScraper) GetBaseURL() string {
	return n.baseURL
}

func (n *NoDeskScraper) GetRateLimit() int {
	return n.rateLimit
}

func (n *NoDeskScraper) ScrapeJobs(ctx context.Context, maxJobs int) ([]domain.Job, error) {
	var jobs []domain.Job
	jobCount := 0
	
	c := n.createCollector()
	
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("NoDesk scraping error: %s\n", err.Error())
	})
	
	// Extract job listings from NoDesk
	c.OnHTML(".job-board-item", func(e *colly.HTMLElement) {
		if maxJobs > 0 && jobCount >= maxJobs {
			return
		}
		
		job := n.extractJobFromItem(e)
		if job != nil {
			jobs = append(jobs, *job)
			jobCount++
		}
	})
	
	err := c.Visit("https://nodesk.co/remote-jobs/")
	if err != nil {
		return nil, fmt.Errorf("failed to scrape NoDesk: %w", err)
	}
	
	c.Wait()
	
	return jobs, nil
}

func (n *NoDeskScraper) extractJobFromItem(e *colly.HTMLElement) *domain.Job {
	// Extract job title
	title := strings.TrimSpace(e.ChildText(".job-title"))
	if title == "" {
		return nil
	}
	
	// Extract company name
	companyName := strings.TrimSpace(e.ChildText(".company-name"))
	if companyName == "" {
		companyName = "Unknown Company"
	}
	
	// Extract job URL
	jobURL := e.ChildAttr("a", "href")
	if jobURL == "" {
		return nil
	}
	
	// Build full URL if relative
	var fullURL string
	if strings.HasPrefix(jobURL, "http") {
		fullURL = jobURL
	} else {
		var err error
		fullURL, err = url.JoinPath("https://nodesk.co", jobURL)
		if err != nil {
			return nil
		}
	}
	
	// Extract category/tags
	category := strings.TrimSpace(e.ChildText(".job-category"))
	var skills []string
	if category != "" {
		skills = append(skills, category)
	}
	
	// Extract salary if available
	salary := strings.TrimSpace(e.ChildText(".salary"))
	
	description := fmt.Sprintf("Remote %s position at %s", category, companyName)
	if salary != "" {
		description += fmt.Sprintf(" - %s", salary)
	}
	
	job := &domain.Job{
		Title:                  title,
		CompanyName:           companyName,
		Location:              "Remote",
		Description:           description,
		ApplyURL:              fullURL,
		Source:                "NoDesk",
		PostedAt:              time.Now(),
		IsSponsorshipAvailable: false,
		ExtractedSkills:       skills,
		Salary:                salary,
	}
	
	return job
}

// Utility function to clean and extract skills from job descriptions
func ExtractSkillsFromDescription(description, title string) []string {
	var skills []string
	
	// Common tech skills to look for
	techSkills := []string{
		"javascript", "typescript", "python", "java", "go", "golang", "rust", "c++", "c#", "csharp",
		"react", "vue", "angular", "node.js", "nodejs", "django", "flask", "spring", "laravel",
		"mysql", "postgresql", "mongodb", "redis", "elasticsearch",
		"aws", "azure", "gcp", "docker", "kubernetes", "jenkins", "git", "github", "gitlab",
		"html", "css", "sass", "scss", "bootstrap", "tailwind",
		"api", "rest", "graphql", "microservices", "devops", "ci/cd",
		"machine learning", "ai", "data science", "analytics", "sql",
	}
	
	combinedText := strings.ToLower(description + " " + title)
	
	for _, skill := range techSkills {
		if strings.Contains(combinedText, strings.ToLower(skill)) {
			skills = append(skills, skill)
		}
	}
	
	return skills
}
