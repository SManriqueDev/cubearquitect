package orchestrator

import (
"fmt"
"strings"
"testing"

"github.com/SManriqueDev/cubearchitect/internal/config"
"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

func TestPostgresBlueprint_SSHKeyInjection(t *testing.T) {
tests := []struct {
name           string
sshKeyNames    string
expectedKeys   []string
expectedLen    int
}{
{
name:           "single SSH key",
sshKeyNames:    "SebastiansMacPro",
expectedKeys:   []string{"SebastiansMacPro"},
expectedLen:    1,
},
{
name:           "multiple SSH keys",
sshKeyNames:    "Key1,Key2,Key3",
expectedKeys:   []string{"Key1", "Key2", "Key3"},
expectedLen:    3,
},
{
name:           "SSH keys with spaces",
sshKeyNames:    "Key1 , Key2 , Key3",
expectedKeys:   []string{"Key1", "Key2", "Key3"},
expectedLen:    3,
},
{
name:           "no SSH keys",
sshKeyNames:    "",
expectedKeys:   nil,
expectedLen:    0,
},
}

for _, tt := range tests {
t.Run(tt.name, func(t *testing.T) {
cfg := &config.Config{
SSHKeyNames: tt.sshKeyNames,
}
bp := NewPostgresBasicBlueprint(cfg)

req, _, err := bp.BuildVPSRequest("db-test", map[string]string{})
if err != nil {
t.Fatalf("failed to build request: %v", err)
}

vpsReq, ok := req.(cubepath.VPSCreateRequest)
if !ok {
t.Fatalf("expected VPSCreateRequest, got %T", req)
}

// Verify SSH keys
if len(vpsReq.SSHKeyNames) != tt.expectedLen {
t.Errorf("expected %d SSH keys, got %d", tt.expectedLen, len(vpsReq.SSHKeyNames))
}

for i, expectedKey := range tt.expectedKeys {
if i >= len(vpsReq.SSHKeyNames) {
t.Errorf("missing key at index %d: expected %s", i, expectedKey)
continue
}
if vpsReq.SSHKeyNames[i] != expectedKey {
t.Errorf("key %d: expected %s, got %s", i, expectedKey, vpsReq.SSHKeyNames[i])
}
}

// Verify other fields are set correctly
if vpsReq.Name == "" {
t.Error("VPS name should not be empty")
}
if vpsReq.PlanName != "gp.nano" {
t.Errorf("expected plan gp.nano, got %s", vpsReq.PlanName)
}
if vpsReq.TemplateName != "ubuntu-24" {
t.Errorf("expected template ubuntu-24, got %s", vpsReq.TemplateName)
}
if !vpsReq.IPv4 {
t.Error("IPv4 should be enabled")
}
if !vpsReq.EnableBackups {
t.Error("Backups should be enabled")
}
if vpsReq.CustomCloudinit == "" {
t.Error("cloud-init script should not be empty")
}
})
}
}

func TestPostgresBlueprint_ConnectionString(t *testing.T) {
cfg := &config.Config{
SSHKeyNames: "TestKey",
}
bp := NewPostgresBasicBlueprint(cfg)

vpsIP := "192.168.1.100"

// Build a VPS request to obtain the per-deployment metadata (which carries
// the generated password).
_, meta, err := bp.BuildVPSRequest("db-test", map[string]string{})
if err != nil {
t.Fatalf("failed to build VPS request: %v", err)
}

connStr, err := bp.ExtractConnectionString(vpsIP, meta)
if err != nil {
t.Fatalf("failed to extract connection string: %v", err)
}

// Verify the connection string has the expected format.
prefix := fmt.Sprintf("postgresql://%s:", postgresUsername)
if !strings.HasPrefix(connStr, prefix) {
t.Errorf("connection string missing expected prefix %q, got %q", prefix, connStr)
}
suffix := fmt.Sprintf("@%s:%d/app_db?sslmode=disable", vpsIP, postgresPort)
if !strings.HasSuffix(connStr, suffix) {
t.Errorf("connection string missing expected suffix %q, got %q", suffix, connStr)
}

// Verify the password is not the old hardcoded value.
if strings.Contains(connStr, "postgres123!") {
t.Error("connection string must not contain the old hardcoded password")
}

// Verify env var name.
if bp.EnvVarName() != "DATABASE_URL" {
t.Errorf("expected env var name DATABASE_URL, got %s", bp.EnvVarName())
}
}

func TestPostgresBlueprint_UniquePasswords(t *testing.T) {
cfg := &config.Config{}
bp := NewPostgresBasicBlueprint(cfg)

_, meta1, err := bp.BuildVPSRequest("node-1", map[string]string{})
if err != nil {
t.Fatalf("first BuildVPSRequest failed: %v", err)
}
_, meta2, err := bp.BuildVPSRequest("node-2", map[string]string{})
if err != nil {
t.Fatalf("second BuildVPSRequest failed: %v", err)
}

pw1, ok1 := meta1[metaKeyDBPassword].(string)
pw2, ok2 := meta2[metaKeyDBPassword].(string)

if !ok1 || !ok2 {
t.Fatal("metaKeyDBPassword is missing or not a string in build metadata")
}
if pw1 == "" || pw2 == "" {
t.Error("generated passwords must not be empty")
}
if pw1 == pw2 {
t.Error("each deployment must receive a unique password")
}
}

func TestPostgresBlueprint_ExtractConnectionStringMissingMetadata(t *testing.T) {
cfg := &config.Config{}
bp := NewPostgresBasicBlueprint(cfg)

_, err := bp.ExtractConnectionString("10.0.0.1", nil)
if err == nil {
t.Error("expected error when metadata is nil, got nil")
}

_, err = bp.ExtractConnectionString("10.0.0.1", map[string]interface{}{})
if err == nil {
t.Error("expected error when metadata has no db_password, got nil")
}
}

func TestPostgresBlueprint_CloudInitFormat(t *testing.T) {
cfg := &config.Config{
SSHKeyNames: "TestKey",
}
bp := NewPostgresBasicBlueprint(cfg)

req, _, err := bp.BuildVPSRequest("db-test", map[string]string{})
if err != nil {
t.Fatalf("failed to build request: %v", err)
}

vpsReq, ok := req.(cubepath.VPSCreateRequest)
if !ok {
t.Fatalf("expected VPSCreateRequest, got %T", req)
}

// Verify cloud-init script starts with #cloud-config (required by CubePath)
if !strings.HasPrefix(vpsReq.CustomCloudinit, "#cloud-config") {
t.Errorf("cloud-init script must start with '#cloud-config', got first 50 chars: %s",
vpsReq.CustomCloudinit[:50])
}

// Verify script contains PostgreSQL setup
if !strings.Contains(vpsReq.CustomCloudinit, "postgresql") {
t.Error("cloud-init script should contain PostgreSQL setup")
}

// Verify script is not empty
if len(vpsReq.CustomCloudinit) < 100 {
t.Error("cloud-init script should not be too short")
}

// Verify the old hardcoded password is not present in the cloud-init script.
if strings.Contains(vpsReq.CustomCloudinit, "postgres123!") {
t.Error("cloud-init script must not contain the old hardcoded password")
}
}
