package helper

import "slices"

func Contains(xs []string, values ...string) bool {
	for _, value := range values {
		if !slices.Contains(xs, value) {
			return false
		}
	}

	return true
}
