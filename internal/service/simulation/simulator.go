package simulation

import (
	"fmt"
	"strings"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

// Simulator controls the movement of trains
type Simulator struct {
	Graph       *domain.Graph // Reference to Ivan's data to get Names
	Trains      []*Train
	TotalTrains int
}

// New creates a simulator instance
// paths: A slice of paths, where each path is a slice of station IDs (e.g. [0, 5, 2])
func New(g *domain.Graph, paths [][]int, numTrains int) *Simulator {
	trains := make([]*Train, numTrains)
	
	// LOAD BALANCING LOGIC
	// We need to decide which train takes which path.
	// For now, we use a simple Round-Robin. 
	// (Anwar might provide a better distribution logic later, but this works for now).
	for i := 0; i < numTrains; i++ {
		// Pick a path in a loop: Path A, Path B, Path A, Path B...
		pathIndex := i % len(paths) 
		
		trains[i] = &Train{
			ID:         i + 1,
			PathIDs:    paths[pathIndex],
			CurrentIdx: 0, // Start at the first station in the path
			Arrived:    false,
		}
	}

	return &Simulator{
		Graph:       g,
		Trains:      trains,
		TotalTrains: numTrains,
	}
}

// Run executes the simulation turn by turn
func (s *Simulator) Run() {
	finishedCount := 0

	// --- THE GAME LOOP (Tick) ---
	for finishedCount < s.TotalTrains {
		
		// 1. Reset occupied map for this turn
		// We use a map to ensure no two trains jump to the same station in one turn
		occupied := make(map[int]bool)
		
		// Collect output strings for this turn (e.g. "T1-waterloo")
		var turnOutput []string

		// 2. Iterate through all trains
		for _, t := range s.Trains {
			if t.Arrived {
				continue
			}

			// Determine where the train wants to go next
			nextStepIdx := t.CurrentIdx + 1
			
			// Safety Check: If path ended
			if nextStepIdx >= len(t.PathIDs) {
				t.Arrived = true
				finishedCount++
				continue
			}

			nextStationID := t.PathIDs[nextStepIdx]

			// --- COLLISION CHECK ---
			// A train can move IF:
			// 1. The destination station is NOT occupied in this turn
			// 2. EXCEPT if it is the End Station (usually infinite capacity).
			//    We assume the last station in the path is the End Station.
			
			isEndStation := (nextStepIdx == len(t.PathIDs)-1)
			
			if !occupied[nextStationID] || isEndStation {
				// MOVE SUCCESSFUL
				t.CurrentIdx = nextStepIdx
				
				// Mark station as occupied (unless it's the infinite end station)
				if !isEndStation {
					occupied[nextStationID] = true
				}

				// Format Output: We need the NAME, not the ID
				// We look it up in Ivan's Graph
				stationName := s.Graph.Stations[nextStationID].Name
				turnOutput = append(turnOutput, fmt.Sprintf("T%d-%s", t.ID, stationName))

				// Check if this move finished the train
				if isEndStation {
					t.Arrived = true
					finishedCount++
				}
			} else {
				// MOVE FAILED (Traffic Jam)
				// The train waits at its CURRENT station.
				// We must mark the CURRENT station as occupied so nobody hits us from behind.
				if t.CurrentIdx != -1 {
					currentStationID := t.PathIDs[t.CurrentIdx]
					occupied[currentStationID] = true
				}
			}
		}

		// 3. Print the turn result
		if len(turnOutput) > 0 {
			fmt.Println(strings.Join(turnOutput, " "))
		}
	}
}