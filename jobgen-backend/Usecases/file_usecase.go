package usecases

import (
	"context"
	"fmt"
	"io"
	domain "jobgen-backend/Domain"
	"time"
)

type fileUsecase struct {
	fileRepo domain.IFileRepository
	s3s      domain.IFileService
}

func NewFileUsecase(fileRepo domain.IFileRepository, s3s domain.IFileService) domain.IFileUsecase {
	return &fileUsecase{fileRepo: fileRepo,
		s3s: s3s}
}

func generateKeyName(file *domain.File) string {

	// If it's a profile picture, just use userID
	// this will override the old profile picture
	if file.BucketName == "profile-pictures" {
		return file.UserID
	}

	// If it's a document, use userID as a folder
	if file.BucketName == "documents" {
		if file.UserID != "" {
			return fmt.Sprintf("%s/%s-%s", file.UserID, file.UniqueID, file.FileName)
		}
		return fmt.Sprintf("%s-%s", file.UniqueID, file.FileName)
	}

	// fallback
	return fmt.Sprintf("%s-%s", file.UniqueID, file.FileName)
}

// Delete implements domain.IFileUsecase.
func (f *fileUsecase) Delete(ctx context.Context, ID, userID string) error {
	file, err := f.fileRepo.FindByID(ctx, ID)
	if err != nil {
		return domain.ErrFileNotFound
	}

	// checks for authorization
	if userID != file.UserID {
		return domain.ErrUnauthorized
	}
	err = f.fileRepo.Delete(ctx, ID)
	if err != nil {
		return fmt.Errorf("failed to delete file from database: %w", domain.ErrInternal)
	}
	key := generateKeyName(file)
	err = f.s3s.Delete(ctx, file.BucketName, key)
	if err != nil {
		return fmt.Errorf("failed to delete file from provider: %w", domain.ErrInternal)
	}
	return nil
}

// Download implements domain.IFileUsecase.
func (f *fileUsecase) Download(ctx context.Context, fileID string, userID string) (string, error) {
	// get the meta data
	file, err := f.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return "", err
	}

	// if it is not a profile-pictures then everything is treated private
	if file.BucketName != "profile-pictures" {
		if userID != file.UserID {
			return "", domain.ErrUnauthorized
		}
	}

	// generate unique name
	key := generateKeyName(file)
	url, err := f.s3s.PresignedURL(ctx, file.BucketName, key)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", domain.ErrInternal)
	}
	return url, err
}

// GetMetaData gets the meta data stored on db about a particular file
func (f *fileUsecase) GetMetaData(ctx context.Context, fileID string) (domain.File, error) {
	data, err := f.fileRepo.FindByID(ctx, fileID)
	if err != nil {
		return domain.File{}, err
	}
	return *data, nil
}

// Exists implements domain.IFileUsecase.
// Exists check if a file with id exist or not
func (f *fileUsecase) Exists(ctx context.Context, ID string) (bool, error) {
	fetchedData, err := f.fileRepo.FindByID(ctx, ID)
	if err != nil {
		return false, err
	}
	return fetchedData != nil, nil
}

// Upload implements domain.IFileUsecase.
func (f *fileUsecase) Upload(ctx context.Context, file io.Reader, metaData *domain.File) (*domain.File, error) {

	// filters and check for valid data type
	switch metaData.BucketName {
	case "profile-pictures":
		if metaData.ContentType != "image/jpeg" && metaData.ContentType != "image/png" {
			return &domain.File{}, domain.ErrInvalidFileFormat
		}
	case "documents":
		if metaData.ContentType != "application/pdf" {
			return &domain.File{}, domain.ErrInvalidFileFormat
		}
	default:
		return &domain.File{}, domain.ErrUnknownFileType
	}

	// generates a unique keyname for the object
	key := generateKeyName(metaData)
	metaData.CreatedAt = time.Now()
	err := f.s3s.Upload(ctx, metaData.BucketName, key, file, metaData.ContentType, metaData.Size)
	if err != nil {
		return nil, err
	}

	// store it to the database
	err = f.fileRepo.Store(ctx, metaData)
	if err != nil {
		return nil, err
	}
	return metaData, nil
}

func (f *fileUsecase) GetProfilePictureByUserID(ctx context.Context, userID string) (string, error) {
	// Fetch the profile picture file from the repository
	file, err := f.fileRepo.FindByUserID(ctx, userID)
	if err != nil || file == nil {
		return "", domain.ErrFileNotFound
	}

	// Only allow bucket "profile-pictures"
	if file.BucketName != "profile-pictures" {
		return "", domain.ErrFileNotFound
	}

	// Generate key for S3
	key := generateKeyName(file)

	// Generate presigned URL
	url, err := f.s3s.PresignedURL(ctx, file.BucketName, key)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", domain.ErrInternal)
	}

	return url, nil
}
