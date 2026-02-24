package domain

type Station struct {
	Name string
	X    float64
	Y    float64
}

type Edge struct {
	To     int
	Weight float64
}

// option 1
type Graph struct {
	Stations []Station      // Index is the station ID
	NameToID map[string]int // Translates "waterloo" -> 0
	Edges    [][]Edge       // Adjacency list: Edges[0] = []int{1, 2}
}

// option 2

// type Edge struct {
// 	To     string
// 	Weight float64
// }

// type Graph struct {
// 	Stations map[string]Station
// 	Edges    map[string][]Edge // e.g., Edges["waterloo"] = []Edge{"victoria", "euston"}
// }
