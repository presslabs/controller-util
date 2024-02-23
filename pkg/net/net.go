package net

import (
	"math/rand"
	"slices"
)

const (
	minPrivatePort = 49152
	maxPrivatePort = 65535
)

// RandomPortInRangeExcluding will suggest a http port in the given range, excluding the given ports.
func RandomPortInRangeExcluding(startPort, stopPort int, exclude []int) int {
	availablePorts := []int{}

	for port := startPort; port <= stopPort; port++ {
		if !slices.Contains(exclude, port) {
			availablePorts = append(availablePorts, port)
		}
	}

	// Note: It panics if len(availablePorts) == 0.
	return availablePorts[rand.Intn(len(availablePorts))] //nolint: gosec
}

// RandomPortInRange will suggest a http port in the given range.
func RandomPortInRange(startPort, stopPort int) int {
	return RandomPortInRangeExcluding(startPort, stopPort, []int{})
}

// RandomPrivatePortExcluding will suggest a http port excluding the given ports.
func RandomPrivatePortExcluding(exclude []int) int {
	return RandomPortInRangeExcluding(minPrivatePort, maxPrivatePort, exclude)
}

// RandomPrivatePort will suggest a private http port.
func RandomPrivatePort() int {
	return RandomPrivatePortExcluding([]int{})
}
