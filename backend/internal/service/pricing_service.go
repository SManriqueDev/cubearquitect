package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// PricingService wraps CubePath pricing retrieval.
type PricingService struct {
	client *cubepath.Client
}

// NewPricingService builds a PricingService.
func NewPricingService(client *cubepath.Client) *PricingService {
	return &PricingService{client: client}
}

// GetPricing proxies pricing data from CubePath.
func (s *PricingService) GetPricing() (json.RawMessage, error) {
	res, err := s.client.Get("/pricing/")
	if err != nil {
		log.Printf("Error fetching pricing: %v", err)
		return nil, err
	}
	return res, nil
}
