package db

import (
	"fmt"

	"github.com/kuruteiru/gotodo/models"
)

func selectTask(id uint64) (*models.Task, error) {
	task := models.Task{}
	query := "SELECT * FROM task WHERE ID = ?"
	if err := db.QueryRow(query, id).Scan(&task.ID, &task.Name, &task.Description, &task.Priority, &task.Created, &task.Completed); err != nil {
		return nil, fmt.Errorf("failed selecting task with id %v: [%w]", id, err)
	}

	return &task, nil
}

func insertTask(task models.Task) (uint64, error) {
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
