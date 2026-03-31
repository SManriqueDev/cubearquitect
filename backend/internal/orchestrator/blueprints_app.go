package orchestrator

import (
	"fmt"
	"log"
	"strings"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

const (
	nodeBasicName = "node-basic"
	nodePort      = 3000
)

type NodeBasicBlueprint struct {
}

func NewNodeBasicBlueprint() *NodeBasicBlueprint {
	return &NodeBasicBlueprint{}
}

func (bp *NodeBasicBlueprint) Type() NodeType     { return NodeTypeApp }
func (bp *NodeBasicBlueprint) Name() string       { return nodeBasicName }
func (bp *NodeBasicBlueprint) EnvVarName() string { return "APP_URL" }

func (bp *NodeBasicBlueprint) BuildVPSRequest(node *DeployNode, params map[string]string) (interface{}, error) {
	cloudInit := bp.generateCloudInit(params)

	truncatedID := node.ID
	if idx := strings.LastIndex(node.ID, "-"); idx > 0 {
		truncatedID = node.ID[:idx]
	}
	if len(truncatedID) > 50 {
		truncatedID = truncatedID[:50]
	}

	name := node.Name
	if name == "" {
		name = fmt.Sprintf("vps-%s", truncatedID)
	}
	label := node.Label
	if label == "" {
		label = fmt.Sprintf("app Node.js %s", truncatedID)
	}

	planName := getStringParam(params, "plan_name", "gp.nano")
	locationName := getStringParam(params, "location_name", "us-mia-1")

	req := cubepath.VPSCreateRequest{
		Name:            name,
		PlanName:        planName,
		TemplateName:    "ubuntu-24",
		LocationName:    locationName,
		Label:           label,
		IPv4:            getBoolParam(params, "ipv4", false),
		EnableBackups:   getBoolParam(params, "enable_backups", false),
		CustomCloudinit: cloudInit,
	}

	if sshKeys := getStringParam(params, "ssh_key_names", ""); sshKeys != "" {
		keyNames := strings.Split(sshKeys, ",")
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
	var keys []string
	for k := range envVars {
		keys = append(keys, k)
	}
	log.Printf("Generating cloud-init with %d env vars: %v", len(envVars), keys)

	var systemParams = map[string]bool{
		"plan_name":      true,
		"location_name":  true,
		"template_name":  true,
		"ipv4":           true,
		"enable_backups": true,
	}

	var appConfigContent string
	var etcEnvContent string

	for k, v := range envVars {
		if systemParams[k] {
			continue
		}
		escapedVal := strings.ReplaceAll(v, `"`, `\"`)
		appConfigContent += fmt.Sprintf("export %s=\"%s\"\n", k, escapedVal)
		etcEnvContent += fmt.Sprintf("%s=\"%s\"\n", k, escapedVal)
	}

	appConfigContent = strings.TrimRight(appConfigContent, "\n")
	etcEnvContent = strings.TrimRight(etcEnvContent, "\n")

	if appConfigContent == "" {
		appConfigContent = "# No environment variables configured"
		etcEnvContent = ""
	}

	return fmt.Sprintf(`#cloud-config
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
      %s
  - path: /etc/environment
    permissions: '0644'
    content: |
      %s
`, appConfigContent, etcEnvContent)
}
