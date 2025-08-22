package config

import (
	"fmt"
	"log"
	"os"
	"runway/logger"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               int
	AppsApiUrl         string
	ReviewsBaseUrl     string
	AppsStorageFile    string
	ReviewsStorageFile string
	TimeoutSecs        int
	Logger             logger.Config
}

func LoadConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}
	timeoutSecs, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT"))
	appPort, _ := strconv.Atoi(os.Getenv("PORT"))

	required := map[string]string{
		"APPLE_API_URL": os.Getenv("APPLE_API_URL"),
		"PORT":          os.Getenv("PORT"),
	}

	for key, value := range required {
		if value == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", key)
		}
	}

	loggerConfig := logger.Config{
		Level:    os.Getenv("LOG_LEVEL"),
		FilePath: os.Getenv("LOG_FILE_PATH"), // Empty means stdout only
	}
	return &Config{
		Port:               appPort,
		AppsApiUrl:         os.Getenv("APPLE_API_URL"),
		ReviewsBaseUrl:     os.Getenv("APPLE_REVIEWS_BASE_URL"),
		AppsStorageFile:    os.Getenv("APPS_STORAGE_FILE"),
		ReviewsStorageFile: os.Getenv("REVIEWS_STORAGE_FILE"),
		TimeoutSecs:        timeoutSecs,
		Logger:             loggerConfig,
	}, nil
}
