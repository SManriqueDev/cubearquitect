package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// HandlerSet groups available HTTP handlers.
type HandlerSet struct {
	Health   *handler.HealthHandler
	Projects *handler.ProjectsHandler
	VPS      *handler.VPSHandler
	Pricing  *handler.PricingHandler
	Deploy   *handler.DeployHandler
	
	// Dependencies for WS
	OrchestratorSvc *service.OrchestratorService
	EventHub        *orchestrator.EventHub
}

// RegisterRoutes wires handlers to routes.
func RegisterRoutes(app *fiber.App, handlers HandlerSet) {
	app.Get("/health", handlers.Health.GetHealth)
	app.Get("/api/projects", handlers.Projects.GetProjects)
	app.Post("/api/vps", handlers.VPS.CreateVPS)
	app.Get("/api/pricing", handlers.Pricing.GetPricing)
	
	// Deployment orchestration routes
	app.Post("/api/deploy", handlers.Deploy.PostDeploy)
	app.Get("/api/deployments", handlers.Deploy.ListDeployments)
	app.Get("/api/deployments/:deployment_id", handlers.Deploy.GetDeploymentStatus)
	app.Get("/api/deployments/:deployment_id/events", websocket.New(handlers.Deploy.WebSocketDeploymentEvents))
}
