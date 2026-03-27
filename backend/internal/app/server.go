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

type App struct {
	*fiber.App
}

func New(cfg *config.Config) *App {
	client := cubepath.NewClient(cfg.BaseURL, cfg.Token)
	fiberApp := fiber.New(fiber.Config{
		AppName: "CubeArchitect API v1",
	})

	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New())

	projectsService := service.NewProjectsService(client)
	vpsService := service.NewVPSService(client, cfg.ProjectID)
	pricingService := service.NewPricingService(client)

	orchestratorService := service.NewOrchestratorService(client, cfg.ProjectID, cfg)
	eventHub := orchestrator.NewEventHub()

	orchestratorService.SetEventHub(eventHub)

	handlerSet := HandlerSet{
		Health:   handler.NewHealthHandler(),
		Projects: handler.NewProjectsHandler(projectsService),
		VPS:      handler.NewVPSHandler(vpsService),
		Pricing:  handler.NewPricingHandler(pricingService),
		Deploy:   handler.NewDeployHandler(orchestratorService, eventHub),
	}

	RegisterRoutes(fiberApp, handlerSet)

	return &App{
		App: fiberApp,
	}
}

func (a *App) Close() error {
	return nil
}
