//go:build integration

package yacloud_test

import (
	"backend/internal/adapter/out/yacloud"
	"bytes"
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYandexS3Storage_Integration(t *testing.T) {
	var (
		accessKey  = "YCAJEArxLfMxD5DYKG-b-6lSs"
		secretKey  = "YCOWxXBb6v21R6SDYw_wyn20iKLDg632t8jBxa6s"
		bucketName = "foodstock-test"
		region     = "ru-central1"
		ctx        = context.Background()
	)

	cred := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(cred),
	)
	require.NoError(t, err)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
	})

	storage := yacloud.NewYandexS3Storage(s3Client, bucketName)

	fakeImageBytes := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15c4\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")

	t.Run("Upload File Success", func(t *testing.T) {
		fileReader := bytes.NewReader(fakeImageBytes)
		filename := "test_product.png"
		contentType := "image/png"

		s3Key, err := storage.Upload(ctx, fileReader, filename, contentType)

		require.NoError(t, err)
		assert.NotEmpty(t, s3Key)

		assert.Contains(t, s3Key, "items/")
		t.Logf("Uploaded successfully. S3 Key: %s", s3Key)

		t.Run("Verify Object Exists in S3", func(t *testing.T) {
			_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(s3Key),
			})
			require.NoError(t, err, "File was uploaded, but cannot be found in bucket")

			_, _ = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(s3Key),
			})
		})
	})

	t.Run("Auth or Bucket Error", func(t *testing.T) {
		badCred := credentials.NewStaticCredentialsProvider("wrong_key", "wrong_secret", "")
		badCfg, _ := config.LoadDefaultConfig(ctx,
			config.WithRegion(region),
			config.WithCredentialsProvider(badCred),
		)
		badS3Client := s3.NewFromConfig(badCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
		})

		badStorage := yacloud.NewYandexS3Storage(badS3Client, bucketName)
		fileReader := bytes.NewReader(fakeImageBytes)

		_, err = badStorage.Upload(ctx, fileReader, "error.png", "image/png")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to upload file to yandex s3 storage")
	})
}
