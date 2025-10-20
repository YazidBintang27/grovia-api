package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	MLAPIURL string
	Aws      AwsConfig
}

type AwsConfig struct {
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

func LoadConfig() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("configs")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &AppConfig{
		MLAPIURL: viper.GetString("ml_api_url"),
		Aws: AwsConfig{
			Region:    os.Getenv("AWS_REGION"),
			Bucket:    os.Getenv("AWS_S3_BUCKET"),
			AccessKey: os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		},
	}
}
