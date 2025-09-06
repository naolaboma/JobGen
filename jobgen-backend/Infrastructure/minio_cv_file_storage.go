package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	domain "jobgen-backend/Domain"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type minioCVFileStorage struct {
	client *minio.Client
	bucket string
}


func NewMinioCVFileStorageService(url, accessKey, secretKey, bucket string) (*minioCVFileStorage, error) {
	cli, err := minio.New(url, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
		Region: "us-east-1",
	})
	if err != nil {
		return nil, err
	}
	// Ensure bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	return &minioCVFileStorage{client: cli, bucket: bucket}, nil
}

func (s *minioCVFileStorage) UploadFile(userID, category, fileName string, file multipart.File) (string, error) {
	// Buffer file to get size (CVs are small; limit enforced at usecase level)
	buf, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	key := buildObjectKey(userID, category)
	contentType := "application/pdf"

	// Upload to MinIO
	_, err = s.client.PutObject(context.Background(), s.bucket, key, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"uploaded_at": time.Now().UTC().Format(time.RFC3339),
			"category":    category,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload CV to storage: %w", domain.ErrInternal)
	}
	return s.encodeID(s.bucket, key), nil
}

func (s *minioCVFileStorage) DeleteFile(fileID string) error {
	bucket, key, err := s.decodeID(fileID)
	if err != nil {
		return err
	}
	return s.client.RemoveObject(context.Background(), bucket, key, minio.RemoveObjectOptions{})
}

func (s *minioCVFileStorage) GetFile(fileID string) (io.ReadCloser, error) {
	bucket, key, err := s.decodeID(fileID)
	if err != nil {
		return nil, err
	}
	obj, err := s.client.GetObject(context.Background(), bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve object: %w", domain.ErrNotFound)
	}
	return obj, nil
}

func (s *minioCVFileStorage) encodeID(bucket, key string) string {
	return bucket + "|" + key
}

func (s *minioCVFileStorage) decodeID(id string) (string, string, error) {
	parts := strings.SplitN(id, "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid file id format")
	}
	return parts[0], parts[1], nil
}

func buildObjectKey(userID, category string) string {
	id := uuid.NewString()
	prefix := "cv"
	if category != "" {
		prefix = strings.ToLower(category)
	}
	if userID != "" {
		return fmt.Sprintf("%s/%s/%s.pdf", prefix, userID, id)
	}
	return fmt.Sprintf("%s/%s.pdf", prefix, id)
}
