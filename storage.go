package storage

import (
	"context"
	"io"
)

type IStorage interface {
	Download(ctx context.Context, bucket, key string) (*[]byte, error)
	Delete(ctx context.Context, bucket, key string) error
	Upload(ctx context.Context, bucket, key string, body io.Reader) error
}
