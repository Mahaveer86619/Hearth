package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server configuration
	HTTPPort int
	TCPPort  int

	// MinIO configuration
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string

	// Redis configuration
	RedisURL string
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	AppConfig = Config{
		HTTPPort: getEnvInt("HTTP_PORT", 4050),
		TCPPort:  getEnvInt("TCP_PORT", 4040),

		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),

		RedisURL: getEnv("REDIS_URL", "redis://redis:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := lookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := lookupEnv(key); exists {
		var intValue int
		_, err := fmt.Sscanf(value, "%d", &intValue)
		if err != nil {
			log.Printf("Error parsing integer from environment variable %s: %v", key, err)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}

var lookupEnv = func(key string) (string, bool) {
	return os.LookupEnv(key)
}