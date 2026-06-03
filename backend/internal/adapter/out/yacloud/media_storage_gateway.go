package yacloud

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type YandexS3Storage struct {
	s3Client   *s3.Client
	bucketName string
}

func NewYandexS3Storage(s3Client *s3.Client, bucketName string) *YandexS3Storage {
	return &YandexS3Storage{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (s *YandexS3Storage) Upload(ctx context.Context, file io.Reader, filename, contentType string) (string, error) {
	var ext string

	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		ext = filepath.Ext(filename)
	} else {
		ext = extensions[0]
	}

	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	s3Key := fmt.Sprintf("items/%s", newFilename)

	_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to yandex s3 storage: %w", err)
	}

	return s3Key, nil
}
