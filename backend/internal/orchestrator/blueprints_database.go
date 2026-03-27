package orchestrator

import (
	"fmt"
	"strings"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

const (
	postgresBasicName = "postgres-basic"
	postgresPort      = 5432
	postgresUsername  = "postgres"
	postgresPassword  = "postgres123!" // TODO: use secrets management
)

type PostgresBasicBlueprint struct {
	config *config.Config
}

func NewPostgresBasicBlueprint(cfg *config.Config) *PostgresBasicBlueprint {
	return &PostgresBasicBlueprint{config: cfg}
}

func (bp *PostgresBasicBlueprint) Type() NodeType     { return NodeTypeDatabase }
func (bp *PostgresBasicBlueprint) Name() string       { return postgresBasicName }
func (bp *PostgresBasicBlueprint) EnvVarName() string { return "DATABASE_URL" }

func (bp *PostgresBasicBlueprint) BuildVPSRequest(node *DeployNode, params map[string]string) (interface{}, error) {
	dbName := "app_db"
	if custom, ok := params["db_name"]; ok {
		dbName = custom
	}

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
		label = fmt.Sprintf("database PostgreSQL %s", truncatedID)
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
		CustomCloudinit: bp.generateCloudInit(dbName),
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

func (bp *PostgresBasicBlueprint) ExtractConnectionString(vpsIP string, _ map[string]interface{}) (string, error) {
	if vpsIP == "" {
		return "", fmt.Errorf("vpsIP cannot be empty")
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/app_db?sslmode=disable",
		postgresUsername, postgresPassword, vpsIP, postgresPort), nil
}

func (bp *PostgresBasicBlueprint) generateCloudInit(dbName string) string {
	cloudInit := fmt.Sprintf(`#cloud-config
package_update: true
packages:
  - postgresql
  - postgresql-contrib

runcmd:
  - systemctl start postgresql
  - echo "[1/7] PostgreSQL started"
  - systemctl enable postgresql
  - echo "[2/7] PostgreSQL enabled"
  - sleep 3
  - echo "[3/7] Creating database %s..."
  - sudo -u postgres psql -c "CREATE DATABASE %s;" || echo "Database may already exist"
  - echo "[4/7] Creating user %s..."
  - sudo -u postgres psql -c "CREATE USER %s WITH PASSWORD '%s';" || echo "User may already exist"
  - echo "[5/7] Granting privileges..."
  - sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE %s TO %s;" || echo "Privileges may already granted"
  - echo "[6/7] Configuring remote access..."
  - echo "host all all 0.0.0.0/0 md5" >> /etc/postgresql/*/main/pg_hba.conf
  - sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /etc/postgresql/*/main/postgresql.conf || true
  - echo "[7/7] Restarting PostgreSQL..."
  - systemctl restart postgresql
  - echo "PostgreSQL setup complete"
`, dbName, dbName, postgresUsername, postgresUsername, postgresPassword, dbName, postgresUsername)
	return cloudInit
}
