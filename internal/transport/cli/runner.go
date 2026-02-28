package cli

import (
	"fmt"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
	"gitea.kood.tech/ivanandreev/pathfinder/internal/service/simulation"
)

// Runner handles the command line visualization
type Runner struct {
	// We can add configuration flags here later if needed
}

// New creates a new CLI runner
func New() *Runner {
	return &Runner{}
}

// Start initializes and runs the simulation
// This is the function Ivan will call from internal/app/app.go
func (r *Runner) Start(graph *domain.Graph, paths [][]int, numTrains int) error {
	// 1. Validation
	if len(paths) == 0 {
		return fmt.Errorf("no paths provided for simulation")
	}
	if numTrains <= 0 {
		return fmt.Errorf("invalid number of trains: %d", numTrains)
	}

	// 2. Initialize the Service
	// We inject the dependencies (Graph and Paths) into the simulator
	sim := simulation.New(graph, paths, numTrains)

	// 3. Run the "Game Loop"
	// The simulator handles the printing to stdout
	sim.Run()

	return nil
}