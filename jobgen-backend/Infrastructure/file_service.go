package infrastructure

import (
	"context"
	"fmt"
	"io"
	domain "jobgen-backend/Domain"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorageService interface {
	GetFile(fileID string) (io.ReadCloser, error)
}

type localFileStorageService struct {
	basePath string
}

func NewFileStorageService(basePath string) FileStorageService {
	return &localFileStorageService{
		basePath: basePath,
	}
}

func (fs *localFileStorageService) GetFile(fileID string) (io.ReadCloser, error) {
	filePath := fs.basePath + "/" + fileID
	return os.Open(filePath)
}

type s3Service struct {
	minIO   *minio.Client
	maxSize int64 // maximum allowed upload size in bytes
	maxLife int64 // maximum url life in seconds
}

func NewFileService(url, accessKey, secretKey string, maxSize, maxLife int64) (domain.IFileService, error) {
	minioClient, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true, // set true if you are using https
		Region: "us-east-1",
	})
	return &s3Service{
		minIO:   minioClient,
		maxSize: maxSize,
		maxLife: maxLife,
	}, err
}

// Delete implements domain.IFileService.
func (s *s3Service) Delete(ctx context.Context, bucket string, key string) error {
	return s.minIO.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

// PresignedURL implements domain.IFileService.
func (s *s3Service) PresignedURL(ctx context.Context, bucket string, key string) (string, error) {
	url, err := s.minIO.PresignedGetObject(ctx, bucket, key, time.Duration(s.maxLife)*time.Second, nil)
	if err != nil {
		return "", fmt.Errorf("failed to check generate url: %w", domain.ErrInternal)
	}
	return url.String(), nil
}

// Upload implements domain.IFileService.
func (s *s3Service) Upload(ctx context.Context, bucket string, key string, file io.Reader, contentType string, size int64) error {

	if size > s.maxSize {
		return domain.ErrFileTooBig
	}

	exists, err := s.minIO.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("failed to check if bucket exist: %w", domain.ErrInternal)
	}

	if !exists {
		if err := s.minIO.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}

	_, err = s.minIO.PutObject(ctx, bucket, key, file, size, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"uploaded_at": time.Now().UTC().Format(time.RFC3339),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", domain.ErrInternal)
	}
	return nil
}
