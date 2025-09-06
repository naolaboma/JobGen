package infrastructure

import (
	"fmt"
	"io"
	domain "jobgen-backend/Domain"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)


type localCVFileStorage struct {
	basePath string
}

func NewLocalCVFileStorageService(basePath string) *localCVFileStorage {
	// Ensure base path exists
	_ = os.MkdirAll(basePath, 0o755)
	return &localCVFileStorage{basePath: basePath}
}

func (s *localCVFileStorage) UploadFile(userID, category, fileName string, file multipart.File) (string, error) {
	id := uuid.NewString()
	// Store flat by id to keep lookups simple
	path := filepath.Join(s.basePath, id)

	// Reset file pointer if possible (best-effort)
	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Touch mtime for traceability (optional)
	_ = os.Chtimes(path, time.Now(), time.Now())
	return id, nil
}

// DeleteFile removes the stored file by its ID.
func (s *localCVFileStorage) DeleteFile(fileID string) error {
	path := filepath.Join(s.basePath, fileID)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %w", domain.ErrFileNotFound)
		}
		return err
	}
	return nil
}

// GetFile returns a read-only stream for a previously stored file.
func (s *localCVFileStorage) GetFile(fileID string) (io.ReadCloser, error) {
	path := filepath.Join(s.basePath, fileID)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", domain.ErrFileNotFound)
		}
		return nil, err
	}
	return f, nil
}
