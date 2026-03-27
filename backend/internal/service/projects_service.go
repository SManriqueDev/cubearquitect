package service

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

type ProjectsService struct{}

func NewProjectsService() *ProjectsService {
	return &ProjectsService{}
}

type VPSItemWithType map[string]interface{}

type ProjectItemWithTypes struct {
	Project    cubepath.ProjectInfo `json:"project"`
	Networks   []interface{}        `json:"networks"`
	Baremetals []interface{}        `json:"baremetals"`
	VPS        []VPSItemWithType    `json:"vps"`
}

type ProjectsResponseWithTypes []ProjectItemWithTypes

func (s *ProjectsService) List(client *cubepath.Client) (cubepath.ProjectResponse, error) {
	res, err := client.Get("/projects/")
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

func (s *ProjectsService) ListWithNodeTypes(client *cubepath.Client) (ProjectsResponseWithTypes, error) {
	projects, err := s.List(client)
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

			if label, ok := vps["label"].(string); ok {
				if strings.HasPrefix(label, "app ") {
					vpsTyped["node_type"] = "app"
				} else if strings.HasPrefix(label, "database ") {
					vpsTyped["node_type"] = "database"
				}
			}

			item.VPS = append(item.VPS, vpsTyped)
		}

		result = append(result, item)
	}

	return result, nil
}
