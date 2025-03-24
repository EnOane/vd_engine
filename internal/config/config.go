package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type grpcConfig struct {
	ApiHost string
	ApiPort int
}

type dbConfig struct {
	DbHost string
	DbPort int
	DbMain string
	DbPass string
	DbUser string
}

type s3Config struct {
	S3Host      string
	S3Port      int
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
	S3Region    string
}

type Config struct {
	GrpcConfig grpcConfig
	DbConfig   dbConfig
	S3Config   s3Config
}

var AppConfig Config

func MustLoad() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("No .env file found")
	}

	AppConfig = Config{
		GrpcConfig: grpcConfig{
			ApiHost: getEnv("API_HOST", ""),
			ApiPort: getEnvAsInt("API_PORT", 50051),
		},
		DbConfig: dbConfig{
			DbHost: getEnv("DB_HOST", ""),
			DbPort: getEnvAsInt("DB_PORT", 5432),
			DbMain: getEnv("DB_MAIN", ""),
			DbUser: getEnv("DB_USERNAME", ""),
			DbPass: getEnv("DB_PASSWORD", ""),
		},
		S3Config: s3Config{
			S3Host:      getEnv("S3_HOST", ""),
			S3Port:      getEnvAsInt("S3_PORT", 9000),
			S3Bucket:    getEnv("S3_BUCKET", ""),
			S3AccessKey: getEnv("S3_ACCESS_KEY", ""),
			S3SecretKey: getEnv("S3_SECRET_KEY", ""),
			S3Region:    getEnv("S3_REGION", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
