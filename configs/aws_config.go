package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AwsConfig struct {
	Region string
	Bucket string
}

func LoadAwsConfig() AwsConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return AwsConfig{
		Region: os.Getenv("AWS_REGION"),
		Bucket: os.Getenv("AWS_S3_BUCKET"),
	}
}