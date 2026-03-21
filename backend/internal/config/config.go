package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Token   string
	BaseURL string
	Port    string
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

	return cfg
}