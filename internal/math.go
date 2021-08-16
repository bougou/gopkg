package internal

// difference represents a - b
func difference(a, b []string) []string {
	// map of b
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func contains(args []string, element string) bool {
	for _, arg := range args {
		if arg == element {
			return true
		}
	}
	return false
}
