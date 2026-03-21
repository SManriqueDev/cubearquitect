package app

import (
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/gofiber/fiber/v2"
)

// HandlerSet groups available HTTP handlers.
type HandlerSet struct {
	Health   *handler.HealthHandler
	Projects *handler.ProjectsHandler
	VPS      *handler.VPSHandler
	Pricing  *handler.PricingHandler
}

// RegisterRoutes wires handlers to routes.
func RegisterRoutes(app *fiber.App, handlers HandlerSet) {
	app.Get("/health", handlers.Health.GetHealth)
	app.Get("/api/projects", handlers.Projects.GetProjects)
	app.Post("/api/vps", handlers.VPS.CreateVPS)
	app.Get("/api/pricing", handlers.Pricing.GetPricing)
}
