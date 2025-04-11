// Package tool provides utility functions such as parallelism calculations,
// work simulation, and panic recovery.
package tool

import (
	"runtime"
	"time"
)

// GetMaxParallelism returns a multiplier-based value derived from GOMAXPROCS.
func GetMaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	return maxProcs * 8
}

// SimulateWork pauses execution for the specified number of seconds to simulate work.
func SimulateWork(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
