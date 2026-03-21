package orchestrator

import (
	"strings"
	"testing"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

func TestNodeBasicBlueprint_Kind(t *testing.T) {
	cfg := &config.Config{}
	bp := NewNodeBasicBlueprint(cfg)

	if bp.Kind() != NodeKindApp {
		t.Errorf("expected kind %s, got %s", NodeKindApp, bp.Kind())
	}
}

func TestNodeBasicBlueprint_Name(t *testing.T) {
	cfg := &config.Config{}
	bp := NewNodeBasicBlueprint(cfg)

	if bp.Name() != nodeBasicName {
		t.Errorf("expected name %s, got %s", nodeBasicName, bp.Name())
	}
}

func TestNodeBasicBlueprint_BuildVPSRequest(t *testing.T) {
	cfg := &config.Config{
		SSHKeyNames: "TestKey",
	}
	bp := NewNodeBasicBlueprint(cfg)

	req, err := bp.BuildVPSRequest("app-1", map[string]string{
		"DATABASE_URL": "postgresql://user:pass@db-ip:5432/app_db",
	})
	if err != nil {
		t.Fatalf("BuildVPSRequest failed: %v", err)
	}

	vpsReq, ok := req.(cubepath.VPSCreateRequest)
	if !ok {
		t.Fatalf("expected VPSCreateRequest, got %T", req)
	}

	// Check basic properties
	if !strings.Contains(vpsReq.Name, "node-app") {
		t.Errorf("expected name to contain 'node-app', got %s", vpsReq.Name)
	}

	if vpsReq.PlanName != "gp.nano" {
		t.Errorf("expected plan gp.nano, got %s", vpsReq.PlanName)
	}

	if vpsReq.TemplateName != "ubuntu-24" {
		t.Errorf("expected template ubuntu-24, got %s", vpsReq.TemplateName)
	}

	// Check cloud-init starts with #cloud-config
	if !strings.HasPrefix(vpsReq.CustomCloudinit, "#cloud-config") {
		t.Errorf("cloud-init must start with #cloud-config, got: %s", vpsReq.CustomCloudinit)
	}

	// Check SSH keys are present
	if len(vpsReq.SSHKeyNames) == 0 {
		t.Errorf("expected SSH keys to be set")
	}
}

func TestNodeBasicBlueprint_ExtractConnectionString(t *testing.T) {
	cfg := &config.Config{}
	bp := NewNodeBasicBlueprint(cfg)

	connStr, err := bp.ExtractConnectionString("192.168.1.100", nil)
	if err != nil {
		t.Fatalf("ExtractConnectionString failed: %v", err)
	}

	expectedURL := "http://192.168.1.100:3000"
	if connStr != expectedURL {
		t.Errorf("expected %s, got %s", expectedURL, connStr)
	}
}

func TestNodeBasicBlueprint_EnvVarName(t *testing.T) {
	cfg := &config.Config{}
	bp := NewNodeBasicBlueprint(cfg)

	envVar := bp.EnvVarName()
	if envVar != "APP_URL" {
		t.Errorf("expected APP_URL, got %s", envVar)
	}
}
