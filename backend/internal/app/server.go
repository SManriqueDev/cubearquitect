package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type App struct {
	*fiber.App
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	fiberApp := fiber.New(fiber.Config{
		AppName: "CubeArchitect API v1",
	})

	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New())

	// Services are now created per-request via handlers using c.Locals
	projectsService := service.NewProjectsService()
	vpsService := service.NewVPSService()
	pricingService := service.NewPricingService()
	sshKeysService := service.NewSSHKeysService()

	// Orchestrator service needs project ID from user, not from config
	// We'll pass it dynamically in the handler
	orchestratorService := service.NewOrchestratorService(nil, 0, cfg)
	eventHub := orchestrator.NewEventHub()

	orchestratorService.SetEventHub(eventHub)

	handlerSet := HandlerSet{
		Health:   handler.NewHealthHandler(),
		Projects: handler.NewProjectsHandler(projectsService),
		VPS:      handler.NewVPSHandler(vpsService),
		Pricing:  handler.NewPricingHandler(pricingService),
		Deploy:   handler.NewDeployHandler(orchestratorService, eventHub),
		SSHKeys:  handler.NewSSHKeysHandler(sshKeysService),
	}

	RegisterRoutes(fiberApp, handlerSet, cfg)

	return &App{
		App: fiberApp,
		cfg: cfg,
	}
}

func (a *App) Close() error {
	return nil
}
