package main

import (
	"log"
	"os"
	"strconv"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.Load()
	client := cubepath.NewClient(cfg.BaseURL, cfg.Token)

	app := fiber.New(fiber.Config{
		AppName: "CubeArchitect API v1",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "alive"})
	})

	app.Get("/api/projects", func(c *fiber.Ctx) error {
		projects, err := client.GetProjects()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(projects)
	})

	app.Post("/api/vps", func(c *fiber.Ctx) error {
		projectID := os.Getenv("CUBE_PROJECT_ID")
		projectIDInt, err := strconv.Atoi(projectID)

		var req cubepath.VPSCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}
		vps, err := client.CreateVPS(projectIDInt, req)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(201).JSON(vps)
	})

	app.Get("/api/pricing", func(c *fiber.Ctx) error {
		pricing, err := client.GetPricing()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(pricing)
	})

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
