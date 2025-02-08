package s3_client

import (
	"context"
	"fitness-trainer/internal/logger"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Client struct {
	client *s3.Client
	bucket string
}

func New(client *s3.Client, bucket string) *Client {
	return &Client{
		client: client,
		bucket: bucket,
	}
}

func (c *Client) GeneratePutPresignedURL(ctx context.Context, key string) (string, error) {
	presignClient := s3.NewPresignClient(c.client)

	req, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	},
		s3.WithPresignExpires(time.Minute*15),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	logger.Debugf("generated presigned URL for key %s: %s", key, req.URL)

	return req.URL, nil
}
