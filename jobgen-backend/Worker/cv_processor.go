package Worker

import (
	"fmt"
	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"
	usecases "jobgen-backend/Usecases"
	"log"
)

type CVProcessor struct {
	queue     infrastructure.QueueService
	repo      domain.CVRepository
	parser    infrastructure.CVParserService
	fileStore infrastructure.FileStorageService
	aiService domain.AIService
}

func NewCVProcessor(q infrastructure.QueueService, r domain.CVRepository, p infrastructure.CVParserService, fs infrastructure.FileStorageService, ai domain.AIService) *CVProcessor {
	return &CVProcessor{
		queue:     q,
		repo:      r,
		parser:    p,
		fileStore: fs,
		aiService: ai,
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

	file, err := w.fileStore.GetFile(cv.FileStorageID)
	if err != nil {
		log.Printf("🔴 Error getting file from storage for job %s: %v", jobID, err)
		w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
		return
	}
	defer file.Close()

	rawText, err := w.parser.ExtractText(file)
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

	// Heuristic low-confidence: if no sections detected, flag for manual review
	if len(parsedResults.Skills) == 0 && len(parsedResults.Experiences) == 0 && len(parsedResults.Educations) == 0 {
		parsedResults.ProcessingError = "low_confidence_parse"
	}

	// Try to get AI suggestions; if it fails (e.g., invalid/absent API key), continue without failing the job
	suggestions, aiErr := w.aiService.AnalyzeCV(rawText)
	if aiErr != nil {
		log.Printf("🟠 AI unavailable for job %s: %v", jobID, aiErr)
		if parsedResults.ProcessingError != "" {
			parsedResults.ProcessingError = fmt.Sprintf("%s; ai_unavailable", parsedResults.ProcessingError)
		} else {
			parsedResults.ProcessingError = "ai_unavailable"
		}
		parsedResults.Suggestions = nil
		parsedResults.Score = usecases.CalculateScore(parsedResults.Suggestions)

		if err := w.repo.UpdateWithResults(jobID, parsedResults); err != nil {
			log.Printf("🔴 Error saving results (no AI) for job %s: %v", jobID, err)
			w.repo.UpdateStatus(jobID, domain.StatusFailed, err.Error())
			return
		}
		w.repo.UpdateStatus(jobID, domain.StatusCompleted)
		log.Printf("✅ Processed job %s without AI suggestions", jobID)
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
	w.repo.UpdateStatus(jobID, domain.StatusCompleted)
	log.Printf("✅ Successfully processed job ID: %s", jobID)
}
