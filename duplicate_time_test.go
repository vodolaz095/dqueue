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
	t.Logf("%v", h.Dump())
	time.Sleep(200 * time.Millisecond)

	task1, ok1 := h.Get()
	if !ok1 {
		t.Fatal("no task1 extracted")
	}
	if task1.Payload.(string) != "task1" {
		t.Error("wrong task1")
	}
	t.Logf("Task 1 extracted as expected")
	t.Logf("%v", h.Dump())

	task2, ok2 := h.Get()
	if !ok2 {
		t.Fatal("no task2 extracted")
	}
	if task2.Payload.(string) != "task2" {
		t.Error("wrong task2")
	}
	t.Logf("Task 2 extracted as expected")
	if h.Len() != 0 {
		t.Fatal("wrong length")
	}
}
