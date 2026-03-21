package orchestrator

import "fmt"

type NodeKind string

const (
	NodeKindApp      NodeKind = "app"
	NodeKindDatabase NodeKind = "database"
	NodeKindCache    NodeKind = "cache"
)

type Blueprint interface {
	Kind() NodeKind
	Name() string
	BuildVPSRequest(nodeID string, params map[string]string) (interface{}, error)
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
	key := fmt.Sprintf("%s:%s", bp.Kind(), bp.Name())
	if _, exists := r.blueprints[key]; exists {
		return fmt.Errorf("blueprint already registered: %s", key)
	}
	r.blueprints[key] = bp
	return nil
}

func (r *BlueprintRegistry) Get(kind NodeKind, name string) (Blueprint, error) {
	key := fmt.Sprintf("%s:%s", kind, name)
	bp, exists := r.blueprints[key]
	if !exists {
		return nil, fmt.Errorf("blueprint not found: %s", key)
	}
	return bp, nil
}

func (r *BlueprintRegistry) GetDefault(kind NodeKind) (Blueprint, error) {
	for _, bp := range r.blueprints {
		if bp.Kind() == kind {
			return bp, nil
		}
	}
	return nil, fmt.Errorf("no blueprint found for kind: %s", kind)
}

func (r *BlueprintRegistry) ListByKind(kind NodeKind) []Blueprint {
	var result []Blueprint
	for _, bp := range r.blueprints {
		if bp.Kind() == kind {
			result = append(result, bp)
		}
	}
	return result
}
