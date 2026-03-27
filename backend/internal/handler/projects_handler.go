package handler

import (
	"github.com/SManriqueDev/cubearchitect/internal/middleware"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/fiber/v2"
)

// ProjectsHandler handles project-related endpoints.
type ProjectsHandler struct {
	service *service.ProjectsService
}

// NewProjectsHandler constructs a ProjectsHandler with the required service.
func NewProjectsHandler(svc *service.ProjectsService) *ProjectsHandler {
	return &ProjectsHandler{service: svc}
}

// GetProjects returns the project list forwarded from CubePath.
func (h *ProjectsHandler) GetProjects(c *fiber.Ctx) error {
	client, err := middleware.MustCubeClient(c)
	if err != nil {
		return err
	}

	projects, err := h.service.ListWithNodeTypes(client)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(projects)
}
