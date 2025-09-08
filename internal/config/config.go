package config

import (
	"os"
)

type Config struct {
	Port           string
	DatabaseURL    string
	RedisAddr      string
	JWTSecret      string
	S3Endpoint     string
	S3Region       string
	S3AccessKey    string
	S3SecretKey    string
	S3Bucket       string
	UploadMaxBytes int64
}

func FromEnv() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/store?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "devsecret"),
		S3Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
		S3Region:    getEnv("S3_REGION", "us-east-1"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		S3Bucket:    getEnv("S3_BUCKET", "stories"),
	}
}

func getEnv(key, def string) string { if v := os.Getenv(key); v != "" { return v }; return def }
