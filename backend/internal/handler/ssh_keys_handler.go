package handler

import (
	"github.com/SManriqueDev/cubearchitect/internal/middleware"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
)

// SSHKeysHandler handles SSH key endpoints.
type SSHKeysHandler struct {
	service *service.SSHKeysService
}

// NewSSHKeysHandler creates a new SSH keys handler.
func NewSSHKeysHandler(svc *service.SSHKeysService) *SSHKeysHandler {
	return &SSHKeysHandler{service: svc}
}

// GetSSHKeys returns the list of SSH keys for the authenticated user.
func (h *SSHKeysHandler) GetSSHKeys(c *fiber.Ctx) error {
	client, err := middleware.MustCubeClient(c)
	if err != nil {
		return err
	}

	keys, err := h.service.List(client)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(keys)
}
