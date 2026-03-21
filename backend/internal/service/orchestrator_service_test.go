package service

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/orchestrator"
)

// stubCubePathClient is a no-op ClientInterface implementation for unit tests.
// It returns an immediate error for all calls so that the background execution
// goroutine spawned by StartDeployment fails fast without making real HTTP calls.
type stubCubePathClient struct{}

func (s *stubCubePathClient) Get(path string) (json.RawMessage, error) {
	return nil, errors.New("stub: no HTTP calls in unit tests")
}
func (s *stubCubePathClient) Post(path string, body interface{}) (json.RawMessage, error) {
	return nil, errors.New("stub: no HTTP calls in unit tests")
}
func (s *stubCubePathClient) Put(path string, body interface{}) (json.RawMessage, error) {
	return nil, errors.New("stub: no HTTP calls in unit tests")
}
func (s *stubCubePathClient) Patch(path string, body interface{}) (json.RawMessage, error) {
	return nil, errors.New("stub: no HTTP calls in unit tests")
}
func (s *stubCubePathClient) Delete(path string) (json.RawMessage, error) {
	return nil, errors.New("stub: no HTTP calls in unit tests")
}

func newTestOrchestratorService() *OrchestratorService {
	cfg := &config.Config{
		Token:       "fake-token",
		BaseURL:     "https://api.example.com",
		Port:        "8080",
		ProjectID:   123,
		SSHKeyNames: "",
	}
	return NewOrchestratorService(&stubCubePathClient{}, 123, cfg)
}

func TestOrchestratorService_StartDeployment(t *testing.T) {
	svc := newTestOrchestratorService()

	payload := &orchestrator.DeployPayload{
		Nodes: []orchestrator.DeployNode{
			{
				ID:        "db-1",
				Kind:      orchestrator.NodeKindDatabase,
				Name:      "postgres-db",
				Label:     "PostgreSQL",
				Blueprint: "postgres-basic",
			},
			{
				ID:    "app-1",
				Kind:  orchestrator.NodeKindApp,
				Name:  "api-server",
				Label: "API Server",
			},
		},
		Edges: []orchestrator.DeployEdge{
			{
				Source: "db-1",
				Target: "app-1",
			},
		},
	}

	deploymentID, err := svc.StartDeployment(payload)
	if err != nil {
		t.Fatalf("failed to start deployment: %v", err)
	}

	if deploymentID == "" {
		t.Fatal("deployment ID should not be empty")
	}

	status, err := svc.GetDeploymentStatus(deploymentID)
	if err != nil {
		t.Fatalf("failed to get status: %v", err)
	}

	if status.DeploymentID != deploymentID {
		t.Errorf("expected deployment ID %s, got %s", deploymentID, status.DeploymentID)
	}

	// Verify execution plan
	if len(status.Plan) != 2 {
		t.Errorf("expected 2 execution levels, got %d", len(status.Plan))
	}

	// Verify level 0 has db-1
	if len(status.Plan[0]) != 1 || status.Plan[0][0] != "db-1" {
		t.Errorf("expected level 0 to be [db-1], got %v", status.Plan[0])
	}

	// Verify level 1 has app-1
	if len(status.Plan[1]) != 1 || status.Plan[1][0] != "app-1" {
		t.Errorf("expected level 1 to be [app-1], got %v", status.Plan[1])
	}
}

func TestOrchestratorService_ValidationErrors(t *testing.T) {
	svc := newTestOrchestratorService()

	tests := []struct {
		name    string
		payload *orchestrator.DeployPayload
		errMsg  string
	}{
		{
			name: "empty payload",
			payload: &orchestrator.DeployPayload{
				Nodes: []orchestrator.DeployNode{},
				Edges: []orchestrator.DeployEdge{},
			},
			errMsg: "at least one node",
		},
		{
			name: "invalid node kind",
			payload: &orchestrator.DeployPayload{
				Nodes: []orchestrator.DeployNode{
					{
						ID:   "node-1",
						Kind: "invalid",
						Name: "test",
					},
				},
			},
			errMsg: "unknown node kind",
		},
		{
			name: "edge to non-existent node",
			payload: &orchestrator.DeployPayload{
				Nodes: []orchestrator.DeployNode{
					{
						ID:   "node-1",
						Kind: orchestrator.NodeKindApp,
						Name: "test",
					},
				},
				Edges: []orchestrator.DeployEdge{
					{
						Source: "node-1",
						Target: "missing",
					},
				},
			},
			errMsg: "unknown node",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.StartDeployment(tt.payload)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
