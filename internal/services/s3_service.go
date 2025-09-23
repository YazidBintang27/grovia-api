package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"grovia/configs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	client     *s3.Client
	bucketName string
}

type S3Uploader interface {
	UploadFile(file *multipart.FileHeader, folder string) (string, error)
}

func NewS3Service(cfg configs.AwsConfig) *S3Service {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			getEnv("AWS_ACCESS_KEY_ID", ""),
			getEnv("AWS_SECRET_ACCESS_KEY", ""),
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
	}
}

func (s *S3Service) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s/%d_%s", folder, time.Now().Unix(), file.Filename)

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fileName),
		Body:   src,
		ACL:    "public-read", 
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, getEnv("AWS_REGION", ""), fileName)

	return url, nil
}

func getEnv(key, fallback string) string {
	if value, exists := lookupEnv(key); exists {
		return value
	}
	return fallback
}

func lookupEnv(key string) (string, bool) {
	value := ""
	if v := getenv(key); v != "" {
		value = v
		return value, true
	}
	return value, false
}

func getenv(key string) string {
	return fmt.Sprintf("%s", (func() string {
		if val, ok := lookup(key); ok {
			return val
		}
		return ""
	})())
}

func lookup(key string) (string, bool) {
	return "", false 
}