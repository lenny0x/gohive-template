package task

import (
	"context"
)

type Task interface {
	Name() string
	Run(ctx context.Context) error
}

type TaskRunner struct {
	tasks []Task
}

func NewTaskRunner() *TaskRunner {
	return &TaskRunner{
		tasks: make([]Task, 0),
	}
}

func (r *TaskRunner) Register(task Task) {
	r.tasks = append(r.tasks, task)
}

func (r *TaskRunner) Tasks() []Task {
	return r.tasks
}
