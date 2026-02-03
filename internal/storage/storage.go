package storage

import (
	"context"
	"io"
)

type Storage interface {
	Save(ctx context.Context, path string, reader io.Reader) (string, error)
	GetURL(path string) string
}
