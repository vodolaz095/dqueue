package dqueue

import (
	"testing"
	"time"
)

var testQueue queue

func TestQueue_Push(t *testing.T) {
	testQueue.nextOn = time.Now().Add(5 * time.Second)
	next := time.Now().Add(time.Second)
	testQueue.Push(&Task{
		Payload:   "Task 1",
		ExecuteAt: next,
	})
	next = next.Add(time.Second)
	testQueue.Push(&Task{
		Payload:   "Task 2",
		ExecuteAt: next,
	})
	testQueue.Push(&Task{
		Payload:   "Task 3",
		ExecuteAt: next.Add(-500 * time.Millisecond),
	})
	testQueue.Push("something to be ignored")
	testQueue.Push(12345)
	if testQueue.tasks[0].Payload.(string) != "Task 1" {
		t.Error("wrong payload for 1st task")
	}
	if testQueue.tasks[1].Payload.(string) != "Task 2" {
		t.Error("wrong payload for 2nd task")
	}
	if testQueue.tasks[2].Payload.(string) != "Task 3" {
		t.Error("wrong payload for 3rd task")
	}
	if testQueue.Len() != 3 {
		t.Error("wrong task accepted")
	}
	if testQueue.nextOn.After(next) {
		t.Error("wrong nextOn")
	}
}

func TestQueue_PushWrong(t *testing.T) {
	testQueue.Push(&Task{
		Payload:   "Task 4",
		ExecuteAt: time.Now().Add(-time.Hour),
	})
	if testQueue.Len() != 3 {
		t.Error("wrong task accepted")
	}
}

func TestQueue_Len(t *testing.T) {
	if testQueue.Len() != 3 {
		t.Error("wrong length")
	}
}

func TestQueue_Less(t *testing.T) {
	if !testQueue.Less(0, 1) {
		t.Error("wrong less behaviour")
	}
	if testQueue.Less(1, 0) {
		t.Error("wrong less behaviour")
	}
}

func TestQueue_Swap(t *testing.T) {
	testQueue.Swap(0, 1)
	if testQueue.tasks[0].Payload.(string) != "Task 2" {
		t.Error("wrong payload for 1st task")
	}
	if testQueue.tasks[1].Payload.(string) != "Task 1" {
		t.Error("wrong payload for 2nd task")
	}
	if testQueue.tasks[2].Payload.(string) != "Task 3" {
		t.Error("wrong payload for 3rd task")
	}
}

func TestQueue_Pop(t *testing.T) {
	raw := testQueue.Pop()
	task := raw.(*Task)
	if task.Payload != "Task 3" {
		t.Error("wrong payload for 3rd task")
	}
	if testQueue.Len() != 2 {
		t.Error("wrong length")
	}
}

func TestQueue_prune(t *testing.T) {
	testQueue.prune()
	if testQueue.Len() != 0 {
		t.Error("wrong length")
	}
}
