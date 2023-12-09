package dqueue

import "time"

// Task is element describing things to be done in future
type Task struct {
	ExecuteAt time.Time
	Payload   any
}
