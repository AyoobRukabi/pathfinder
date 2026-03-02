package sur

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
