package models

import (
	"fmt"
	"time"
)

type TaskPriority uint8

const (
    TaskPriorityNone TaskPriority = iota
    TaskPriorityLow
    TaskPriorityMedium
    TaskPriorityHigh
)

var taskPriorityString = map[TaskPriority]string{
    TaskPriorityNone:   "none",
    TaskPriorityLow:    "low",
    TaskPriorityMedium: "medium",
    TaskPriorityHigh:   "high",
}

func (tp *TaskPriority) String() string {
    return taskPriorityString[*tp]
}

type Task struct {
    ID          uint
    Name        string
    Description string
    Priority    TaskPriority
    Created     time.Time
    Completed   *time.Time
}

func NewTask(id uint, name string, description string, priority TaskPriority) Task {
    return Task{
        ID:          id,
        Name:        name,
        Description: description,
        Priority:    priority,
        Created:     time.Now(),    
    }
}

func (t *Task) Complete() {
    now := time.Now().UTC()
    t.Completed = &now
}

func GenerateTasks() []Task {
    count := 10
    tasks := []Task{}

    for i := range count {
        tasks = append(tasks, Task{
            ID: uint(i),
            Name: fmt.Sprintf("task %v", uint(i)),
            Description: fmt.Sprintf("task desc %v", uint(i)),
            Created: time.Now().UTC(),
            Priority: TaskPriorityNone,
        })
    }

    // <-time.After(3 * time.Second)
    // tasks[2].Complete()
    // <-time.After(1 * time.Second)
    // tasks[3].Complete()
    // <-time.After(2 * time.Second)
    // tasks[7].Complete()
    // <-time.After(1 * time.Second)
    // tasks[9].Complete()

    return tasks
}

func PrintTasks(tasks []Task) {
    for i, task := range tasks {
        fmt.Printf("%v: %+v\n", i, task)
    }
}
