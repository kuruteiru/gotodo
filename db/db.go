package db

import (
	"encoding/csv"
	"fmt"
	// "math/rand/v2"
	"os"
	"strconv"
	"time"

	"github.com/kuruteiru/gotodo/models"
)

func Main() {
	t, err := ReadTasks(1)
	if err != nil {
		fmt.Printf("errr1: %v\n", err.Error())
	}
	models.PrintTasksTable(t)
	fmt.Println()

	// if err := WriteTasks(1, models.GenerateTasks(5)); err != nil {
	// 	fmt.Printf("errw: %v\n", err.Error())
	// }
	//
	// t, err = ReadTasks(1)
	// if err != nil {
	// 	fmt.Printf("errr2: %v\n", err.Error())
	// }
	// models.PrintTasksTable(t)
}

// todo: handle errors
func ReadTasks(id uint64) ([]models.Task, error) {
	path := fmt.Sprintf("db/tasks/%v.csv", id)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	for i, record := range records {
		if i == 0 {
			continue
		}

		id, _ := strconv.ParseUint(record[0], 10, 64)
		name := record[1]
		description := record[2]
		priority, _ := strconv.ParseUint(record[3], 10, 8)
		created, _ := time.Parse(time.DateTime, record[4])
		var completed *time.Time
		if record[5] != "" {
			c, _ := time.Parse(time.DateTime, record[5])
			completed = &c
		}

		tasks = append(tasks, models.Task{
			ID:          id,
			Name:        name,
			Description: description,
			Priority:    models.TaskPriority(priority),
			Created:     created,
			Completed:   completed,
		})
	}

	return tasks, nil
}

// todo: handle errors
func WriteTasks(id uint64, tasks []models.Task) error {
	path := fmt.Sprintf("db/tasks/%v.csv", id)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("err: %v\n", err.Error())
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	for _, t := range tasks {
		record := []string{
			strconv.FormatUint(t.ID, 10),
			t.Name,
			t.Description,
			strconv.FormatUint(uint64(t.Priority), 10),
			t.Created.Format(time.DateTime),
			// t.Created.Add(time.Duration(rand.IntN(100)+1) * time.Second).Format(time.DateTime),
		}

		if t.Completed != nil {
			record = append(record, t.Completed.Format(time.DateTime))
		} else {
			record = append(record, "")
		}

		if err := w.Write(record); err != nil {
			return err
		}
	}

	return nil
}
