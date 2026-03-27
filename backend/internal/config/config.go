package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL string
	Port    string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		BaseURL: os.Getenv("CUBE_API_URL"),
		Port:    os.Getenv("PORT"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.cubepath.com"
		log.Println("ℹ️  CUBE_API_URL not set, using default: https://api.cubepath.com")
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	log.Printf("ℹ️  Server configured for CubePath API: %s", cfg.BaseURL)

	return cfg
}
