package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

type Config struct {
	FilePath  string `yaml:"path"`
	Port      int    `yaml:"port"`
	SecretKey string
}

func LoadConfig(filepath string) (Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("Failed to decode file :%w", err)
	}
	err = godotenv.Load("./my.env")
	if err != nil {
		return Config{}, fmt.Errorf("Error loading .env file: %w", err)
	}
	cfg.SecretKey = os.Getenv("JWT_SECRET")
	if cfg.SecretKey == "" {
		return Config{}, fmt.Errorf("Failed to get secret key: %w", err)
	}

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("port cannot be less than zero: %d", cfg.Port)
	}
	return cfg, nil
}
