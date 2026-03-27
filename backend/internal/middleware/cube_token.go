package middleware

import (
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/gofiber/fiber/v2"
)

const (
	CubeClientKey = "cubeClient"
	CubeTokenKey  = "X-Cube-Token"
)

// CubeTokenMiddleware extracts the API token from X-Cube-Token header
// or query parameter (for WebSocket connections) and creates a CubePath client.
func CubeTokenMiddleware(defaultBaseURL string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get(CubeTokenKey)

		// Fallback to query parameter for WebSocket connections
		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "X-Cube-Token header or token query parameter is required",
			})
		}

		baseURL := defaultBaseURL
		if customURL := c.Get("X-Cube-API-URL"); customURL != "" {
			baseURL = customURL
		}

		client := cubepath.NewClient(baseURL, token)
		c.Locals(CubeClientKey, client)

		return c.Next()
	}
}

// GetCubeClient retrieves the CubePath client from context.
// Panics if middleware was not applied.
func GetCubeClient(c *fiber.Ctx) *cubepath.Client {
	client, ok := c.Locals(CubeClientKey).(*cubepath.Client)
	if !ok || client == nil {
		panic("CubeTokenMiddleware was not applied")
	}
	return client
}

// MustCubeClient retrieves the CubePath client or returns an error.
func MustCubeClient(c *fiber.Ctx) (*cubepath.Client, error) {
	client, ok := c.Locals(CubeClientKey).(*cubepath.Client)
	if !ok || client == nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "CubePath client not initialized")
	}
	return client, nil
}
