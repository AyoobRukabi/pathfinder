package sur

func extractPaths(graph [][]Edge, names []string) [][]string {
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
