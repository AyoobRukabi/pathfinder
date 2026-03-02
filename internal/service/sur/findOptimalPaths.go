package sur

import (
	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

type Edge struct {
	To, Cap, Flow, Cost, Rev int
}

type Service struct {
	startStation string
	endStation   string
	numTrains    int
	networkMap   domain.MapData
}

func New(start, end string, trains int, data domain.MapData) *Service {
	return &Service{
		startStation: start,
		endStation:   end,
		numTrains:    trains,
		networkMap:   data,
	}
}
func (s *Service) getReindexedMap() ([][]int, []string) {
	oldToNew := make(map[int]int)
	newStationNames := make([]string, len(s.networkMap.Stations))

	//  Get original IDs for Start and End
	startID := s.networkMap.StationsNameToID[s.startStation]
	endID := s.networkMap.StationsNameToID[s.endStation]

	//  Map Start to 0 and End to 1
	oldToNew[startID] = 0
	oldToNew[endID] = 1
	newStationNames[0] = s.startStation
	newStationNames[1] = s.endStation

	//  Map everything else to 2, 3, 4...
	nextID := 2
	for i, station := range s.networkMap.Stations {
		if i == startID || i == endID {
			continue
		}
		oldToNew[i] = nextID
		newStationNames[nextID] = station.Name
		nextID++
	}

	//  Create the new AdjList using the new IDs
	newAdjList := make([][]int, len(s.networkMap.AdjList))
	for oldU, neighbors := range s.networkMap.AdjList {
		for _, oldV := range neighbors {
			newU := oldToNew[oldU]
			newV := oldToNew[oldV]
			newAdjList[newU] = append(newAdjList[newU], newV)
		}
	}

	return newAdjList, newStationNames
}
func (s *Service) FindOptimalPaths() [][]string {
	// We ensure start is 0 and end is 1 as algorithm expects
	adjList, stationNames := s.getReindexedMap()

	// Setup Graph
	numNodes := 2*len(stationNames) - 2
	graph := make([][]Edge, numNodes)
	addEdge := func(from, to, cap, cost int) {
		graph[from] = append(graph[from], Edge{to, cap, 0, cost, len(graph[to])})
		graph[to] = append(graph[to], Edge{from, 0, 0, -cost, len(graph[from]) - 1})
	}

	// Add vertex capacities (split nodes) to enforce vertex-disjoint constraint
	for i := 2; i < len(stationNames); i++ {
		addEdge(inNode(i), outNode(i), 1, 0)
	}

	// Add track connections from AdjList
	for u, neighbors := range adjList {
		for _, v := range neighbors {
			// Since AdjList is undirected (both u-v and v-u exist),
			// we only process each pair once to avoid duplicate edges.
			if u > v {
				continue
			}

			// Logic for connecting nodes based on their roles
			if u == 0 {
				// Waterloo (Source) to any station V
				// We connect to the "In" side of V
				addEdge(0, inNode(v), 1, 1)
			} else if v == 0 {
				// Case where V is Waterloo
				addEdge(0, inNode(u), 1, 1)
			} else if u == 1 {
				// Any station V to Victoria (Sink)
				// We leave the "Out" side of V to hit Victoria
				addEdge(outNode(v), 1, 1, 1)
			} else if v == 1 {
				// Case where V is Victoria
				addEdge(outNode(u), 1, 1, 1)
			} else {
				// Intermediate to Intermediate (e.g., Euston to St Pancras)
				// Always: OutNode(From) -> InNode(To)
				addEdge(outNode(u), inNode(v), 1, 1)
				addEdge(outNode(v), inNode(u), 1, 1)
			}
		}
	}
	bestTurns := int(1e9)
	var bestPaths [][]string
	// Repeatedly find the shortest augmenting path allowing negative edge costs
	for spfa(graph, 0, 1, numNodes) {
		paths := extractPaths(graph, stationNames)
		turns, _ := optimizeTrainAllocation(paths, s.numTrains)

		if turns < bestTurns {
			bestTurns = turns
			bestPaths = paths
		}
	}
	return bestPaths
}
