package sur

func inNode(idx int) int {
	if idx <= 1 {
		return idx
	}
	return 2*idx - 2
}
func outNode(idx int) int {
	if idx <= 1 {
		return idx
	}
	return 2*idx - 1
}
