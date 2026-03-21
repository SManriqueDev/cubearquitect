package orchestrator

import (
	"testing"
)

func TestGeneratePlan_Simple(t *testing.T) {
	services := []Service{
		{ID: "db", DependsOn: []string{}},
		{ID: "app", DependsOn: []string{"db"}},
	}

	plan, err := GeneratePlan(services)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plan) != 2 {
		t.Errorf("expected 2 levels, got %d", len(plan))
	}

	if len(plan[0]) != 1 || plan[0][0] != "db" {
		t.Errorf("level 0 should be [db], got %v", plan[0])
	}

	if len(plan[1]) != 1 || plan[1][0] != "app" {
		t.Errorf("level 1 should be [app], got %v", plan[1])
	}
}

func TestGeneratePlan_Parallel(t *testing.T) {
	services := []Service{
		{ID: "db", DependsOn: []string{}},
		{ID: "cache", DependsOn: []string{}},
		{ID: "app", DependsOn: []string{"db", "cache"}},
	}

	plan, err := GeneratePlan(services)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plan) != 2 {
		t.Errorf("expected 2 levels, got %d", len(plan))
	}

	if len(plan[0]) != 2 {
		t.Errorf("level 0 should have 2 nodes, got %d", len(plan[0]))
	}

	if len(plan[1]) != 1 || plan[1][0] != "app" {
		t.Errorf("level 1 should be [app], got %v", plan[1])
	}
}

func TestGeneratePlan_Circle(t *testing.T) {
	services := []Service{
		{ID: "a", DependsOn: []string{"b"}},
		{ID: "b", DependsOn: []string{"a"}},
	}

	_, err := GeneratePlan(services)
	if err == nil {
		t.Fatal("expected circular dependency error, got nil")
	}
}

func TestGeneratePlan_MissingDependency(t *testing.T) {
	services := []Service{
		{ID: "app", DependsOn: []string{"missing_db"}},
	}

	_, err := GeneratePlan(services)
	if err == nil {
		t.Fatal("expected error for missing dependency, got nil")
	}
}

func TestGeneratePlan_Empty(t *testing.T) {
	plan, err := GeneratePlan([]Service{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plan) != 0 {
		t.Errorf("expected empty plan, got %v", plan)
	}
}

func TestGeneratePlan_ComplexDAG(t *testing.T) {
	// Level 0: db, cache
	// Level 1: app (depends on db, cache)
	// Level 2: worker (depends on app)
	services := []Service{
		{ID: "db", DependsOn: []string{}},
		{ID: "cache", DependsOn: []string{}},
		{ID: "app", DependsOn: []string{"db", "cache"}},
		{ID: "worker", DependsOn: []string{"app"}},
	}

	plan, err := GeneratePlan(services)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(plan) != 3 {
		t.Errorf("expected 3 levels, got %d", len(plan))
	}

	if len(plan[0]) != 2 {
		t.Errorf("level 0 should have 2 nodes, got %d", len(plan[0]))
	}

	if len(plan[1]) != 1 {
		t.Errorf("level 1 should have 1 node, got %d", len(plan[1]))
	}

	if len(plan[2]) != 1 {
		t.Errorf("level 2 should have 1 node, got %d", len(plan[2]))
	}
}
