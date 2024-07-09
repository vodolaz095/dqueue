D(EFERRED) Queue
======================

[![Go](https://github.com/vodolaz095/dqueue/actions/workflows/go.yml/badge.svg)](https://github.com/vodolaz095/dqueue/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/vodolaz095/dqueue)](https://pkg.go.dev/github.com/vodolaz095/dqueue?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/vodolaz095/dqueue)](https://goreportcard.com/report/github.com/vodolaz095/dqueue)


It was a test task i finished in 2 hours in 2017 year, i polished code a little, created
example with contexts and added 100% unit tests coverage in 2023.

What does it do?
======================
With this package we can make deferred queue of tasks to be executed, like
`execute this in 3 minutes`, `execute that in 15 seconds from now` and so on.
Then, we can consume this tasks by concurrent goroutines and they (tasks) will be
provided to consumers in proper order, like first task will be `that` to be executed in
15 seconds from now.



Usage
======================
See full example at [example.go](example%2Fexample.go)

```

handler := dqueue.New() // import "github.com/vodolaz095/dqueue"

// Publish tasks
handler.ExecuteAt(something, time.Now().Add(time.Minute))
handler.ExecuteAfter(something, time.Minute)

// make global context to be canceled
wg := sync.WaitGroup{}
mainCtx, mainCancel := context.WithTimeout(context.Background(), 3*time.Second)
defer mainCancel()

// Start concurrent consumers
wg := sync.WaitGroup{}
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
                if ready { // task is ready
                    err := ProcessTask(task)
                    if err != nil { // aka, requeue message to be delivered in 1 minute
                      handler.ExecuteAfter(something, time.Minute)
                    }
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

// See tasks left
tasks, err := handler.Dump()


```
