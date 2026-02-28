package cli

import (
	"fmt"
	"testing"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

func TestRunner_Start(t *testing.T) {
	// --- MOCK DATA ---
	mockGraph := &domain.Graph{
		Stations: []domain.Station{
			{Name: "start", X: 0, Y: 0},
			{Name: "A", X: 1, Y: 1},
			{Name: "B", X: 2, Y: 2},
			{Name: "end", X: 10, Y: 10},
		},
	}

	// Path 1: start -> A -> end
	// Path 2: start -> B -> end
	mockPaths := [][]int{
		{0, 1, 3},
		{0, 2, 3},
	}

	fmt.Println("--- CLI Runner Integration Test ---")

	// --- TEST EXECUTION ---
	runner := New()
	
	// We just want to ensure it runs without error and prints to stdout
	err := runner.Start(mockGraph, mockPaths, 4)
	
	if err != nil {
		t.Errorf("Runner.Start() failed with error: %v", err)
	}
}