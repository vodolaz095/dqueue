package dqueue

import (
	"container/heap"
	"sort"
	"sync"
	"time"
)

type Handler struct {
	mu          sync.Mutex
	data        *queue
	initialized bool
}

// New creates new deferred queue Handler.
func New() Handler {
	data := queue{}
	data.nextOn = time.Now().Add(maxNextInterval)
	heap.Init(&data)
	return Handler{data: &data, initialized: true}
}

func (h *Handler) checkInitialized() {
	if !h.initialized {
		panic("in order to use dqueue.Handler it should be created via dqueue.New()")
	}
}

// Len is thread save function to see how much tasks are in queue.
func (h *Handler) Len() int {
	h.checkInitialized()
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.data.Len()
}

// Get extracts one of tasks from queue in thread safe manner. If second
// argument is true, it means task is ready to be executed and is removed from
// queue. If there are no ready tasks in queue, first argument is zero Task struct
// and second one - false.
func (h *Handler) Get() (task Task, ready bool) {
	h.checkInitialized()
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.data.Len() == 0 {
		return Task{}, false
	}
	if h.data.nextOn.After(time.Now()) {
		return Task{}, false
	}
	item := heap.Pop(h.data).(*Task)
	if item.ExecuteAt.Before(time.Now()) { // ready
		return *item, true
	}
	heap.Push(h.data, item)
	return Task{}, false
}

// ExecuteAt schedules task for execution on time desired, it returns true,
// if task is accepted.
func (h *Handler) ExecuteAt(payload any, when time.Time) (ok bool) {
	h.checkInitialized()
	if when.Before(time.Now()) {
		return false
	}
	task := Task{ExecuteAt: when, Payload: payload}
	h.mu.Lock()
	heap.Push(h.data, &task)
	h.mu.Unlock()
	return true
}

// ExecuteAfter schedules task for execution on after time.Duration provided, it
// returns true, if task is accepted.
func (h *Handler) ExecuteAfter(payload any, after time.Duration) (ok bool) {
	return h.ExecuteAt(payload, time.Now().Add(after))
}

// Dump returns copy of contents of the queue in sorted manner, leaving queue intact.
func (h *Handler) Dump() (ret []*Task) {
	h.checkInitialized()
	ret = make([]*Task, h.data.Len())
	for i := range h.data.tasks {
		ret[i] = h.data.tasks[i]
	}
	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].ExecuteAt.Before(ret[j].ExecuteAt)
	})
	return ret
}

// Prune resets queue.
func (h *Handler) Prune() {
	h.checkInitialized()
	h.mu.Lock()
	h.data.prune()
	h.mu.Unlock()
}
