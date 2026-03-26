package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// NodeTypeStoreInterface defines the interface for node type storage
type NodeTypeStoreInterface interface {
	Get(vpsID int) (string, bool)
	GetAll() map[int]string
}

// ProjectsService wraps the Cubepath client for project operations.
type ProjectsService struct {
	client        *cubepath.Client
	nodeTypeStore NodeTypeStoreInterface
}

// NewProjectsService builds a ProjectsService.
func NewProjectsService(client *cubepath.Client) *ProjectsService {
	return &ProjectsService{client: client}
}

// SetNodeTypeStore sets the node type store for enriching VPS data.
func (s *ProjectsService) SetNodeTypeStore(store NodeTypeStoreInterface) {
	s.nodeTypeStore = store
}

// VPSItemWithType represents a VPS with node_type field
type VPSItemWithType map[string]interface{}

// ProjectItemWithTypes extends ProjectItem with typed VPS items
type ProjectItemWithTypes struct {
	Project    cubepath.ProjectInfo `json:"project"`
	Networks   []interface{}        `json:"networks"`
	Baremetals []interface{}        `json:"baremetals"`
	VPS        []VPSItemWithType    `json:"vps"`
}

// ProjectsResponseWithTypes is the full response with node types
type ProjectsResponseWithTypes []ProjectItemWithTypes

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

// ListWithNodeTypes returns projects with node_type information for each VPS.
func (s *ProjectsService) ListWithNodeTypes() (ProjectsResponseWithTypes, error) {
	projects, err := s.List()
	if err != nil {
		return nil, err
	}

	var result ProjectsResponseWithTypes
	for _, proj := range projects {
		item := ProjectItemWithTypes{
			Project:    proj.Project,
			Networks:   proj.Networks,
			Baremetals: proj.Baremetals,
			VPS:        make([]VPSItemWithType, 0),
		}

		for _, vpsRaw := range proj.VPS {
			vps, ok := vpsRaw.(map[string]interface{})
			if !ok {
				continue
			}

			vpsTyped := VPSItemWithType(vps)

			if s.nodeTypeStore != nil {
				if id, ok := vps["id"].(float64); ok {
					vpsID := int(id)
					if nodeType, exists := s.nodeTypeStore.Get(vpsID); exists {
						vpsTyped["node_type"] = nodeType
					}
				}
			}

			item.VPS = append(item.VPS, vpsTyped)
		}

		result = append(result, item)
	}

	return result, nil
}
