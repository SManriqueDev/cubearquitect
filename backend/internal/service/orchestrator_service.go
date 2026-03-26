package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
)

type OrchestratorService struct {
	engine            *orchestrator.DeploymentEngine
	blueprintRegistry *orchestrator.BlueprintRegistry
	deployments       map[string]*deploymentState
	config            *config.Config
	mu                sync.RWMutex
}

type deploymentState struct {
	Context   *orchestrator.DeploymentContext
	StartedAt time.Time
	Error     error
}

func NewOrchestratorService(client *cubepath.Client, projectID int, cfg *config.Config) *OrchestratorService {
	registry := orchestrator.NewBlueprintRegistry()

	registry.Register(orchestrator.NewPostgresBasicBlueprint(cfg))
	registry.Register(orchestrator.NewNodeBasicBlueprint(cfg))

	engine := orchestrator.NewDeploymentEngine(client, projectID, registry)

	return &OrchestratorService{
		engine:            engine,
		blueprintRegistry: registry,
		config:            cfg,
		deployments:       make(map[string]*deploymentState),
	}
}

func (s *OrchestratorService) SetEventHub(hub *orchestrator.EventHub) {
	s.engine.SetEventHub(hub)
}

func (s *OrchestratorService) SetNodeTypeStore(store *orchestrator.NodeTypeStore) {
	s.engine.SetNodeTypeStore(store)
}

func (s *OrchestratorService) StartDeployment(payload *orchestrator.DeployPayload) (string, error) {
	deploymentID := fmt.Sprintf("deploy-%d", time.Now().UnixNano())

	if err := s.validatePayload(payload); err != nil {
		return "", fmt.Errorf("invalid payload: %w", err)
	}

	ctx := orchestrator.NewDeploymentContext(deploymentID, payload)

	services := make([]orchestrator.Service, len(payload.Nodes))
	for i, node := range payload.Nodes {
		deps := []string{}
		for _, edge := range payload.Edges {
			if edge.Target == node.ID {
				deps = append(deps, edge.Source)
			}
		}
		services[i] = orchestrator.Service{
			ID:        node.ID,
			DependsOn: deps,
		}
	}

	plan, err := orchestrator.GeneratePlan(services)
	if err != nil {
		return "", fmt.Errorf("failed to generate execution plan: %w", err)
	}

	ctx.Plan = plan

	s.mu.Lock()
	s.deployments[deploymentID] = &deploymentState{
		Context:   ctx,
		StartedAt: time.Now(),
	}
	s.mu.Unlock()

	go s.executeDeployment(deploymentID, ctx)

	return deploymentID, nil
}

func (s *OrchestratorService) GetDeploymentStatus(deploymentID string) (*orchestrator.DeploymentContext, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, exists := s.deployments[deploymentID]
	if !exists {
		return nil, fmt.Errorf("deployment not found: %s", deploymentID)
	}

	return state.Context, nil
}

func (s *OrchestratorService) ListDeployments() []*orchestrator.DeploymentContext {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*orchestrator.DeploymentContext
	for _, state := range s.deployments {
		result = append(result, state.Context)
	}
	return result
}

func (s *OrchestratorService) executeDeployment(deploymentID string, ctx *orchestrator.DeploymentContext) {
	log.Printf("Starting deployment: %s", deploymentID)

	deployCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	err := s.engine.ExecuteDeployment(deployCtx, ctx)

	s.mu.Lock()
	if state, exists := s.deployments[deploymentID]; exists {
		state.Error = err
	}
	s.mu.Unlock()

	if err != nil {
		log.Printf("Deployment failed: %s - %v", deploymentID, err)
	} else {
		log.Printf("Deployment completed: %s", deploymentID)
	}
}

func (s *OrchestratorService) validatePayload(payload *orchestrator.DeployPayload) error {
	if len(payload.Nodes) == 0 {
		return fmt.Errorf("payload must contain at least one node")
	}

	nodeIDs := make(map[string]bool)
	for _, node := range payload.Nodes {
		if node.ID == "" {
			return fmt.Errorf("node ID cannot be empty")
		}
		if _, exists := nodeIDs[node.ID]; exists {
			return fmt.Errorf("duplicate node ID: %s", node.ID)
		}
		nodeIDs[node.ID] = true

		if node.Type != orchestrator.NodeTypeApp &&
			node.Type != orchestrator.NodeTypeDatabase {
			return fmt.Errorf("unknown node type: %s", node.Type)
		}
	}

	for _, edge := range payload.Edges {
		if !nodeIDs[edge.Source] {
			return fmt.Errorf("edge source references unknown node: %s", edge.Source)
		}
		if !nodeIDs[edge.Target] {
			return fmt.Errorf("edge target references unknown node: %s", edge.Target)
		}
	}

	return nil
}
