package sur

func spfa(graph [][]Edge, S, T, numNodes int) bool {
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
		revIdx := graph[p][idx].Rev

		graph[p][idx].Flow++
		graph[curr][revIdx].Flow--
		curr = p
	}
	return true
}
