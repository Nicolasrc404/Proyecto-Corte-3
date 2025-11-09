package server

import (
	"fmt"
	"time"
)

type Task struct {
	ID        int
	Name      string
	StartedAt time.Time
	Done      bool
}

type TaskQueue struct {
	tasks chan Task
}

func NewTaskQueue() *TaskQueue {
	return &TaskQueue{
		tasks: make(chan Task, 100),
	}
}

func (q *TaskQueue) Start() {
	go func() {
		for task := range q.tasks {
			fmt.Printf("⏳ Ejecutando tarea: %s (ID %d)\n", task.Name, task.ID)
			time.Sleep(2 * time.Second)
			fmt.Printf("✅ Tarea completada: %s\n", task.Name)
		}
	}()
}

func (q *TaskQueue) AddTask(t Task) {
	q.tasks <- t
}
