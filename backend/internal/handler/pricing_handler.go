package handler

import (
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
)

// PricingHandler handles pricing-related endpoints.
type PricingHandler struct {
	service *service.PricingService
}

// NewPricingHandler builds a PricingHandler with its dependencies.
func NewPricingHandler(svc *service.PricingService) *PricingHandler {
	return &PricingHandler{service: svc}
}

// GetPricing returns the pricing payload from CubePath.
func (h *PricingHandler) GetPricing(c *fiber.Ctx) error {
	pricing, err := h.service.GetPricing()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(pricing)
}
