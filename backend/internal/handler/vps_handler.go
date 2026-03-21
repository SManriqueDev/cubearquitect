package handler

import (
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
)

// VPSHandler handles VPS creation requests.
type VPSHandler struct {
	service *service.VPSService
}

// NewVPSHandler constructs a VPSHandler with the required service.
func NewVPSHandler(svc *service.VPSService) *VPSHandler {
	return &VPSHandler{service: svc}
}

// CreateVPS proxies the request body to CubePath and returns the result.
func (h *VPSHandler) CreateVPS(c *fiber.Ctx) error {
	var req cubepath.VPSCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	vps, err := h.service.Create(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(vps)
}
