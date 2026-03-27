package orchestrator

import "fmt"

type NodeType string

const (
	NodeTypeApp      NodeType = "app"
	NodeTypeDatabase NodeType = "database"
	NodeTypeCache    NodeType = "cache"
)

type Blueprint interface {
	Type() NodeType
	Name() string
	BuildVPSRequest(node *DeployNode, params map[string]string) (interface{}, error)
	ExtractConnectionString(vpsIP string, metadata map[string]interface{}) (string, error)
	EnvVarName() string
}

type BlueprintRegistry struct {
	blueprints map[string]Blueprint
}

func NewBlueprintRegistry() *BlueprintRegistry {
	return &BlueprintRegistry{blueprints: make(map[string]Blueprint)}
}

func (r *BlueprintRegistry) Register(bp Blueprint) error {
	key := fmt.Sprintf("%s:%s", bp.Type(), bp.Name())
	if _, exists := r.blueprints[key]; exists {
		return fmt.Errorf("blueprint already registered: %s", key)
	}
	r.blueprints[key] = bp
	return nil
}

func (r *BlueprintRegistry) Get(nodeType NodeType, name string) (Blueprint, error) {
	key := fmt.Sprintf("%s:%s", nodeType, name)
	bp, exists := r.blueprints[key]
	if !exists {
		return nil, fmt.Errorf("blueprint not found: %s", key)
	}
	return bp, nil
}

func (r *BlueprintRegistry) GetDefault(nodeType NodeType) (Blueprint, error) {
	for _, bp := range r.blueprints {
		if bp.Type() == nodeType {
			return bp, nil
		}
	}
	return nil, fmt.Errorf("no blueprint found for type: %s", nodeType)
}

func (r *BlueprintRegistry) ListByType(nodeType NodeType) []Blueprint {
	var result []Blueprint
	for _, bp := range r.blueprints {
		if bp.Type() == nodeType {
			result = append(result, bp)
		}
	}
	return result
}
