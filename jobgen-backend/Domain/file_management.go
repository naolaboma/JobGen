package domain

import (
	"context"
	"io"
	"time"
)

// File entity
type File struct {
	ID          string    `json:"_id" bson:"_id,omitempty"`
	UserID      string    `json:"user_id" bson:"user_id"`
	UniqueID    string    `json:"unique_id" bson:"unique_id"`
	FileName    string    `json:"file_name" bson:"file_name"`
	BucketName  string    `json:"bucket_name" bson:"bucket_name"`
	ContentType string    `json:"content_type" bson:"content_type"`
	Size        int64     `json:"size" bson:"size"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

type IFileService interface {
	Upload(ctx context.Context, bucket, key string, file io.Reader, contentType string, size int64) error
	Delete(ctx context.Context, bucket, key string) error
	PresignedURL(ctx context.Context, bucket, key string) (string, error)
}

// Usecase Interface
type IFileUsecase interface {
	// Upload uploads a file to storage and returns its metadata
	Upload(ctx context.Context, file io.Reader, metaData *File) (*File, error)

	// Download retrieves a file from storage as a stream
	Download(ctx context.Context, fileID string, userID string) (string, error)

	// Delete removes a file from storage
	Delete(ctx context.Context, ID string, userID string) error

	// Exists checks whether a file exists in storage
	Exists(ctx context.Context, ID string) (bool, error)

	// Gets user profile url
	GetProfilePictureByUserID(ctx context.Context, ID string) (string, error)
}

// File Repository Interface
type IFileRepository interface {
	Store(ctx context.Context, file *File) error
	FindByID(ctx context.Context, ID string) (*File, error)
	FindByUserID(ctx context.Context, ID string) (*File, error)
	Delete(ctx context.Context, ID string) error
}
