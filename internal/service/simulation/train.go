package simulation

// Train represents a single train unit
type Train struct {
	ID          int
	PathIDs     []int // The sequence of Station IDs this train must visit
	CurrentIdx  int   // Current index in PathIDs (-1 = waiting at start)
	Arrived     bool
}