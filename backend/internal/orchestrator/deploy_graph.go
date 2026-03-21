package orchestrator

// GraphStep describes an individual deployment step.
type GraphStep struct {
	Name         string
	Dependencies []string
}

// DeployGraph models a series of steps to deploy resources.
type DeployGraph struct {
	Steps []GraphStep
}
