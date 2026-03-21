package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// VPSService encapsulates CubePath VPS creation logic.
type VPSService struct {
	client    *cubepath.Client
	projectID int
}

// NewVPSService returns a VPSService.
func NewVPSService(client *cubepath.Client, projectID int) *VPSService {
	return &VPSService{client: client, projectID: projectID}
}

// Create creates a VPS for the configured project.
func (s *VPSService) Create(req cubepath.VPSCreateRequest) (*cubepath.VPS, error) {
	res, err := s.client.Post(fmt.Sprintf("/vps/create/%d", s.projectID), req)
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
