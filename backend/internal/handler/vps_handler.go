package handler

import (
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/middleware"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
)

type VPSCreateRequest struct {
	ProjectID int                       `json:"project_id"`
	VPS       cubepath.VPSCreateRequest `json:"vps"`
}

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
	client, err := middleware.MustCubeClient(c)
	if err != nil {
		return err
	}

	var req VPSCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if req.ProjectID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
	}

	vps, err := h.service.Create(client, req.ProjectID, req.VPS)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(vps)
}
