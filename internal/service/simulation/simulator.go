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
	
	// Track how many trains are assigned to each path to calculate waiting time
	pathUsage := make([]int, len(paths))

	for i := 0; i < numTrains; i++ {
		bestPathIndex := 0
		minCost := 999999999 // Start with a high number

		// Find the path with the lowest effective cost (Length + Queue)
		for pIdx, path := range paths {
			// Cost = Length of Path + Number of trains already queued on it
			// This estimates when this specific train would arrive
			cost := len(path) + pathUsage[pIdx]

			if cost < minCost {
				minCost = cost
				bestPathIndex = pIdx
			}
		}

		// Assign train to the best path
		trains[i] = &Train{
			ID:         i + 1,
			PathIDs:    paths[bestPathIndex],
			CurrentIdx: 0, // Starts at the Start Station (index 0)
			Arrived:    false,
		}
		
		// Increment usage for that path
		pathUsage[bestPathIndex]++
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