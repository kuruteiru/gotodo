package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kuruteiru/gotodo/models"
)

func SelectTask(id uint64) (*models.Task, error) {
	if err := checkDB(); err != nil {
		return nil, fmt.Errorf("cannot select task: [%w]", err)
	}

	query := "SELECT * FROM task WHERE id = ?"

	task := models.Task{}
	if err := db.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Created, &task.Completed, &task.TodolistID); err != nil {
		return nil, fmt.Errorf("failed selecting task with id %v: [%w]", id, err)
	}

	return &task, nil
}

func SelectTasks() ([]models.Task, error) {
	if err := checkDB(); err != nil {
		return nil, fmt.Errorf("cannot select tasks: [%w]", err)
	}

	query := "SELECT * FROM task"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed selecting tasks: [%w]", err)
	}
	defer rows.Close()

	tasks := []models.Task{}
	var scanningErrs []error
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Created, &task.Completed, &task.TodolistID); err != nil {
			scanningErrs = append(scanningErrs, fmt.Errorf("failed scanning row [%v]", err))
			continue
		}
		tasks = append(tasks, task)
	}

	if scanningErrs != nil {
		return tasks, errors.Join(scanningErrs...)
	}

	return tasks, nil
}

func InsertTask(task *models.Task) error {
	if err := checkDB(); err != nil {
		return fmt.Errorf("cannot insert task: [%w]", err)
	}

	query := `INSERT INTO 
	task (name, description, priority, created, completed, todolist_id)
	VALUES ($1, $2, $3, $4, $5, $6)`

	res, err := db.Exec(query, task.Name, task.Description, task.Priority, task.Created, task.Completed, task.TodolistID)
	if err != nil {
		return fmt.Errorf("failed inserting task %+v: [%w]", task, err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed getting affected rows: [%w]", err)
	}
	if ra != 1 {
		return fmt.Errorf("wrong number of rows affected: %v", ra)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed getting last inserted id: [%w]", err)
	}

	task.ID = uint64(id)
	return nil
}

func DeleteTask(id uint64) error {
	if err := checkDB(); err != nil {
		return fmt.Errorf("cannot delete task: [%w]", err)
	}
	
	query := "DELETE FROM task WHERE id = ?"

	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed deleting task with id %v: [%w]", id, err)
	}

	ra, err := res.RowsAffected() 
	if err != nil {
		return fmt.Errorf("failed getting affected rows: [%w]", err)
	}
	if ra != 1 {
		return fmt.Errorf("wrong number of rows affected: %v", ra)
	}

	return nil
}

//todo: fix/remove reflection, test
func UpdateTask(task models.Task) error {
	if err := checkDB(); err != nil {
		return fmt.Errorf("cannot delete task: [%w]", err)
	}
	
	buff := "UPDATE task SET "

	var sb strings.Builder
	if _, err := sb.WriteString(buff); err != nil {
		return fmt.Errorf("failed writing part of query (%v) to string builder (%v): [%w]", buff, sb.String(), err)
	}

	current, err := SelectTask(task.ID)
	if err != nil || task.ID == 0 {
		var err error
		if err = InsertTask(&task); err != nil {
			return fmt.Errorf("failed inserting a nonexistant task that should have been updated: [%w]", err)
		}
		return nil
	}

	currentValue := reflect.ValueOf(current)
	taskValue := reflect.ValueOf(task)
	if currentValue.Kind() == reflect.Ptr {
		currentValue = currentValue.Elem()
	}
	if taskValue.Kind() == reflect.Ptr {
		taskValue = taskValue.Elem()
	}
	t := reflect.TypeOf(task)

	var diff []any
	for i := range t.NumField() {
		field := t.Field(i).Name
		cf := currentValue.FieldByName(field)
		tf := taskValue.FieldByName(field)

		if reflect.DeepEqual(cf.Interface(), tf.Interface()) {
			continue
		}

		if len(diff) != 0 && i < t.NumField() - 1 {
			buff := fmt.Sprintf(", ")
			if _, err := sb.WriteString(buff); err != nil {
				return fmt.Errorf("failed writing part of query (%v) to string builder (%v): [%w]", buff, sb.String(), err)
			}
		}

		buff := fmt.Sprintf("%v = ?", strings.ToLower(field))
		if _, err := sb.WriteString(buff); err != nil {
			return fmt.Errorf("failed writing part of query (%v) to string builder (%v): [%w]", buff, sb.String(), err)
		}

		diff = append(diff, tf.Interface())
	}

	sb.WriteString(" WHERE id = ?")
	diff = append(diff, task.ID)

	fmt.Printf("diff: %v\n", diff)

	fmt.Printf("query: %v\n", sb.String())

	res, err := db.Exec(sb.String(), diff...)
	if err != nil {
		return fmt.Errorf("failed modifying task with id %v: [%w]", task.ID, err)
	}

	ra, err := res.RowsAffected() 
	if err != nil {
		return fmt.Errorf("failed getting affected rows: [%w]", err)
	}
	if ra != 1 {
		return fmt.Errorf("wrong number of rows affected: %v", ra)
	}

	return nil
}
