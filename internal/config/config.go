package config

import (
	"fmt"
	"os"

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
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load("../my.env")
	if err != nil {
		return &Config{}, fmt.Errorf("Error loading .env file: %w", err)
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
	}, nil
}
