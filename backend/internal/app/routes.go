package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/middleware"
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
	SSHKeys  *handler.SSHKeysHandler
}

// RegisterRoutes wires handlers to routes.
func RegisterRoutes(app *fiber.App, handlers HandlerSet, cfg *config.Config) {
	app.Get("/health", handlers.Health.GetHealth)

	// Protected routes - require X-Cube-Token header
	protected := app.Group("", middleware.CubeTokenMiddleware(cfg.BaseURL))

	protected.Get("/api/projects", handlers.Projects.GetProjects)
	protected.Get("/api/ssh-keys", handlers.SSHKeys.GetSSHKeys)
	protected.Get("/api/pricing", handlers.Pricing.GetPricing)
	protected.Post("/api/deploy", handlers.Deploy.PostDeploy)
	protected.Get("/api/deployments", handlers.Deploy.ListDeployments)
	protected.Get("/api/deployments/:deployment_id", handlers.Deploy.GetDeploymentStatus)
	protected.Get("/api/deployments/:deployment_id/events", websocket.New(handlers.Deploy.WebSocketDeploymentEvents))
}
