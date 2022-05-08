package main

func Product(a, b []string) [][]string {
	result := make([][]string, 0, len(a))
	for _, v := range a {
		for _, vv := range b {
			result = append(result, []string{v, vv})
		}
	}
	return result
}
