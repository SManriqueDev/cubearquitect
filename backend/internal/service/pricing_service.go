package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// PricingService wraps CubePath pricing retrieval.
type PricingService struct{}

// NewPricingService builds a PricingService.
func NewPricingService() *PricingService {
	return &PricingService{}
}

// GetPricing proxies pricing data from CubePath using the provided client.
func (s *PricingService) GetPricing(client *cubepath.Client) (json.RawMessage, error) {
	res, err := client.Get("/pricing/")
	if err != nil {
		log.Printf("Error fetching pricing: %v", err)
		return nil, err
	}
	return res, nil
}
