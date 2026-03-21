package handler

import "github.com/gofiber/fiber/v2"

// HealthHandler handles health check requests.
type HealthHandler struct{}

// NewHealthHandler returns a new HealthHandler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// GetHealth responds with a simple status payload.
func (h *HealthHandler) GetHealth(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "alive"})
}
