package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kuruteiru/gotodo/models"
)

var mutex sync.RWMutex

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

func ReadTasks(id uint64) ([]models.Task, error) {
	path := fmt.Sprintf("db/tasks/%v.csv", id)

	mutex.RLock()
	defer mutex.RUnlock()

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed reading csv file %s: %w", path, err)
	}

	var tasks []models.Task
	var errs []error

	for i, record := range records {
		if i == 0 {
			continue
		}

		task, err := parseTask(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("can't read record %v: [%w]", i, err))
			continue
		}

		if task == nil {
			continue
		}

		tasks = append(tasks, *task)
	}

	if len(tasks) == 0 {
		return nil, errors.New("no valid tasks found")
	}
	
	if errs != nil {
		return tasks, fmt.Errorf("invalid records: [%w]", errors.Join(errs...))
	}

	return tasks, nil
}

//todo: write into tmp file instead
func WriteTasks(id uint64, tasks []models.Task) error {
	path := fmt.Sprintf("db/tasks/%v.csv", id)

	mutex.Lock()
	defer mutex.Unlock()

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)

	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	var errs []error
	for i, t := range tasks {
		if err := w.Write(formatTask(t)); err != nil {
			errs = append(errs, fmt.Errorf("couldn't write record %v: [%w]", i, err))
			continue
		}
	}

	return nil
}

func parseTask(record []string) (*models.Task, error) {
	if len(record) > 6 {
		return nil, fmt.Errorf("invalid record: [%v]", record)
	}

	id, err := strconv.ParseUint(record[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid id: [%w]", err)
	}

	name := record[1]
	description := record[2]

	priority, err := strconv.ParseUint(record[3], 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid priority: [%w]", err)
	}

	created, err := time.Parse(time.DateTime, record[4])
	if err != nil {
		return nil, fmt.Errorf("invalid creation date: [%w]", err)
	}

	var completed *time.Time
	if record[5] != "" {
		if c, err := time.Parse(time.DateTime, record[5]); err != nil {
			return nil, fmt.Errorf("invalid completion date: [%w]", err)
		} else {
			completed = &c
		}
	}

	return &models.Task{
		ID:          id,
		Name:        name,
		Description: description,
		Priority:    models.TaskPriority(priority),
		Created:     created,
		Completed:   completed,
	}, nil
}

func formatTask(task models.Task) []string {
	record := []string{
		strconv.FormatUint(task.ID, 10),
		task.Name,
		task.Description,
		strconv.FormatUint(uint64(task.Priority), 10),
		task.Created.Format(time.DateTime),
	}

	if task.Completed != nil {
		record = append(record, task.Completed.Format(time.DateTime))
	} else {
		record = append(record, "")
	}

	return record
}
