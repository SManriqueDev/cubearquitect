package app

import (
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// App holds all dependencies and cleanup functions
type App struct {
	*fiber.App
	nodeTypeStore *orchestrator.NodeTypeStore
}

// New constructs a Fiber app wired with dependencies.
func New(cfg *config.Config) *App {
	client := cubepath.NewClient(cfg.BaseURL, cfg.Token)
	fiberApp := fiber.New(fiber.Config{
		AppName: "CubeArchitect API v1",
	})

	fiberApp.Use(logger.New())
	fiberApp.Use(cors.New())

	// Initialize node type store for persisting node types
	nodeTypeStore, err := orchestrator.NewNodeTypeStore(cfg.DataDir)
	if err != nil {
		log.Printf("Warning: Failed to initialize node type store: %v", err)
	}

	projectsService := service.NewProjectsService(client)
	vpsService := service.NewVPSService(client, cfg.ProjectID)
	pricingService := service.NewPricingService(client)

	// Initialize orchestrator
	orchestratorService := service.NewOrchestratorService(client, cfg.ProjectID, cfg)
	eventHub := orchestrator.NewEventHub()

	// Set event hub on the engine for event publishing
	orchestratorService.SetEventHub(eventHub)

	// Set node type store for tracking node types
	if nodeTypeStore != nil {
		orchestratorService.SetNodeTypeStore(nodeTypeStore)
		projectsService.SetNodeTypeStore(nodeTypeStore)
	}

	handlerSet := HandlerSet{
		Health:   handler.NewHealthHandler(),
		Projects: handler.NewProjectsHandler(projectsService),
		VPS:      handler.NewVPSHandler(vpsService),
		Pricing:  handler.NewPricingHandler(pricingService),
		Deploy:   handler.NewDeployHandler(orchestratorService, eventHub),
	}

	RegisterRoutes(fiberApp, handlerSet)

	return &App{
		App:           fiberApp,
		nodeTypeStore: nodeTypeStore,
	}
}

// Close cleans up resources
func (a *App) Close() error {
	if a.nodeTypeStore != nil {
		return a.nodeTypeStore.Close()
	}
	return nil
}
