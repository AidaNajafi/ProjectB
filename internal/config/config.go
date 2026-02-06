package config

import (
	"flag"
	"fmt"
)

type Config struct {
	FilePath string
	Port     int
}

func LoadConfig(args []string) (Config, error) {
	fs := flag.NewFlagSet("authentication", flag.ContinueOnError)
	var cfg Config
	fs.StringVar(&cfg.FilePath, "file", "/home/aida-najafi/Go/project2.0", "our mock database")
	fs.IntVar(&cfg.Port, "port", 8080, "optional port")
	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}	

	if cfg.Port <= 0 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("port cannot be less than zero: %d", cfg.Port)
	}
	return cfg, nil
}
