package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

type DeploymentEngine struct {
	client    *cubepath.Client
	registry  *BlueprintRegistry
	projectID int
	eventHub  *EventHub
}

func NewDeploymentEngine(client *cubepath.Client, projectID int, registry *BlueprintRegistry) *DeploymentEngine {
	return &DeploymentEngine{
		client:    client,
		projectID: projectID,
		registry:  registry,
	}
}

func (e *DeploymentEngine) SetEventHub(hub *EventHub) {
	e.eventHub = hub
}

func (e *DeploymentEngine) ExecuteDeployment(ctx context.Context, deployCtx *DeploymentContext) error {
	if deployCtx.Plan == nil || len(deployCtx.Plan) == 0 {
		return fmt.Errorf("execution plan is empty")
	}

	for levelIdx, level := range deployCtx.Plan {
		log.Printf("[Deployment %s] Starting level %d with %d nodes: %v", deployCtx.DeploymentID, levelIdx, len(level), level)

		if e.eventHub != nil {
			e.eventHub.Publish(EventLevelStart(deployCtx.DeploymentID, levelIdx))
		}

		if err := e.executeLevel(ctx, deployCtx, level, levelIdx); err != nil {
			log.Printf("[Deployment %s] Level %d failed: %v", deployCtx.DeploymentID, levelIdx, err)
			if e.eventHub != nil {
				e.eventHub.Publish(EventError(deployCtx.DeploymentID, fmt.Sprintf("Level %d failed: %v", levelIdx, err)))
			}
			return err
		}

		log.Printf("[Deployment %s] Level %d completed successfully", deployCtx.DeploymentID, levelIdx)

		if e.eventHub != nil {
			e.eventHub.Publish(EventLevelComplete(deployCtx.DeploymentID, levelIdx))
		}
	}

	return nil
}

