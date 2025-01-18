package models

import (
    "fmt"
    "time"
    "net/http"
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
    Id          uint
    Name        string
    Description string
    Priority    TaskPriority
    Created     time.Time
    Completed   *time.Time
}

func NewTask(id uint, name string, description string, priority TaskPriority) Task {
    return Task{
        Id:          id,
        Name:        name,
        Description: description,
        Priority:    priority,
        Created:     time.Now(),    
    }
}

func (t *Task) Delete() error {
    return nil
}

func (t *Task) Update() error {
    return nil
}

func (t *Task) Complete() {
    now := time.Now()
    t.Completed = &now
}

func GetTask(id uint) Task {
    return Task{}
}

func GetTasks() []Task {
    return []Task{}
}

func GenerateTasks() []Task {
    count := 10
    tasks := []Task{}

    for i := 0; i < count; i++ {
        tasks = append(tasks, Task{
            Id: uint(i),
            Name: fmt.Sprintf("task %v", uint(i)),
            Description: fmt.Sprintf("task desc %v", uint(i)),
            Created: time.Now().UTC().Truncate(time.Second),
            Priority: TaskPriorityNone,
        })
        fmt.Printf("%+v\n", tasks[i])
        fmt.Printf("task created: %v\n", tasks[i].Created)
        fmt.Printf("task priority: %v = %v\n\n", tasks[i].Priority, tasks[i].Priority.String())
    }

    <-time.After(3 * time.Second)
    tasks[2].Complete()
    <-time.After(4 * time.Second)
    tasks[3].Complete()
    <-time.After(2 * time.Second)
    tasks[7].Complete()
    <-time.After(1 * time.Second)
    tasks[9].Complete()

    return tasks
}

func ViewTodolist(w http.ResponseWriter, r *http.Request) {
    tasks := GenerateTasks()
}
