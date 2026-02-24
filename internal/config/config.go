package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBhost       string
	DBport       string
	DBuser       string
	DBpass       string
	DBname       string
	DBsslmode    string
	SecretKey    string
	ServerPort   string
	HealthURL    string
	BaseURL      string
	ApiKey       string
	TimeOut      time.Duration
	CBname       string
	CBMaxRequest uint32
	CBInterval   time.Duration
	CBTimeOut    time.Duration
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file: %w", err)
	}
	timeStr := os.Getenv("TIME_OUT")
	timeOut, err := time.ParseDuration(timeStr)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse time duration %w", err)
	}
	intervalStr := os.Getenv("CB_INTERVAL")
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return nil, err
	}

	timeOutCBStr := os.Getenv("CB_TIMEOUT")
	timeOutCB, err := time.ParseDuration(timeOutCBStr)
	if err != nil {
		return nil, err
	}

	maxReqStr := os.Getenv("CB_MAX_REQUEST")
	maxReqInt, err := strconv.ParseFloat(maxReqStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid CB_MAX_REQUEST: %w", err)
	}
	if maxReqInt <= 0 {
		return nil, fmt.Errorf("CB_MAX_REQUEST must be > 0")
	}
	maxReq := uint32(maxReqInt)
	return &Config{
		DBhost:       os.Getenv("DATABASE_HOST"),
		DBport:       os.Getenv("DATABASE_PORT"),
		DBuser:       os.Getenv("DATABASE_USER"),
		DBpass:       os.Getenv("DATABASE_PASSWORD"),
		DBname:       os.Getenv("DATABASE_NAME"),
		DBsslmode:    os.Getenv("DATABASE_SSLMODE"),
		SecretKey:    os.Getenv("JWT_SECRET"),
		ServerPort:   os.Getenv("SERVER_PORT"),
		HealthURL:    os.Getenv("API_HEALTH_URL"),
		ApiKey:       os.Getenv("PROVIDER_KEY"),
		BaseURL:      os.Getenv("BASE_URL"),
		TimeOut:      timeOut,
		CBname:       os.Getenv("CB_NAME"),
		CBMaxRequest: maxReq,
		CBInterval:   interval,
		CBTimeOut:    timeOutCB,
	}, nil
}
