package db

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/kuruteiru/gotodo/models"
)

var mutexes = sync.Map{}

func getMutex(id uint64) *sync.RWMutex {
	mutex, _ := mutexes.LoadOrStore(id, &sync.RWMutex{})
	return mutex.(*sync.RWMutex)
}

func GetTask(id uint64) {}
func UpdateTask(id uint64) {}
func SaveTask(id uint64) {}
func DeleteTask(id uint64) {}

func ReadTasks(id uint64) ([]models.Task, error) {
	path := fmt.Sprintf("db/tasks/%v.csv", id)

	mutex := getMutex(id)
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

	if len(records) == 0 || len(records[0]) != reflect.TypeOf(models.Task{}).NumField() {
		return nil, fmt.Errorf("invalid header count")
	}

	if len(records) == 1 {
		return nil, fmt.Errorf("no tasks in %s", path)
	}

	var tasks []models.Task
	var errs []error

	for i, record := range records {
		if i == 0 {
			continue
		}

		task, err := parseTask(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed reading record %v: [%w]", i, err))
			continue
		}

		if task == nil {
			continue
		}

		tasks = append(tasks, *task)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no valid tasks found")
	}

	if errs != nil {
		return tasks, fmt.Errorf("invalid records: [%w]", errors.Join(errs...))
	}

	return tasks, nil
}

func WriteTasks(id uint64, tasks []models.Task) error {
	path := fmt.Sprintf("db/tasks/%v.csv", id)
	tmp := path + ".tmp"

	mutex := getMutex(id)
	mutex.Lock()
	defer mutex.Unlock()

	old, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)

	}
	defer old.Close()

	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer func() {
		f.Close()
		os.Remove(tmp)
	}()

	r := csv.NewReader(old)
	w := csv.NewWriter(f)

	oldTasks, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("failed reading csv file %s: %w", path, err)
	}

	for i, t := range oldTasks {
		if err := w.Write(t); err != nil {
			return fmt.Errorf("failed copying old record %v: [%w]", i, err)
		}
	}

	for i, t := range tasks {
		if err := w.Write(formatTask(t)); err != nil {
			return fmt.Errorf("failed writing record %v: [%w]", i, err)
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("failed flushing: [%w]", err)
	}

	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("failed swapping tmp file %v: [%w]", tmp, err)
	}

	return nil
}

func parseTask(record []string) (*models.Task, error) {
	if len(record) != reflect.TypeOf(models.Task{}).NumField() {
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
