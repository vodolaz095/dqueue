package dqueue

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestHandler_Concurrent(t *testing.T) {
	var counter int
	ch := make(chan Task, 1000)
	wg := sync.WaitGroup{}

	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	// Publish tasks
	timestamp := time.Now()
	for i := 0; i < 100; i++ {
		h.ExecuteAfter(i, time.Second+100*time.Millisecond+10*time.Millisecond*time.Duration(i))
	}

	tasks := h.Dump()
	for i := range tasks {
		t.Logf("Task '%v' to be executed in %s from now.\n",
			tasks[i].Payload, tasks[i].ExecuteAt.Sub(time.Now()).String(),
		)
	}
	// consume with external context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func(workerNumber int) {
			ticker := time.NewTicker(time.Millisecond)
			for {
				select {
				case c := <-ticker.C:
					task, ready := h.Get()
					if ready {
						// notify task executed
						t.Logf("Worker %v executed `%v` sheduled for %s on %s (Get() delay %s).\n",
							workerNumber,
							task.Payload,
							task.ExecuteAt.Format("15:04:05.000"),
							time.Now().Format("15:04:05.000"),
							task.ExecuteAt.Sub(c).String(),
						)
						ch <- task
					}
					break
				case <-ctx.Done():
					t.Logf("Closing worker %v, there are %v tasks in queue\n", workerNumber, h.Len())
					wg.Done()
					ticker.Stop()
					return
				}
			}
		}(j)
	}
	wg.Wait()
	close(ch)

	for task := range ch {
		if task.Payload.(int) != counter {
			t.Errorf("wrong order %v", task.Payload)
		}
		counter++

		if task.ExecuteAt.Before(timestamp) {
			t.Errorf("wrong task order for %s %s",
				task.Payload,
				task.ExecuteAt.Format("15:04:05.000"),
			)
		}
		timestamp = task.ExecuteAt
	}
}
