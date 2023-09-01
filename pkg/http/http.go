package http

import (
	"math/rand"
	"slices"
)

const (
	minPrivatePort = 49152
	maxPrivatePort = 65535
)

// SuggestPortInRangeExcluding will suggest a http port in the given range, excluding the given ports.
func SuggestPortInRangeExcluding(minPort, maxPort int, exclude []int) int {
	availablePorts := []int{}

	for port := minPort; port <= maxPort; port++ {
		if !slices.Contains(exclude, port) {
			availablePorts = append(availablePorts, port)
		}
	}

	return availablePorts[rand.Intn(len(availablePorts))] //nolint: gosec It panics if len(availablePorts) == 0.
}

// SuggestPortInRange will suggest a http port in the given range.
func SuggestPortInRange(minPort, maxPort int) int {
	return SuggestPortInRangeExcluding(minPort, maxPort, []int{}) //nolint: gosec
}

// SuggestPrivatePortExcluding will suggest a http port excluding the given ports.
func SuggestPrivatePortExcluding(exclude []int) int {
	return SuggestPortInRangeExcluding(minPrivatePort, maxPrivatePort, exclude)
}

// SuggestPrivatePort will suggest a private http port.
func SuggestPrivatePort() int {
	return SuggestPrivatePortExcluding([]int{})
}
