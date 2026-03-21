package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
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
	
	// Initialize orchestrator
	orchestratorService := service.NewOrchestratorService(client, cfg.ProjectID, cfg)
	eventHub := orchestrator.NewEventHub()
	
	// Set event hub on the engine for event publishing
	orchestratorService.SetEventHub(eventHub)

	handlerSet := HandlerSet{
		Health:        handler.NewHealthHandler(),
		Projects:      handler.NewProjectsHandler(projectsService),
		VPS:           handler.NewVPSHandler(vpsService),
		Pricing:       handler.NewPricingHandler(pricingService),
		Deploy:        handler.NewDeployHandler(orchestratorService, eventHub),
		OrchestratorSvc: orchestratorService,
		EventHub:       eventHub,
	}

	RegisterRoutes(app, handlerSet)

	return app
}
