package orchestrator

import "github.com/SManriqueDev/cubearchitect/internal/cubepath"

type DeployNode struct {
	ID            string            `json:"id"`
	Type          NodeType          `json:"type"`
	Name          string            `json:"name"`
	Label         string            `json:"label,omitempty"`
	Blueprint     string            `json:"blueprint,omitempty"`
	PlanName      string            `json:"plan_name,omitempty"`
	TemplateName  string            `json:"template_name,omitempty"`
	LocationName  string            `json:"location_name,omitempty"`
	IPv4          bool              `json:"ipv4,omitempty"`
	EnableBackups bool              `json:"enable_backups,omitempty"`
	Params        map[string]string `json:"params,omitempty"`
	Status        string            `json:"status,omitempty"`
}

type DeployEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type DeployPayload struct {
	ProjectID   int          `json:"project_id"`
	SSHKeyNames []string     `json:"ssh_key_names,omitempty"`
	Nodes       []DeployNode `json:"nodes"`
	Edges       []DeployEdge `json:"edges"`
}

type DeployResponse struct {
	Success       bool   `json:"success"`
	DeploymentID  string `json:"deployment_id"`
	Message       string `json:"message"`
	NodesCount    int    `json:"nodes_count"`
	EdgesCount    int    `json:"edges_count"`
	ExecutionPlan int    `json:"execution_levels"`
}

type NodeStatus struct {
	NodeID    string
	Status    string
	Error     string
	VPSInfo   *VPSDeploymentInfo
	Timestamp int64
}

type VPSDeploymentInfo struct {
	VPSID            int
	Name             string
	Status           string
	IPAddress        string
	ConnectionString string
	NodeType         string
	OriginalNodeID   string
}

type DeploymentContext struct {
	DeploymentID string
	Client       *cubepath.Client
	ProjectID    int
	SSHKeyNames  []string
	Nodes        map[string]*DeployNode
	Edges        []DeployEdge
	Plan         ExecutionPlan
	NodeStatuses map[string]*NodeStatus
	ErrorsChan   chan error
	StopChan     chan struct{}
}

func NewDeploymentContext(deploymentID string, payload *DeployPayload) *DeploymentContext {
	ctx := &DeploymentContext{
		DeploymentID: deploymentID,
		Nodes:        make(map[string]*DeployNode),
		Edges:        payload.Edges,
		NodeStatuses: make(map[string]*NodeStatus),
		ErrorsChan:   make(chan error, 10),
		StopChan:     make(chan struct{}),
	}

	for i := range payload.Nodes {
		node := &payload.Nodes[i]
		ctx.Nodes[node.ID] = node
		ctx.NodeStatuses[node.ID] = &NodeStatus{
			NodeID: node.ID,
			Status: "pending",
		}
	}

	return ctx
}

func (ctx *DeploymentContext) GetNodeDependencies(nodeID string) []string {
	var deps []string
	for _, edge := range ctx.Edges {
		if edge.Target == nodeID {
			deps = append(deps, edge.Source)
		}
	}
	return deps
}
