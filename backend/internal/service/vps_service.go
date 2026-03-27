package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// VPSService encapsulates CubePath VPS creation logic.
type VPSService struct{}

// NewVPSService returns a VPSService.
func NewVPSService() *VPSService {
	return &VPSService{}
}

// Create creates a VPS for the specified project.
func (s *VPSService) Create(client *cubepath.Client, projectID int, req cubepath.VPSCreateRequest) (*cubepath.VPS, error) {
	res, err := client.Post(fmt.Sprintf("/vps/create/%d", projectID), req)
	if err != nil {
		log.Printf("Error creating VPS via client: %v", err)
		return nil, err
	}

	var vps cubepath.VPS
	if err := json.Unmarshal(res, &vps); err != nil {
		return nil, err
	}

	log.Printf("VPS created: %+v", vps)
	return &vps, nil
}
