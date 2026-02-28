package simulation

import (
	"fmt"
	"testing"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

func TestSmartRouting(t *testing.T) {
	// MOCK GRAPH
	mockGraph := &domain.Graph{
		Stations: []domain.Station{
			{Name: "start", X: 0, Y: 0},
			{Name: "short1", X: 1, Y: 1},
			{Name: "long1", X: 2, Y: 2},
			{Name: "long2", X: 3, Y: 3},
			{Name: "long3", X: 4, Y: 4},
			{Name: "end", X: 10, Y: 10},
		},
	}

	// MOCK PATHS
	// Path 0: start -> short1 -> end (Length 3)
	// Path 1: start -> long1 -> long2 -> long3 -> end (Length 5)
	mockPaths := [][]int{
		{0, 1, 5},       // Short
		{0, 2, 3, 4, 5}, // Long
	}

	fmt.Println("--- Testing Smart Routing (5 Trains) ---")
	// We expect the first few trains to take Path 0 (Short).
	// Only when Path 0 gets crowded should a train take Path 1 (Long).
	
	sim := New(mockGraph, mockPaths, 5)
	
	// Print assignments to verify logic
	for _, tr := range sim.Trains {
		pathName := "Short"
		if len(tr.PathIDs) > 3 {
			pathName = "Long"
		}
		fmt.Printf("Train %d assigned to: %s Path\n", tr.ID, pathName)
	}

	fmt.Println("\n--- Simulation Output ---")
	sim.Run()
}