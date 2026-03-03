package sur

import (
	"fmt"

	"gitea.kood.tech/ivanandreev/pathfinder/internal/domain"
)

type Service struct {
	startStation string
	endStation   string
	numTrains    int
	networkMap   domain.MapData
	//edge         domain.Edge
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
	fmt.Println(s.networkMap)
	// Setup Graph
	numNodes := 2*len(stationNames) - 2
	graph := make([][]domain.Edge, numNodes)
	addEdge := func(from, to, cap, cost int) {
		graph[from] = append(graph[from], domain.Edge{
			To:       to,
			Cap:      cap,
			Flow:     0,
			Cost:     cost,
			RevIndex: len(graph[to]),
		})
		graph[to] = append(graph[to], domain.Edge{
			To:       from,
			Cap:      0,
			Flow:     0,
			Cost:     -cost,
			RevIndex: len(graph[from]) - 1,
		})
	}
	inNode := func(idx int) int {
		if idx <= 1 {
			return idx
		}
		return 2*idx - 2
	}
	outNode := func(idx int) int {
		if idx <= 1 {
			return idx
		}
		return 2*idx - 1
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

// spfa implements the Shortest Path Faster Algorithm for finding paths with negative costs.
func spfa(graph [][]domain.Edge, S, T, numNodes int) bool {
	dist := make([]int, numNodes)
	for i := range dist {
		dist[i] = 1e9
	}
	parentnode := make([]int, numNodes)
	parentEdge := make([]int, numNodes)
	inQueue := make([]bool, numNodes)

	queue := []int{S}
	dist[S] = 0
	inQueue[S] = true

	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		inQueue[u] = false

		for i, e := range graph[u] {
			if e.Cap-e.Flow > 0 && dist[u]+e.Cost < dist[e.To] {
				dist[e.To] = dist[u] + e.Cost
				parentnode[e.To] = u
				parentEdge[e.To] = i
				if !inQueue[e.To] {
					queue = append(queue, e.To)
					inQueue[e.To] = true
				}
			}
		}
	}

	if dist[T] == 1e9 {
		return false
	}

	// Augment flow along the shortest path
	curr := T
	for curr != S {
		p := parentnode[curr]
		idx := parentEdge[curr]
		revIdx := graph[p][idx].RevIndex

		graph[p][idx].Flow++
		graph[curr][revIdx].Flow--
		curr = p
	}
	return true
}
func extractPaths(graph [][]domain.Edge, names []string) [][]string {
	var paths [][]string
	for _, e := range graph[0] {
		if e.Cap > 0 && e.Flow == 1 {
			path := []string{names[0]}
			curr := e.To

			visited := make(map[int]bool)
			for curr != 1 {
				if visited[curr] {
					break
				}
				visited[curr] = true

				idx := (curr + 2) / 2
				path = append(path, names[idx])

				out := curr + 1
				for _, nextE := range graph[out] {
					if nextE.Cap > 0 && nextE.Flow == 1 {
						curr = nextE.To
						break
					}
				}
			}
			path = append(path, names[1])
			paths = append(paths, path)
		}
	}
	return paths
}

type PathInfo struct {
	Nodes  []string
	Trains int
}

func optimizeTrainAllocation(paths [][]string, numTrains int) (int, []*PathInfo) {
	allocation := make([]int, len(paths))
	trainsLeft := numTrains

	for trainsLeft > 0 {
		bestIdx := 0
		minCost := len(paths[0]) + allocation[0]
		for i := 1; i < len(paths); i++ {
			cost := len(paths[i]) + allocation[i]
			if cost < minCost {
				minCost = cost
				bestIdx = i
			}
		}
		allocation[bestIdx]++
		trainsLeft--
	}

	maxTurns := 0
	var result []*PathInfo
	for i := range paths {
		if allocation[i] > 0 {
			turns := len(paths[i]) - 1 + allocation[i] - 1
			if turns > maxTurns {
				maxTurns = turns
			}
			result = append(result, &PathInfo{Nodes: paths[i], Trains: allocation[i]})
		}
	}
	return maxTurns, result
}
