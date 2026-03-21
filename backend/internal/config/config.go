package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Token     string
	BaseURL   string
	Port      string
	ProjectID int
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Token:   os.Getenv("CUBE_API_TOKEN"),
		BaseURL: os.Getenv("CUBE_API_URL"),
		Port:    os.Getenv("PORT"),
	}

	if cfg.Token == "" {
		log.Fatal("❌ CUBE_API_TOKEN is required")
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.cubepath.com"
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	projectIDStr := os.Getenv("CUBE_PROJECT_ID")
	if projectIDStr == "" {
		log.Fatal("❌ CUBE_PROJECT_ID is required")
	}

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		log.Fatalf("❌ CUBE_PROJECT_ID must be an integer: %v", err)
	}

	cfg.ProjectID = projectID

	return cfg
}
