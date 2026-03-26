package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type NodeTypeStore struct {
	mu       sync.RWMutex
	data     map[int]string
	filePath string
}

func NewNodeTypeStore(baseDir string) (*NodeTypeStore, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	store := &NodeTypeStore{
		data:     make(map[int]string),
		filePath: filepath.Join(baseDir, "node_types.json"),
	}

	if err := store.Load(); err != nil {
		log.Printf("[NodeTypeStore] No existing data file, starting fresh: %v", err)
	}

	return store, nil
}

func (s *NodeTypeStore) Set(vpsID int, nodeType string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[vpsID] = nodeType
	log.Printf("[NodeTypeStore] Registered VPS %d as type: %s", vpsID, nodeType)

	if err := s.saveLocked(); err != nil {
		log.Printf("[NodeTypeStore] Failed to persist: %v", err)
	}
}

func (s *NodeTypeStore) Get(vpsID int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodeType, exists := s.data[vpsID]
	return nodeType, exists
}

func (s *NodeTypeStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var loaded map[int]string
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	s.data = loaded
	log.Printf("[NodeTypeStore] Loaded %d entries from %s", len(s.data), s.filePath)
	return nil
}

func (s *NodeTypeStore) saveLocked() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *NodeTypeStore) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.saveLocked()
}

func (s *NodeTypeStore) Close() error {
	log.Printf("[NodeTypeStore] Closing, final save of %d entries", len(s.data))
	return s.Save()
}

func (s *NodeTypeStore) GetAll() map[int]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[int]string, len(s.data))
	for k, v := range s.data {
		result[k] = v
	}
	return result
}
