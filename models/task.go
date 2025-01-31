package models

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

var (
	currentTaskID uint64 = 0
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
	ID          uint64
	Name        string
	Description string
	Priority    TaskPriority
	Created     time.Time
	Completed   *time.Time
}

func NewTask(name string, description string, priority TaskPriority) Task {
	t := Task{
		ID:          currentTaskID,
		Name:        name,
		Description: description,
		Priority:    priority,
		Created:     time.Now().UTC(),
	}
	currentTaskID++
	return t
}

func (t *Task) Complete() {
	now := time.Now().UTC()
	t.Completed = &now
}

func GetTask(id uint64) *Task {
	t := &GenerateTasks(1)[0]
	t.ID = id
	return t
}

func GenerateTasks(count int) []Task {
	tasks := []Task{}
	for i := range count {
		tasks = append(tasks, Task{
			ID:          uint64(i),
			Name:        fmt.Sprintf("task %v", uint(i)),
			Description: fmt.Sprintf("task desc %v", uint(i)),
			Created:     time.Now().UTC(),
			Priority:    TaskPriorityNone,
		})
	}
	return tasks
}

func PrintTasks(tasks []Task) {
	for i, task := range tasks {
		fmt.Printf("%v: %+v\n", i, task)
	}
}

func PrintTasksTable(tasks []Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', tabwriter.TabIndent)
	defer w.Flush()

	fmt.Fprintf(w, "id\tname\tdescription\tpriority\tcreated\tcompleted\t\n")
	for _, t := range tasks {
		completed := "false"
		if t.Completed != nil {
			completed = t.Completed.Format(time.DateTime)
		}

		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t\n", t.ID, t.Name, t.Description,
			t.Priority.String(), t.Created.Format(time.DateTime), completed)
	}
}
