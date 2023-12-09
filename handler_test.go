package dqueue

import (
	"testing"
	"time"
)

func TestHandlerEmpty(t *testing.T) {
	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	task, ready := h.Get()
	if ready {
		t.Error("task is ready")
	}
	if task.Payload != nil || !task.ExecuteAt.IsZero() {
		t.Error("task returned")
	}

	tasks := h.Dump()
	if len(tasks) != 0 {
		t.Error("tasks returned")
	}
}

func TestHandlerPublishOK(t *testing.T) {
	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	if !h.ExecuteAt("task1", time.Now().Add(time.Second)) {
		t.Error("wrong publish status")
	}
	if !h.ExecuteAfter("task2", time.Second+100*time.Millisecond) {
		t.Error("wrong publish status")
	}
	task, ready := h.Get()
	if ready {
		t.Error("task is ready")
	}
	if task.Payload != nil || !task.ExecuteAt.IsZero() {
		t.Error("task returned")
	}
	tasks := h.Dump()
	if len(tasks) != 2 {
		t.Error("wrong tasks returned")
	}
	if tasks[0].Payload.(string) != "task1" {
		t.Errorf("wrong order %s", tasks[0].Payload)
	}
	if tasks[1].Payload.(string) != "task2" {
		t.Errorf("wrong order %s", tasks[1].Payload)
	}
	time.Sleep(time.Second)
	firstTask, ok := h.Get()
	if !ok {
		t.Error("task not ready")
	}
	if firstTask.Payload.(string) != "task1" {
		t.Errorf("wrong order %s", firstTask.Payload)
	}
	if h.Len() != 1 {
		t.Error("wrong length")
	}
	time.Sleep(time.Second + 100*time.Millisecond)
	secondTask, ok := h.Get()
	if !ok {
		t.Error("task not ready")
	}
	if secondTask.Payload.(string) != "task2" {
		t.Errorf("wrong order %s", secondTask.Payload)
	}
	if firstTask.ExecuteAt.After(secondTask.ExecuteAt) {
		t.Error("wrong order")
	}
}

func TestHandlerPublishReverse(t *testing.T) {
	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	if !h.ExecuteAt("task1", time.Now().Add(time.Second+100*time.Millisecond)) {
		t.Error("wrong publish status")
	}
	if !h.ExecuteAfter("task2", time.Second) {
		t.Error("wrong publish status")
	}
	task, ready := h.Get()
	if ready {
		t.Error("task is ready")
	}
	if task.Payload != nil || !task.ExecuteAt.IsZero() {
		t.Error("task returned")
	}
	tasks := h.Dump()
	if len(tasks) != 2 {
		t.Error("wrong tasks returned")
	}
	if tasks[0].Payload.(string) != "task2" {
		t.Errorf("wrong order %s", tasks[0].Payload)
	}
	if tasks[1].Payload.(string) != "task1" {
		t.Errorf("wrong order %s", tasks[1].Payload)
	}
	time.Sleep(time.Second)
	firstTask, ok := h.Get()
	if !ok {
		t.Error("task not ready")
	}
	if firstTask.Payload.(string) != "task2" {
		t.Errorf("wrong order %s", firstTask.Payload)
	}
	if h.Len() != 1 {
		t.Error("wrong length")
	}
	time.Sleep(time.Second + 100*time.Millisecond)
	secondTask, ok := h.Get()
	if !ok {
		t.Error("task not ready")
	}
	if secondTask.Payload.(string) != "task1" {
		t.Errorf("wrong order %s", secondTask.Payload)
	}
	if firstTask.ExecuteAt.After(secondTask.ExecuteAt) {
		t.Error("wrong order")
	}
}

func TestHandlerPublishMalformed(t *testing.T) {
	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	if h.ExecuteAt("task1", time.Now().Add(-time.Second)) {
		t.Error("malformed task published")
	}
	if h.ExecuteAfter("task2", -time.Second-100*time.Millisecond) {
		t.Error("malformed task published")
	}
	task, ready := h.Get()
	if ready {
		t.Error("task is ready")
	}
	if task.Payload != nil || !task.ExecuteAt.IsZero() {
		t.Error("task returned")
	}
	tasks := h.Dump()
	if len(tasks) != 0 {
		t.Error("tasks returned")
	}
	if !h.ExecuteAfter("task3", 100*time.Millisecond) {
		t.Error("good task 3 is not published")
	}
	time.Sleep(10 * time.Millisecond)
	if !h.ExecuteAfter("task4", 10*time.Millisecond) {
		t.Error("good task 4 is not published")
	}
}

func TestHandler_Prune(t *testing.T) {
	h := New()
	if h.Len() != 0 {
		t.Error("wrong length")
	}
	h.ExecuteAt("task1", time.Now().Add(time.Second))
	h.ExecuteAfter("task2", time.Second+100*time.Millisecond)
	task, ready := h.Get()
	if ready {
		t.Error("task is ready")
	}
	if task.Payload != nil || !task.ExecuteAt.IsZero() {
		t.Error("task returned")
	}
	h.Prune()
	tasks := h.Dump()
	if len(tasks) != 0 {
		t.Error("tasks returned")
	}
}

func TestHandlerUnInitializedLen(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	h.Len()
}

func TestHandlerUnInitializedGet(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	_, _ = h.Get()
}

func TestHandlerUnInitializedExecuteAt(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	_ = h.ExecuteAt("1", time.Now().Add(time.Minute))
}

func TestHandlerUnInitializedExecuteAfter(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	_ = h.ExecuteAfter("1", time.Minute)
}

func TestHandlerUnInitializedDump(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	h.Dump()
}

func TestHandlerUnInitializedPrune(t *testing.T) {
	defer func() {
		r := recover()
		if r.(string) != "in order to use dqueue.Handler it should be created via dqueue.New()" {
			t.Error("wrong panic")
		}
	}()
	h := Handler{}
	h.Prune()
}
