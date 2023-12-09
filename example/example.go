package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vodolaz095/dqueue"
)

func main() {
	wg := sync.WaitGroup{}
	handler := dqueue.New()

	// Publish tasks
	for i := 0; i < 10; i++ {
		handler.ExecuteAt(
			fmt.Sprintf("Task %v", i),
			time.Now().Add(time.Second+10*time.Millisecond*time.Duration(i)),
		)
		handler.ExecuteAfter(
			fmt.Sprintf("Task %v_bis", i),
			time.Second+100*time.Millisecond+10*time.Millisecond*time.Duration(i),
		)
	}

	tasks := handler.Dump()
	for i := range tasks {
		fmt.Printf("Task %v `%s` to be executed in %s from now.\n",
			i, tasks[i].Payload, tasks[i].ExecuteAt.Sub(time.Now()).String(),
		)
	}

	// make global context to be canceled
	mainCtx, mainCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer mainCancel()

	// consume tasks
	for j := 0; j < 10; j++ {
		wg.Add(1)
		go func(workerNumber int, initialCtx context.Context) {
			ctx, cancel := context.WithCancel(initialCtx)
			defer cancel()
			ticker := time.NewTicker(time.Millisecond)
			for {
				select {
				case t := <-ticker.C:
					task, ready := handler.Get()
					if ready {
						// notify task executed
						fmt.Printf("Worker %v executed `%s` sheduled for %s on %s (Get() delay %s).\n",
							workerNumber,
							task.Payload,
							task.ExecuteAt.Format("15:04:05.000"),
							time.Now().Format("15:04:05.000"),
							task.ExecuteAt.Sub(t).String(),
						)
					}
					break
				case <-ctx.Done():
					fmt.Printf("Closing worker %v, there are %v tasks in queue\n", workerNumber, handler.Len())
					wg.Done()
					ticker.Stop()
					return
				}
			}
		}(j, mainCtx)
	}
	wg.Wait()

	tasks = handler.Dump()
	for i := range tasks {
		fmt.Printf("Task %v left: `%v` to be executed at %s\n", i, tasks[i].Payload,
			tasks[i].ExecuteAt.Format("15:04:05"))
	}
}
