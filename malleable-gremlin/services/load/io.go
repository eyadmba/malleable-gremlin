package load

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func GenerateIOLoad(ctx context.Context, tasks int, wait time.Duration, parallel int) (*LoadResult, error) {
	if tasks <= 0 {
		return nil, fmt.Errorf("number of tasks must be positive")
	}

	if wait <= 0 {
		return nil, fmt.Errorf("wait duration must be positive")
	}

	if parallel <= 0 {
		return nil, fmt.Errorf("parallel tasks must be positive")
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, parallel)
	wg.Add(tasks)

	startTime := time.Now()
	ctxDone := ctx.Done()

	// Start IO-intensive goroutines
	for i := 0; i < tasks; i++ {
		select {
		case semaphore <- struct{}{}:
			go func() {
				defer func() {
					<-semaphore
					wg.Done()
				}()

				select {
				case <-time.After(wait):
				case <-ctxDone:
					return
				}
			}()
		case <-ctxDone:
			wg.Done()
		}
	}

	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(waitChan)
	}()

	select {
	case <-waitChan:
		return &LoadResult{
			TasksStarted: tasks,
			Duration:     time.Since(startTime),
		}, nil
	case <-ctxDone:
		return nil, ctx.Err()
	}
}
