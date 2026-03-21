package main

import (
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/app"
	"github.com/SManriqueDev/cubearchitect/internal/config"
)

func main() {
	cfg := config.Load()
	app := app.New(cfg)

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
