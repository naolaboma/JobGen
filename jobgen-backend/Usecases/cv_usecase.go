package usecases

import (
	"errors"
	"jobgen-backend/Domain"
	"jobgen-backend/Infrastructure"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CVUsecase interface {
	CreateParsingJob(userID string, fileHeader *multipart.FileHeader) (string, error)
	GetJobStatusAndResult(jobID string) (*domain.CV, error)
}

type cvUsecase struct {
	repo      domain.CVRepository
	queue     infrastructure.QueueService
	fileStore domain.FileStorageService // From file_management.go
}

func NewCVUsecase(repo domain.CVRepository, q infrastructure.QueueService, fs domain.FileStorageService) CVUsecase {
	return &cvUsecase{repo: repo, queue: q, fileStore: fs}
}

func (uc *cvUsecase) CreateParsingJob(userID string, fileHeader *multipart.FileHeader) (string, error) {
	// Input validation
	if fileHeader.Size > 5*1024*1024 { // 5 MB limit
		return "", errors.New("file size exceeds the 5MB limit")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Integrate with File Storage Service
	fileID, err := uc.fileStore.UploadFile(userID, "CV", fileHeader.Filename, file)
	if err != nil {
		return "", err
	}

	jobID := uuid.NewString()
	cv := &domain.CV{
		ID:            jobID,
		UserID:        userID,
		FileStorageID: fileID,
		FileName:      fileHeader.Filename,
		Status:        domain.StatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	if err := uc.repo.Create(cv); err != nil {
		return "", err
	}

	if err := uc.queue.Enqueue(jobID); err != nil {
		return "", err
	}

	return jobID, nil
}

func (uc *cvUsecase) GetJobStatusAndResult(jobID string) (*omain.CV, error) {
	return uc.repo.GetByID(jobID)
}
