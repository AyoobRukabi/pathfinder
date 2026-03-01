package domain

type Station struct {
	Name string
	X    int
	Y    int
}

type Edge struct {
	To       int
	Cap      int
	Flow     int
	Cost     int
	RevIndex int
}

type MapData struct {
	Stations         []Station
	StationsNameToID map[string]int
	AdjList          [][]int
}

// option 2

// type Edge struct {
// 	To       int // Destination node ID
// 	Cap      int // 1 for normal stations, math.MaxInt32 for start/end
// 	Flow     int // Current flow (0 or 1)
// 	RevIndex int // The index of the reverse edge in the 'To' node's adjacency list
// }

// // Graph holds the 2N nodes required for vertex splitting
// type Graph struct {
// 	Nodes [][]Edge // Size is 2 * Number of Stations
// }
