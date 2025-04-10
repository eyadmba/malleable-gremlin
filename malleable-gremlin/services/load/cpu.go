package load

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func GenerateCPULoad(ctx context.Context, tasks int, duration time.Duration) (*LoadResult, error) {
	if tasks <= 0 {
		return nil, fmt.Errorf("number of tasks must be positive")
	}

	if duration <= 0 {
		return nil, fmt.Errorf("duration must be positive")
	}

	var wg sync.WaitGroup
	wg.Add(tasks)

	// Start CPU-intensive goroutines
	for i := 0; i < tasks; i++ {
		go func() {
			defer wg.Done()
			deadline := time.Now().Add(duration)
			for time.Now().Before(deadline) {
				select {
				case <-ctx.Done():
					return // Exit the goroutine if context is cancelled
				default:
					// CPU-intensive operation
					_ = 1 + 1
				}
			}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()

	return &LoadResult{
		TasksStarted: tasks,
		Duration:     duration,
	}, nil
}
