package load

import (
	"fmt"
	"runtime"
	"time"
)

func GenerateMemoryLoad(size int64, gcAfter time.Duration) (*LoadResult, error) {
	if size <= 0 {
		return nil, fmt.Errorf("memory size must be positive")
	}

	// Allocate memory
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	// Wait for specified duration before GC
	if gcAfter > 0 {
		time.Sleep(gcAfter)
	} else if gcAfter == 0 {
		// Immediate GC
		runtime.GC()
	}

	return &LoadResult{
		TasksStarted: 1,
		Duration:     gcAfter,
	}, nil
}
