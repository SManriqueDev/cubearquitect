package app

import (
	"time"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/handler"
	"github.com/SManriqueDev/cubearchitect/internal/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type HandlerSet struct {
	Health   *handler.HealthHandler
	Projects *handler.ProjectsHandler
	VPS      *handler.VPSHandler
	Pricing  *handler.PricingHandler
	Deploy   *handler.DeployHandler
	SSHKeys  *handler.SSHKeysHandler
}

func RegisterRoutes(app *fiber.App, handlers HandlerSet, cfg *config.Config) {
	app.Get("/health", handlers.Health.GetHealth)

	api := app.Group("/api")

	api.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Please try again later.",
			})
		},
	}))

	api.Use(middleware.CubeTokenMiddleware(cfg.BaseURL))

	api.Get("/projects", handlers.Projects.GetProjects)
	api.Get("/ssh-keys", handlers.SSHKeys.GetSSHKeys)
	api.Get("/pricing", handlers.Pricing.GetPricing)
	api.Post("/deploy", handlers.Deploy.PostDeploy)
	api.Get("/deployments", handlers.Deploy.ListDeployments)
	api.Get("/deployments/:deployment_id", handlers.Deploy.GetDeploymentStatus)
	api.Get("/deployments/:deployment_id/events", websocket.New(handlers.Deploy.WebSocketDeploymentEvents))
}