func (e *DeploymentEngine) executeLevel(ctx context.Context, deployCtx *DeploymentContext, level []string, levelIdx int) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(level))

	for _, nodeID := range level {
		wg.Add(1) // Increment WaitGroup counter for each node in the level
		go func(nID string) {
			defer wg.Done()

			deps := deployCtx.GetNodeDependencies(nID)
			for _, depID := range deps {
				depStatus := deployCtx.NodeStatuses[depID]
				if depStatus.Status == "error" || depStatus.Status == "cancelled" {
					deployCtx.NodeStatuses[nID].Status = "cancelled"
					deployCtx.NodeStatuses[nID].Error = fmt.Sprintf("dependency %s failed", depID)
					log.Printf("[Deployment %s] Node %s cancelled due to dependency failure", deployCtx.DeploymentID, nID)

					if e.eventHub != nil {
						e.eventHub.Publish(EventFromNodeStatus(deployCtx.DeploymentID, nID, deployCtx.NodeStatuses[nID]))
					}
					return
				}
			}

			if err := e.deployNode(ctx, deployCtx, nID); err != nil {
				deployCtx.NodeStatuses[nID].Status = "error"
				deployCtx.NodeStatuses[nID].Error = err.Error()

				if e.eventHub != nil {
					e.eventHub.Publish(EventFromNodeStatus(deployCtx.DeploymentID, nID, deployCtx.NodeStatuses[nID]))
				}
				errChan <- fmt.Errorf("node %s: %w", nID, err)
			} else {
				deployCtx.NodeStatuses[nID].Status = "healthy"
				log.Printf("[Deployment %s] Node %s deployment completed successfully", deployCtx.DeploymentID, nID)

				if e.eventHub != nil {
					e.eventHub.Publish(EventFromNodeStatus(deployCtx.DeploymentID, nID, deployCtx.NodeStatuses[nID]))
				}
			}
		}(nodeID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *DeploymentEngine) deployNode(ctx context.Context, deployCtx *DeploymentContext, nodeID string) error {
	node := deployCtx.Nodes[nodeID]
	if node == nil {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	deployCtx.NodeStatuses[nodeID].Status = "deploying"

	blueprintName := node.Blueprint
	if blueprintName == "" {
		bp, err := e.registry.GetDefault(node.Kind)
		if err != nil {
			return fmt.Errorf("no default blueprint for kind %s: %w", node.Kind, err)
		}
		blueprintName = bp.Name()
	}

	blueprint, err := e.registry.Get(node.Kind, blueprintName)
	if err != nil {
		return fmt.Errorf("blueprint not found: %s:%s: %w", node.Kind, blueprintName, err)
	}

	vpsReqInterface, err := blueprint.BuildVPSRequest(nodeID, node.Params)
	if err != nil {
		return fmt.Errorf("failed to build VPS request: %w", err)
	}

	vpsReq, ok := vpsReqInterface.(cubepath.VPSCreateRequest)
	if !ok {
		return fmt.Errorf("unexpected VPS request type from blueprint")
	}

	log.Printf("[Deployment %s] Creating VPS for node %s", deployCtx.DeploymentID, nodeID)
	vpsID, vpsIP, err := e.createAndWaitVPS(ctx, vpsReq)
	if err != nil {
		return fmt.Errorf("failed to create VPS: %w", err)
	}

	deployCtx.NodeStatuses[nodeID].VPSInfo = &VPSDeploymentInfo{
		VPSID:     vpsID,
		Name:      vpsReq.Name,
		IPAddress: vpsIP,
	}

	if node.Kind == NodeKindDatabase || node.Kind == NodeKindCache {
		connStr, err := blueprint.ExtractConnectionString(vpsIP, nil)
		if err != nil {
			return fmt.Errorf("failed to extract connection string: %w", err)
		}
		deployCtx.NodeStatuses[nodeID].VPSInfo.ConnectionString = connStr
		log.Printf("[Deployment %s] Node %s connection string extracted: %s", deployCtx.DeploymentID, nodeID, connStr[:20]+"...")

		if err := e.injectConnectionStringToDependents(deployCtx, nodeID, blueprint.EnvVarName(), connStr); err != nil {
			return fmt.Errorf("failed to inject connection string: %w", err)
		}
	}

	return nil
}

func (e *DeploymentEngine) createAndWaitVPS(ctx context.Context, req cubepath.VPSCreateRequest) (int, string, error) {
	respBytes, err := e.client.Post(fmt.Sprintf("/vps/create/%d", e.projectID), req)
	if err != nil {
		return 0, "", fmt.Errorf("CubePath VPS creation failed: %w", err)
	}

	respStr := string(respBytes)
	if len(respStr) > 500 {
		log.Printf("[createAndWaitVPS] API response (first 500 chars): %s...", respStr[:500])
	} else {
		log.Printf("[createAndWaitVPS] API response: %s", respStr)
	}

	createResp := &cubepath.VPSCreateResponse{}
	var vps *cubepath.VPS
	if err := json.Unmarshal(respBytes, createResp); err == nil && createResp.VPSID != 0 {
		log.Printf("[createAndWaitVPS] Parsed CubePath format: VPS ID=%d, Status=%s, IPv4=%s",
			createResp.VPSID, createResp.Status, createResp.IPv4Address)

		vps = &cubepath.VPS{
			ID:     createResp.VPSID,
			Name:   createResp.Name,
			Status: createResp.Status,
			IPv4:   createResp.IPv4Address,
			IPv6:   createResp.IPv6Address,
		}
	} else {
		vps = &cubepath.VPS{}

		if err := json.Unmarshal(respBytes, vps); err == nil && vps.ID != 0 {
			log.Printf("[createAndWaitVPS] Parsed direct format: VPS ID=%d, Status=%s", vps.ID, vps.Status)
			vps.ExtractIPs()
		} else {
			var wrapped struct {
				VPS *cubepath.VPS `json:"vps"`
			}
			if err := json.Unmarshal(respBytes, &wrapped); err == nil && wrapped.VPS != nil && wrapped.VPS.ID != 0 {
				vps = wrapped.VPS
				vps.ExtractIPs()
				log.Printf("[createAndWaitVPS] Parsed wrapped format: VPS ID=%d, Status=%s", vps.ID, vps.Status)
			} else {
				var arr []cubepath.VPS
				if err := json.Unmarshal(respBytes, &arr); err == nil && len(arr) > 0 && arr[0].ID != 0 {
					vps = &arr[0]
					vps.ExtractIPs()
					log.Printf("[createAndWaitVPS] Parsed array format: VPS ID=%d, Status=%s", vps.ID, vps.Status)
				} else {
					return 0, "", fmt.Errorf("failed to parse VPS response in any supported format")
				}
			}
		}
	}

	if vps.ID == 0 {
		return 0, "", fmt.Errorf("received invalid VPS ID (0) from API response")
	}

	if (vps.Status == "running" || vps.Status == "active") && vps.IPv4 != "" {
		log.Printf("[VPS %d] Running/Active immediately with IP %s", vps.ID, vps.IPv4)
		return vps.ID, vps.IPv4, nil
	}

	maxRetries := 120
	retryInterval := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		select {
		case <-ctx.Done():
			return 0, "", ctx.Err()
		default:
		}

		if i%10 == 0 && i > 0 {
			log.Printf("[VPS %d] Polling status (attempt %d/%d, status=%s)", vps.ID, i, maxRetries, vps.Status)
		}

		statusBytes, err := e.client.Get("/vps/")
		if err != nil {
			log.Printf("[VPS %d] Failed to fetch VPS list: %v, retrying...", vps.ID, err)
			time.Sleep(retryInterval)
			continue
		}

		var vpsList []*cubepath.VPS
		if err := json.Unmarshal(statusBytes, &vpsList); err != nil {
			log.Printf("[VPS %d] Failed to parse VPS list: %v, retrying...", vps.ID, err)
			time.Sleep(retryInterval)
			continue
		}

		var updatedVPS *cubepath.VPS
		for _, v := range vpsList {
			if v.ID == vps.ID {
				updatedVPS = v
				break
			}
		}

		if updatedVPS == nil {
			log.Printf("[VPS %d] VPS not found in list, retrying...", vps.ID)
			time.Sleep(retryInterval)
			continue
		}

		updatedVPS.ExtractIPs()
		vps = updatedVPS

		if i%10 == 0 {
			log.Printf("[VPS %d] Current status: %s (attempt %d/%d)", vps.ID, vps.Status, i+1, maxRetries)
		}

		if (vps.Status == "running" || vps.Status == "active") && vps.IPv4 != "" {
			log.Printf("[VPS %d] Ready! IP: %s", vps.ID, vps.IPv4)
			return vps.ID, vps.IPv4, nil
		}

		time.Sleep(retryInterval)
	}

	return 0, "", fmt.Errorf("VPS %d did not reach running state in time (waited %v)", vps.ID, maxRetries*2)
}

func (e *DeploymentEngine) injectConnectionStringToDependents(deployCtx *DeploymentContext, sourceNodeID string, envVarName, connStr string) error {
	for _, edge := range deployCtx.Edges {
		if edge.Source == sourceNodeID {
			targetNodeID := edge.Target
			targetNode := deployCtx.Nodes[targetNodeID]

			if targetNode.Kind != NodeKindApp {
				continue
			}

			if targetNode.Params == nil {
				targetNode.Params = make(map[string]string)
			}
			targetNode.Params[envVarName] = connStr

			log.Printf("[Deployment %s] Injected %s into node %s", deployCtx.DeploymentID, envVarName, targetNodeID)
		}
	}

	return nil
}
