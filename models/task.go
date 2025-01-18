package models

import (
    "fmt"
    "time"
)

type Task struct {
    id          uint
    name        string
    description string
    done        bool
    created     time.Time
    priority    uint8
}

func NewTask(id uint, name string, description string) Task {
    return Task{
        id: id,
        name: name,
        description: description,
    }
}

func (t *Task) Delete() error {
    return nil
}

func (t *Task) Update() error {
    return nil
}

func GetTask(id uint) Task {
    return Task{}
}

func GetTasks() []Task {
    return []Task{}
}

func GenerateTasks() {
    count := 10
    tasks := []Task{}

    for i := 0; i < count; i++ {
        tasks = append(tasks, Task{
            id: uint(i),
            name: fmt.Sprintf("task %v", uint(i)),
            description: fmt.Sprintf("task desc %v", uint(i)),
        })
        fmt.Printf("%+v\n", tasks[i])
    }
}
