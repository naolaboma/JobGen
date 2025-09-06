package usecases

import (
	"errors"
	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CVUsecase interface {
	CreateParsingJob(userID string, fileHeader *multipart.FileHeader) (string, error)
	CreateParsingJobFromFileID(userID string, fileID string) (string, error)
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

func (uc *cvUsecase) GetJobStatusAndResult(jobID string) (*domain.CV, error) {
	return uc.repo.GetByID(jobID)
}

func (uc *cvUsecase) CreateParsingJobFromFileID(userID string, fileID string) (string, error) {
	if strings.TrimSpace(fileID) == "" {
		return "", errors.New("fileId is required")
	}
	jobID := uuid.NewString()
	cv := &domain.CV{
		ID:            jobID,
		UserID:        userID,
		FileStorageID: fileID,
		FileName:      "",
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
