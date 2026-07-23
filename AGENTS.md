# dqueue Agent Guide

This document provides essential context for agents working in the dqueue codebase - a Go package for managing deferred task queues.

## Project Overview

Dqueue is a Go package that implements a deferred queue for executing tasks at specified future times. Tasks are stored with a payload (any type) and execution time. The queue ensures tasks are delivered in chronological order when they become ready.

Key characteristics:
- Pure Go implementation with no external dependencies
- Uses Go's `container/heap` for efficient priority queue operations
- Thread-safe via mutex protection
- Tasks can have payloads of any type (`any` interface{})
- Maximum task delay capped at 24 hours by design

## Essential Commands

```bash
# Run the example application
make start

# Run tests with verbose output
make test

# Run tests with coverage reporting
make cover

# Run linting (formatting, golint, vet)
make lint

# Run race detector on example
make race
```

## Code Organization

The codebase is minimal and focused:

- `handler.go`: Main interface with thread-safe methods for queue operations
- `queue.go`: Internal heap implementation that satisfies `heap.Interface`
- `task.go`: Task structure definition
- `constants.go`: Package constants (currently only `maxNextInterval`)
- `doc.go`: Package documentation

## Architecture and Data Flow

The system uses a two-layer architecture:

1. **Handler layer** (`Handler` struct): Provides thread-safe public API with mutex protection
2. **Queue layer** (`queue` struct): Internal heap implementation that manages task ordering

Data flow:
1. External code calls `ExecuteAt()` or `ExecuteAfter()` on a `Handler`
2. Handler acquires mutex and pushes task to internal `queue`
3. When `Get()` is called, Handler acquires mutex and checks if the next task is ready
4. If ready, task is popped from heap and returned; otherwise, zero value is returned

The `nextOn` field in `queue` is a performance optimization that tracks when the next task will be ready, avoiding unnecessary heap operations.

## Key Patterns and Conventions

### Type System
- Task payloads use `any` (interface{}) type, allowing any Go value
- This introduces type assertion overhead when consuming tasks
- Consider the generic alternative (dgqueue) for performance-critical applications

### Error Handling
- Methods return boolean `ok` values rather than errors
- Invalid operations (past execution times) fail silently by returning `false`
- Uninitialized handlers panic when used (defensive programming)

### Initialization
- Structs must be created via `New()` constructor
- `initialized` flag prevents misuse of zero-value structs
- Direct struct initialization will cause panics

## Testing Approach

- 100% unit test coverage (as stated in README)
- Tests focus on:
  - Basic queue operations (add, get, length)
  - Time-based execution ordering
  - Duplicate task handling
  - Concurrent access safety
- Example in `example.go` serves as integration test

## Gotchas and Non-Obvious Behavior

### 24-Hour Limitation
- Tasks cannot be scheduled more than 24 hours in the future
- This is enforced by `maxNextInterval = 24 * time.Hour`
- Rationale: Long-running applications may experience memory leaks or GC issues

### Silent Failure on Invalid Times
- `ExecuteAt()` returns `false` (not an error) when given a past time
- The task is simply not added to the queue
- No indication beyond the boolean return value

### Heap Re-Push on Unready Tasks
- When `Get()` is called but the next task isn't ready, it's temporarily popped and re-pushed
- This maintains heap invariants but adds minor overhead
- See `Get()` method implementation

### Mutex Scope
- Mutex is held only for the minimal critical section
- Heap operations (`heap.Pop`, `heap.Push`) occur outside mutex lock
- This reduces contention in concurrent scenarios

## Performance Considerations

- Type assertions on `any` payloads have runtime cost
- For high-throughput scenarios, consider the generic `dgqueue` alternative
- The 24-hour limit prevents unbounded queue growth
- `nextOn` optimization reduces unnecessary heap operations when no tasks are ready

## Future Development

The README mentions a generic-based alternative at https://github.com/vodolaz095/dgqueue for cases where "less agile and more performant" behavior is needed, suggesting this package prioritizes flexibility over raw performance.