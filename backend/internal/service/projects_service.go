package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// ProjectsService wraps the Cubepath client for project operations.
type ProjectsService struct {
	client *cubepath.Client
}

// NewProjectsService builds a ProjectsService.
func NewProjectsService(client *cubepath.Client) *ProjectsService {
	return &ProjectsService{client: client}
}

// List returns the projects fetched from CubePath.
func (s *ProjectsService) List() (cubepath.ProjectResponse, error) {
	res, err := s.client.Get("/projects/")
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		return nil, err
	}

	var projects cubepath.ProjectResponse
	if err := json.Unmarshal(res, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}
