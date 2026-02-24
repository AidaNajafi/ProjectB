package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBhost     string
	DBport     string
	DBuser     string
	DBpass     string
	DBname     string
	DBsslmode  string
	SecretKey  string
	ServerPort string
	HealthURL  string
	BaseURL    string
	ApiKey     string
	TimeOut    time.Duration
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		return &Config{}, fmt.Errorf("Error loading .env file: %w", err)
	}
	timeStr := os.Getenv("TIME_OUT")
	timeOut, err := time.ParseDuration(timeStr)
	if err != nil {
		return &Config{}, fmt.Errorf("Failed to parse time duration %w", err)
	}
	return &Config{
		DBhost:     os.Getenv("DATABASE_HOST"),
		DBport:     os.Getenv("DATABASE_PORT"),
		DBuser:     os.Getenv("DATABASE_USER"),
		DBpass:     os.Getenv("DATABASE_PASSWORD"),
		DBname:     os.Getenv("DATABASE_NAME"),
		DBsslmode:  os.Getenv("DATABASE_SSLMODE"),
		SecretKey:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("SERVER_PORT"),
		HealthURL:  os.Getenv("API_HEALTH_URL"),
		ApiKey:     os.Getenv("PROVIDER_KEY"),
		BaseURL:    os.Getenv("BASE_URL"),
		TimeOut:    timeOut,
	}, nil
}
