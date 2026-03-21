package orchestrator

import (
	"errors"
	"fmt"
)

// Service represents a node in the deployment graph.
type Service struct {
	ID        string
	DependsOn []string
}

// ExecutionPlan represents levels of parallelizable deployment.
type ExecutionPlan [][]string

// GeneratePlan uses Kahn's algorithm to create an execution plan by levels.
func GeneratePlan(services []Service) (ExecutionPlan, error) {
	if len(services) == 0 {
		return ExecutionPlan{}, nil
	}

	// Build adjacency list and in-degree map
	adj := make(map[string][]string)
	inDegree := make(map[string]int)
	allNodes := make(map[string]bool)

	for _, s := range services {
		allNodes[s.ID] = true
		if _, exists := inDegree[s.ID]; !exists {
			inDegree[s.ID] = 0
		}
		for _, dep := range s.DependsOn {
			adj[dep] = append(adj[dep], s.ID)
			inDegree[s.ID]++
		}
	}

	// Validate all dependencies exist
	for _, s := range services {
		for _, dep := range s.DependsOn {
			if !allNodes[dep] {
				return nil, fmt.Errorf("service %s depends on non-existent service %s", s.ID, dep)
			}
		}
	}

	var plan ExecutionPlan
	visitedCount := 0

	// Kahn's algorithm by levels
	for len(allNodes) > visitedCount {
		var currentLevel []string

		// Find nodes with in-degree 0 (no pending dependencies)
		for id := range allNodes {
			if inDegree[id] == 0 {
				currentLevel = append(currentLevel, id)
			}
		}

		// If no nodes with in-degree 0 remain and we haven't visited all nodes, there's a cycle
		if len(currentLevel) == 0 {
			return nil, errors.New("circular dependency detected in architecture")
		}

		plan = append(plan, currentLevel)

		// Mark nodes as processed and update in-degrees of neighbors
		for _, id := range currentLevel {
			for _, neighbor := range adj[id] {
				inDegree[neighbor]--
			}
			inDegree[id] = -1 // Mark as processed
			visitedCount++
		}
	}

	return plan, nil
}
