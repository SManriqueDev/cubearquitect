package orchestrator

import (
	"testing"
)

func TestGeneratePlanWithIndependentServices(t *testing.T) {
	services := []Service{
		{ID: "1", DependsOn: []string{}},
		{ID: "2", DependsOn: []string{}},
	}
	got, err := GeneratePlan(services)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 levels, got %d", len(got))
	}

	if len(got[0]) != 2 {
		t.Fatalf("expected 2 services in level 0, got %d", len(got[0]))
	}
}

func TestGeneratePlanWithLinearDependencies(t *testing.T) {
	services := []Service{
		{ID: "1", DependsOn: []string{}},
		{ID: "2", DependsOn: []string{"1"}},
	}
	got, err := GeneratePlan(services)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 levels, got %d", len(got))
	}

	if got[0][0] != "1" {
		t.Fatalf("level 0: expected [1], got %v", got[0])
	}

	if got[1][0] != "2" {
		t.Errorf("level 1: expected [2], got %v", got[1])
	}
}

func TestGeneratePlanWithEmptyServices(t *testing.T) {
	services := []Service{}
	got, err := GeneratePlan(services)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected 0 levels, got %d", len(got))
	}
}

func TestGeneratePlanWithCircularDependency(t *testing.T) {
	services := []Service{
		{ID: "1", DependsOn: []string{"2"}},
		{ID: "2", DependsOn: []string{"1"}},
	}
	got, err := GeneratePlan(services)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	expectedErrorMsg := "circular dependency detected in architecture"
	if err.Error() != expectedErrorMsg {
		t.Fatalf("expected error message: %s, got %s", expectedErrorMsg, err.Error())
	}

	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}
