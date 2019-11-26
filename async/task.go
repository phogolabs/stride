package async

import (
	"context"
	"sync"
)

// TaskFunc is the task's func
type TaskFunc func(ctx context.Context) error

// Task represents a task
type Task struct {
	ctx    context.Context
	cancel func()
	err    error
	wg     *sync.WaitGroup
	exec   TaskFunc
	data   []interface{}
}

// NewTask creates a new task
func NewTask(fn TaskFunc, data ...interface{}) *Task {
	ctx, cancel := context.WithCancel(context.Background())
	return &Task{
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
		exec:   fn,
		data:   data,
	}
}

// Data returns the tasks data
func (t *Task) Data() interface{} {
	switch len(t.data) {
	case 0:
		return nil
	case 1:
		return t.data[0]
	default:
		return t.data
	}
}

// Run runs the task
func (t *Task) Run() {
	t.wg.Add(1)

	go func() {
		if err := t.exec(t.ctx); err != nil {
			t.err = err
		}

		t.wg.Done()
	}()
}

// Stop stops the task
func (t *Task) Stop() error {
	t.cancel()
	t.wg.Wait()
	return t.err
}

// Wait waits the task to be stopped
func (t *Task) Wait() error {
	t.wg.Wait()
	return t.err
}
