package orchestrator

import (
	"fmt"
	"strings"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

const (
	nodeBasicName = "node-basic"
	nodePort      = 3000
)

type NodeBasicBlueprint struct {
	config *config.Config
}

func NewNodeBasicBlueprint(cfg *config.Config) *NodeBasicBlueprint {
	return &NodeBasicBlueprint{config: cfg}
}

func (bp *NodeBasicBlueprint) Kind() NodeKind     { return NodeKindApp }
func (bp *NodeBasicBlueprint) Name() string       { return nodeBasicName }
func (bp *NodeBasicBlueprint) EnvVarName() string { return "APP_URL" }

func (bp *NodeBasicBlueprint) BuildVPSRequest(nodeID string, params map[string]string) (interface{}, error) {
	cloudInit := bp.generateCloudInit(params)

	truncatedID := nodeID
	if len(nodeID) > 8 {
		truncatedID = nodeID[:8]
	}

	req := cubepath.VPSCreateRequest{
		Name:            fmt.Sprintf("node-app-%s", truncatedID),
		PlanName:        "gp.nano",
		TemplateName:    "ubuntu-24",
		LocationName:    "us-mia-1",
		Label:           fmt.Sprintf("Node.js App (%s)", nodeID),
		IPv4:            true,
		EnableBackups:   false,
		CustomCloudinit: cloudInit,
	}

	if bp.config.SSHKeyNames != "" {
		keyNames := strings.Split(bp.config.SSHKeyNames, ",")
		for i := range keyNames {
			keyNames[i] = strings.TrimSpace(keyNames[i])
		}
		req.SSHKeyNames = keyNames
	}

	return req, nil
}

func (bp *NodeBasicBlueprint) ExtractConnectionString(vpsIP string, _ map[string]interface{}) (string, error) {
	if vpsIP == "" {
		return "", fmt.Errorf("vpsIP cannot be empty")
	}
	return fmt.Sprintf("http://%s:%d", vpsIP, nodePort), nil
}

func (bp *NodeBasicBlueprint) generateCloudInit(envVars map[string]string) string {
	var envSection string
	if len(envVars) > 0 {
		envSection = "\nruncmd:\n"
		envSection += "  - echo 'export NODE_ENV=production' >> /root/.bashrc\n"
		for k, v := range envVars {
			escapedVal := strings.ReplaceAll(v, "'", "'\\''")
			envSection += fmt.Sprintf("  - echo 'export %s=%q' >> /root/.bashrc\n", k, escapedVal)
		}
		envSection += "  - echo 'Starting Node.js application...' >> /var/log/app-init.log\n"
		envSection += "  - echo 'Run: npm install && npm start' >> /var/log/app-init.log\n"
	}

	return `#cloud-config
package_update: true
packages:
  - curl
  - wget
  - git
  - build-essential
  - python3
  - nodejs
  - npm
write_files:
  - path: /root/.app-config
    permissions: '0644'
    content: |
      # Application configuration
` + envSection
}
