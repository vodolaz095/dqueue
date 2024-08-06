package dqueue

import (
	"testing"
	"time"
)

func TestPushDuplicateTime(t *testing.T) {
	h := New()
	now := time.Now().Add(200 * time.Millisecond) // same time, important!
	h.ExecuteAt("task1", now)
	h.ExecuteAt("task2", now)
	t.Logf("%v", h.Dump())
	if h.Len() != 2 {
		t.Fatal("wrong length")
	}
	dumped := h.Dump()
	t.Logf("%v", dumped)
	if dumped[0].Payload.(string) != "task1" {
		t.Errorf("element 0 has wrong payload %s", dumped[0].Payload.(string))
	}
	if dumped[1].Payload.(string) != "task2" {
		t.Errorf("element 1 has wrong payload %s", dumped[0].Payload.(string))
	}
	if dumped[0].ExecuteAt != now {
		t.Errorf("element 0 has wrong execution time")
	}
	if dumped[1].ExecuteAt != now {
		t.Errorf("element 1 has wrong execution time")
	}
	if h.Len() != 2 {
		t.Errorf("wrong length - %v", h.Len())
	}
	if h.data.nextOn != now {
		t.Errorf("wrong nextOn")
	}
	time.Sleep(200 * time.Millisecond)

	task1, ok1 := h.Get()
	if !ok1 {
		t.Fatal("no task1 extracted")
	}
	if task1.Payload.(string) != "task1" {
		t.Error("wrong task1 payload")
	}
	if task1.ExecuteAt != now {
		t.Error("wrong task1 executed at")
	}
	t.Logf("Task 1 extracted as expected")
	t.Logf("Tasks left: %v", h.Dump())
	if h.Len() != 1 {
		t.Errorf("wrong length - %v", h.Len())
	}
	if h.data.nextOn != now {
		t.Errorf("wrong nextOn")
	}

	task2, ok2 := h.Get()
	if !ok2 {
		t.Fatal("no task2 extracted")
	}
	if task2.Payload.(string) != "task2" {
		t.Error("wrong task2")
	}
	if task2.ExecuteAt != now {
		t.Error("wrong task2 executed at")
	}
	t.Logf("Task 2 extracted as expected")
	if h.Len() != 0 {
		t.Fatal("wrong length")
	}
}
