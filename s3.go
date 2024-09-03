package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	config   *AwsConfig
	s3Client *s3.Client
}

type AwsConfig struct {
	Region    string
	AccessKey string
	SecretKey string
}

func NewS3(cfg *AwsConfig) IStorage {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, "")),
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		panic(err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	return &S3Storage{
		config:   cfg,
		s3Client: s3Client,
	}
}

func (r *S3Storage) Download(ctx context.Context, bucket, key string) (*[]byte, error) {
	out, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return nil, fmt.Errorf("couldn't get object from S3 %v:%v: %v", bucket, key, err)
	}
	defer out.Body.Close()

	body, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read object body from %v: %v", key, err)
	}
	return &body, nil
}

func (r *S3Storage) Delete(ctx context.Context, bucket, key string) error {
	_, err := r.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return fmt.Errorf("couldn't delete object from S3 %v:%v: %v", bucket, key, err)
	}
	return nil
}

func (r *S3Storage) Upload(ctx context.Context, bucket, key string, body io.Reader) error {
	_, err := r.s3Client.PutObject(ctx, &s3.PutObjectInput{Bucket: &bucket, Key: &key, Body: body})
	if err != nil {
		return fmt.Errorf("couldn't get upload object to S3 %v:%v: %v", bucket, key, err)
	}
	return nil
}
