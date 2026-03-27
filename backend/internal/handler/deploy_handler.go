package handler

import (
	"log"
	"time"

	"github.com/SManriqueDev/cubearchitect/internal/middleware"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
	"github.com/SManriqueDev/cubearchitect/internal/service"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// DeployHandler handles deployment requests.
type DeployHandler struct {
	orchestratorSvc *service.OrchestratorService
	eventHub        *orchestrator.EventHub
}

// NewDeployHandler creates a new deploy handler.
func NewDeployHandler(orchestratorSvc *service.OrchestratorService, eventHub *orchestrator.EventHub) *DeployHandler {
	return &DeployHandler{
		orchestratorSvc: orchestratorSvc,
		eventHub:        eventHub,
	}
}

// PostDeploy initiates a deployment.
func (h *DeployHandler) PostDeploy(c *fiber.Ctx) error {
	client, err := middleware.MustCubeClient(c)
	if err != nil {
		return err
	}

	var payload orchestrator.DeployPayload
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if payload.ProjectID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "project_id is required")
	}

	deploymentID, err := h.orchestratorSvc.StartDeployment(client, &payload)
	if err != nil {
		log.Printf("Deployment failed: %v", err)
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	deployCtx, _ := h.orchestratorSvc.GetDeploymentStatus(deploymentID)

	response := orchestrator.DeployResponse{
		Success:       true,
		DeploymentID:  deploymentID,
		Message:       "Deployment initiated",
		NodesCount:    len(payload.Nodes),
		EdgesCount:    len(payload.Edges),
		ExecutionPlan: len(deployCtx.Plan),
	}

	return c.Status(fiber.StatusAccepted).JSON(response)
}

// GetDeploymentStatus returns the status of a deployment.
func (h *DeployHandler) GetDeploymentStatus(c *fiber.Ctx) error {
	deploymentID := c.Params("deployment_id")
	if deploymentID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "deployment_id is required")
	}

	deployCtx, err := h.orchestratorSvc.GetDeploymentStatus(deploymentID)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	// Build response with node statuses
	nodeStatuses := make(map[string]interface{})
	for nodeID, status := range deployCtx.NodeStatuses {
		nodeStatuses[nodeID] = map[string]interface{}{
			"status": status.Status,
			"error":  status.Error,
			"vps":    status.VPSInfo,
		}
	}

	return c.JSON(fiber.Map{
		"deployment_id": deploymentID,
		"node_statuses": nodeStatuses,
		"plan_levels":   len(deployCtx.Plan),
	})
}

// WebSocketDeploymentEvents establishes a WebSocket for deployment events.
func (h *DeployHandler) WebSocketDeploymentEvents(c *websocket.Conn) {
	deploymentID := c.Params("deployment_id")
	if deploymentID == "" {
		c.WriteMessage(websocket.TextMessage, []byte(`{"error":"deployment_id is required"}`))
		c.Close()
		return
	}

	log.Printf("WebSocket client connected for deployment: %s", deploymentID)

	// Send initial connection message
	c.WriteJSON(fiber.Map{
		"type":      "connected",
		"message":   "Connected to deployment stream",
		"timestamp": time.Now().UnixMilli(),
	})

	// Atomically subscribe and drain any buffered events so no events are missed
	// between clearing the buffer and registering the subscriber.
	eventCh, bufferedEvents := h.eventHub.SubscribeAndDrain(deploymentID)
	defer h.eventHub.Unsubscribe(deploymentID, eventCh)

	// Send buffered events first (for late subscribers)
	for _, event := range bufferedEvents {
		if err := c.WriteJSON(event); err != nil {
			log.Printf("WebSocket error sending buffered event: %v", err)
			return
		}
	}

	// Forward events to client
	for event := range eventCh {
		if err := c.WriteJSON(event); err != nil {
			log.Printf("WebSocket write error for deployment %s: %v", deploymentID, err)
			break
		}
	}
}

// ListDeployments returns all active deployments.
func (h *DeployHandler) ListDeployments(c *fiber.Ctx) error {
	deployments := h.orchestratorSvc.ListDeployments()

	result := make([]map[string]interface{}, len(deployments))
	for i, deploy := range deployments {
		result[i] = map[string]interface{}{
			"deployment_id": deploy.DeploymentID,
			"nodes_count":   len(deploy.Nodes),
			"plan_levels":   len(deploy.Plan),
		}
	}

	return c.JSON(result)
}
