package db

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kuruteiru/gotodo/models"
)

func SelectTask(id uint64) (*models.Task, error) {
	query := "SELECT * FROM task WHERE id = ?"

	task := models.Task{}
	if err := db.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Created, &task.Completed); err != nil {
		return nil, fmt.Errorf("failed selecting task with id %v: [%w]", id, err)
	}

	return &task, nil
}

func InsertTask(task models.Task) (uint64, error) {
	query := `INSERT INTO 
	task (name, description, priority, created, completed)
	VALUES ($1, $2, $3, $4, $5)`

	res, err := db.Exec(query, task.Name, task.Description, task.Priority, task.Created, task.Completed)
	if err != nil {
		return 0, fmt.Errorf("failed inserting task %+v: [%w]", task, err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed getting affected rows: [%w]", err)
	}
	if ra != 1 {
		return 0, fmt.Errorf("wrong number of rows affected: %v", ra)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed getting last inserted id: [%w]", err)
	}

	return uint64(id), nil
}

//todo: test
func DeleteTask(id uint64) error {
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

//todo: test
func UpdateTask(task models.Task) error {
	buff := "UPDATE task SET "

	var sb strings.Builder
	if _, err := sb.WriteString(buff); err != nil {
		return fmt.Errorf("failed writing part of query (%v) to string builder (%v): [%w]", buff, sb.String(), err)
	}

	current, err := SelectTask(task.ID)
	if err != nil || task.ID == 0 {
		var err error
		if task.ID, err = InsertTask(task); err != nil {
			return fmt.Errorf("failed inserting a nonexistant task that should have been updated: [%w]", err)
		}
		return nil
	}

	currentValue := reflect.ValueOf(current)
	taskValue := reflect.ValueOf(task)
	t := reflect.TypeOf(current)

	var diff []any
	for i := range t.NumField() {
		field := t.Field(i).Name
		cf := currentValue.FieldByName(field)
		tf := taskValue.FieldByName(field)

		if cf == tf {
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

		diff = append(diff, tf)
	}

	sb.WriteString(" WHERE id = ?")

	res, err := db.Exec(sb.String(), diff, task.ID)
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
