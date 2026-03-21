package orchestrator

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/SManriqueDev/cubearchitect/internal/config"
	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

const (
	postgresBasicName = "postgres-basic"
	postgresPort      = 5432
	postgresUsername  = "postgres"

	metaKeyDBPassword = "db_password"
)

type PostgresBasicBlueprint struct {
	config *config.Config
}

func NewPostgresBasicBlueprint(cfg *config.Config) *PostgresBasicBlueprint {
	return &PostgresBasicBlueprint{config: cfg}
}

func (bp *PostgresBasicBlueprint) Kind() NodeKind     { return NodeKindDatabase }
func (bp *PostgresBasicBlueprint) Name() string       { return postgresBasicName }
func (bp *PostgresBasicBlueprint) EnvVarName() string { return "DATABASE_URL" }

func (bp *PostgresBasicBlueprint) BuildVPSRequest(nodeID string, params map[string]string) (interface{}, map[string]interface{}, error) {
	dbName := "app_db"
	if custom, ok := params["db_name"]; ok {
		dbName = custom
	}

	password, err := generatePassword(24)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate DB password: %w", err)
	}

	truncatedID := nodeID
	if len(nodeID) > 8 {
		truncatedID = nodeID[:8]
	}

	req := cubepath.VPSCreateRequest{
		Name:            fmt.Sprintf("postgres-%s", truncatedID),
		PlanName:        "gp.nano",
		TemplateName:    "ubuntu-24",
		LocationName:    "us-mia-1",
		Label:           fmt.Sprintf("PostgreSQL (%s)", nodeID),
		IPv4:            true,
		EnableBackups:   true,
		CustomCloudinit: bp.generateCloudInit(dbName, password),
	}

	if bp.config.SSHKeyNames != "" {
		keyNames := strings.Split(bp.config.SSHKeyNames, ",")
		for i := range keyNames {
			keyNames[i] = strings.TrimSpace(keyNames[i])
		}
		req.SSHKeyNames = keyNames
	}

	meta := map[string]interface{}{
		metaKeyDBPassword: password,
	}
	return req, meta, nil
}

func (bp *PostgresBasicBlueprint) ExtractConnectionString(vpsIP string, metadata map[string]interface{}) (string, error) {
	if vpsIP == "" {
		return "", fmt.Errorf("vpsIP cannot be empty")
	}

	password, ok := metadata[metaKeyDBPassword].(string)
	if !ok || password == "" {
		return "", fmt.Errorf("db_password missing or invalid in metadata")
	}

	return fmt.Sprintf("postgresql://%s:%s@%s:%d/app_db?sslmode=disable",
		postgresUsername, password, vpsIP, postgresPort), nil
}

func (bp *PostgresBasicBlueprint) generateCloudInit(dbName, password string) string {
	return fmt.Sprintf(`#cloud-config
package_update: true
packages:
  - postgresql
  - postgresql-contrib

runcmd:
  - systemctl start postgresql
  - systemctl enable postgresql
  - |
    sudo -u postgres psql << 'EOFPSQL'
    CREATE DATABASE %s;
    CREATE USER %s WITH PASSWORD '%s';
    ALTER ROLE %s SET client_encoding TO 'utf8';
    ALTER ROLE %s SET default_transaction_isolation TO 'read committed';
    ALTER ROLE %s SET default_transaction_deferrable TO on;
    GRANT ALL PRIVILEGES ON DATABASE %s TO %s;
    ALTER USER %s CREATEDB;
    EOFPSQL
  - echo "host    all             all             0.0.0.0/0               md5" >> /etc/postgresql/*/main/pg_hba.conf
  - sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" /etc/postgresql/*/main/postgresql.conf
  - systemctl restart postgresql
  - echo "PostgreSQL setup complete"
`, dbName, postgresUsername, password, postgresUsername, postgresUsername, postgresUsername, dbName, postgresUsername, postgresUsername)
}

// generatePassword returns a URL-safe base64-encoded random string of the
// requested byte length (the encoded string will be longer).
func generatePassword(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
