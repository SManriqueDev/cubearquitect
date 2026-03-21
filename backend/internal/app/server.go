package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// New constructs a Fiber app wired with dependencies.
func New(cfg *config.Config) *fiber.App {
	client := cubepath.NewClient(cfg.BaseURL, cfg.Token)
	app := fiber.New(fiber.Config{
		AppName: "CubeArchitect API v1",
	})

	app.Use(logger.New())
	app.Use(cors.New())

	projectsService := service.NewProjectsService(client)
	vpsService := service.NewVPSService(client, cfg.ProjectID)
	pricingService := service.NewPricingService(client)

	handlerSet := HandlerSet{
		Health:   handler.NewHealthHandler(),
		Projects: handler.NewProjectsHandler(projectsService),
		VPS:      handler.NewVPSHandler(vpsService),
		Pricing:  handler.NewPricingHandler(pricingService),
	}

	RegisterRoutes(app, handlerSet)

	return app
}
