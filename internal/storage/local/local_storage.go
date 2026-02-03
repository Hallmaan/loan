package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type LocalStorage struct {
	basePath string
	baseURL  string
}

func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

func (s *LocalStorage) Save(ctx context.Context, filename string, reader io.Reader) (string, error) {
	ext := filepath.Ext(filename)
	newFilename := uuid.New().String() + ext

	fullPath := filepath.Join(s.basePath, newFilename)

	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return newFilename, nil
}

func (s *LocalStorage) GetURL(path string) string {
	return fmt.Sprintf("%s/uploads/%s", s.baseURL, path)
}
