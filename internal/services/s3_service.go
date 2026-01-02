package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"grovia/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
	region     string
}

type S3Uploader interface {
	UploadFile(file *multipart.FileHeader, folder string) (string, error)
}

func NewS3Service(cfg configs.AwsConfig) *S3Service {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if accessKey == "" || secretKey == "" {
		log.Fatal("AWS credentials are not set. Please set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY")
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Service{
		client:     client,
		bucketName: cfg.Bucket,
		region:     cfg.Region,
	}
}

func (s *S3Service) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (string, error) {
	const MaxFileSize = 10 << 20
	if file.Size > MaxFileSize {
		return "", fmt.Errorf("file too large")
	}

	allowedExts := map[string]bool{".jpg": true, ".png": true, ".jpeg": true}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("file type not allowed")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s/%s%s", folder, uuid.New().String(), ext)

	contentType := file.Header.Get("Content-Type")

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(fileName),
		Body:        src,
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPrivate,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.bucketName, s.region, fileName)

	return fileURL, nil
}
