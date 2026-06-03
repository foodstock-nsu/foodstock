package port

import (
	"context"
	"io"
)

type MediaStorageGateway interface {
	Upload(ctx context.Context, file io.Reader, filename, contentType string) (string, error)
}
