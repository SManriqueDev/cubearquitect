package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

// SSHKey represents an SSH key from CubePath.
type SSHKey struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	KeyType     string `json:"key_type"`
	Fingerprint string `json:"fingerprint"`
}

// SSHKeysResponse represents the list of SSH keys.
type SSHKeysResponse struct {
	SSHKeys []SSHKey `json:"sshkeys"`
}

// SSHKeysService handles SSH key operations.
type SSHKeysService struct{}

func NewSSHKeysService() *SSHKeysService {
	return &SSHKeysService{}
}

// List retrieves SSH keys from CubePath using the provided client.
func (s *SSHKeysService) List(client *cubepath.Client) ([]SSHKey, error) {
	res, err := client.Get("/sshkey/user/sshkeys")
	if err != nil {
		log.Printf("Error fetching SSH keys: %v", err)
		return nil, err
	}

	var keys SSHKeysResponse
	if err := json.Unmarshal(res, &keys); err != nil {
		log.Printf("Error unmarshaling SSH keys: %v", err)
		return nil, err
	}

	return keys.SSHKeys, nil
}
