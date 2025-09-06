package usecases

import (
	"context"
	"errors"
	domain "jobgen-backend/Domain"
	infrastructure "jobgen-backend/Infrastructure"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type CVUsecase interface {
	CreateParsingJob(userID string, fileHeader *multipart.FileHeader) (string, error)
	GetJobStatusAndResult(jobID string) (*domain.CV, error)
}

type cvUsecase struct {
	repo        domain.CVRepository
	queue       infrastructure.QueueService
	fileUsecase domain.IFileUsecase // use higher-level file usecase
}

func NewCVUsecase(repo domain.CVRepository, q infrastructure.QueueService, fu domain.IFileUsecase) CVUsecase {
	return &cvUsecase{repo: repo, queue: q, fileUsecase: fu}
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

	// Upload via File Usecase and store metadata
	meta := &domain.File{
		UserID:     userID,
		FileName:   fileHeader.Filename,
		BucketName: "documents",
		Size:       fileHeader.Size,
	}
	savedMeta, err := uc.fileUsecase.Upload(
		context.Background(), file, meta,
	)
	if err != nil {
		return "", err
	}

	jobID := uuid.NewString()
	cv := &domain.CV{
		ID:            jobID,
		UserID:        userID,
		FileStorageID: savedMeta.ID, // store DB ID for later download
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
