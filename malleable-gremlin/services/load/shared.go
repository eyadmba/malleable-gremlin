package load

import "time"

type LoadResult struct {
	TasksStarted int           `json:"tasks_started"`
	Duration     time.Duration `json:"duration"`
	Error        string        `json:"error,omitempty"`
}
