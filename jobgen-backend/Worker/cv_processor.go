package Worker

import (
	"context"
	"io"
	"log"
	"net/http"

	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"
	usecases "jobgen-backend/Usecases"
)

type CVProcessor struct {
	queue       infrastructure.QueueService
	repo        domain.CVRepository
	parser      infrastructure.CVParserService
	fileUsecase domain.IFileUsecase
	aiService   domain.AIService
}

func NewCVProcessor(q infrastructure.QueueService, r domain.CVRepository, p infrastructure.CVParserService, fu domain.IFileUsecase, ai domain.AIService) *CVProcessor {
	return &CVProcessor{
		queue:       q,
		repo:        r,
		parser:      p,
		fileUsecase: fu,
		aiService:   ai,
	}
}

// Start runs the worker loop. This should be run in a separate goroutine.
func (w *CVProcessor) Start() {
	log.Println("✅ CV Processing Worker started and waiting for jobs...")
	for {
		jobID, err := w.queue.Dequeue()
		if err != nil {
			log.Printf("🔴 Error dequeuing job: %v", err)
			continue
		}
		log.Printf("🔵 Processing job ID: %s", jobID)
		w.processJob(jobID)
	}
}

func (w *CVProcessor) processJob(jobID string) {
	w.repo.UpdateStatus(jobID, domain.StatusProcessing)

	cv, err := w.repo.GetByID(jobID)
	if err != nil {
		log.Printf("🔴 Error fetching CV for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}

	// Generate a presigned URL and fetch file content
	url, err := w.fileUsecase.Download(context.Background(), cv.FileStorageID, cv.UserID)
	if err != nil {
		log.Printf("🔴 Error generating download URL for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		if err == nil {
			err = io.EOF
		}
		log.Printf("🔴 Error downloading file for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, "failed to download file")
		if resp != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()

	rawText, err := w.parser.ExtractText(resp.Body)
	if err != nil {
		log.Printf("🔴 Error parsing PDF for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}

	parsedResults, err := usecases.ParseTextToCVSections(rawText)
	if err != nil {
		log.Printf("🔴 Error structuring text for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}
	parsedResults.RawText = rawText

	suggestions, err := w.aiService.AnalyzeCV(rawText)
	if err != nil {
		log.Printf("🔴 Error from AI service for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}

	// Ensure suggestions are of the correct type
	parsedResults.Suggestions = suggestions
	parsedResults.Score = usecases.CalculateScore(parsedResults.Suggestions)

	if err := w.repo.UpdateWithResults(jobID, parsedResults); err != nil {
		log.Printf("🔴 Error saving final results for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}
	log.Printf("✅ Successfully processed job ID: %s", jobID)
}
