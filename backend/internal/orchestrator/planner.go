package orchestrator

// Planner defines how to compose a DeployGraph.
type Planner interface {
	Plan() (*DeployGraph, error)
}
