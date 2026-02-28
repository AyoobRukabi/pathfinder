package simulation

import (
	"fmt"
	"testing"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

// TestRunSimulation acts as a manual integration test to see the simulator in action
func TestRunSimulation(t *testing.T) {
	// --- MOCK DATA (Fake Graph) ---
	mockGraph := &domain.Graph{
		Stations: []domain.Station{
			{Name: "start", X: 0, Y: 0},   // ID 0
			{Name: "jungle", X: 5, Y: 5},  // ID 1
			{Name: "desert", X: 5, Y: 5},  // ID 2
			{Name: "end", X: 10, Y: 10},   // ID 3
		},
	}

	// --- MOCK PATHS (Fake Algorithms) ---
	// Path 1: start(0) -> jungle(1) -> end(3)
	// Path 2: start(0) -> desert(2) -> end(3)
	mockPaths := [][]int{
		{0, 1, 3}, // Short path 1
		{0, 2, 3}, // Short path 2
	}

	fmt.Println("--- Starting Simulation Test (4 Trains) ---")

	// --- RUN SIMULATOR ---
	sim := New(mockGraph, mockPaths, 4)
	
	// We run it. In a real test, we might check output, 
	// but for now, we just want to see it print to stdout.
	sim.Run()
}