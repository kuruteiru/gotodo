package db

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/kuruteiru/gotodo/models"

	_ "modernc.org/sqlite"
)

const (
	dsn           = "db/gotodo.db"
	createDBQuery = "createdb"
)

var (
	//go:embed queries/*.sql
	queries embed.FS
	db      *sql.DB
)

func Main() {
	if db == nil {
		if err := openDB(); err != nil {
			fmt.Printf("failed opening db: [%v]\n", err.Error())
			return
		}
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Printf("database is not alive: %v\n", err)
	}

	if err := createDB(); err != nil {
		fmt.Printf("failed creating db: [%v]\n", err)
	}

	var query string
	// task := models.NewTask("task1", "desc1", models.TaskPriorityLow)
	// query = "INSERT INTO tasks (name, description, priority, created, completed) VALUES ($1, $2, $3, $4, $5)"
	// if _, err := db.Exec(query, task.Name, task.Description, uint8(task.Priority), task.Created, task.Completed); err != nil {
	// 	fmt.Printf("query failed: %v\n", err)
	// }

	query = "SELECT * FROM tasks"
	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("failed selecting rows from tasks: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.Priority, &t.Created, &t.Completed); err != nil {
			fmt.Printf("failed scanning rows %v", rows.Err().Error())
			break
		}
		fmt.Printf("%+v\n", t)
	}
}

func openDB() error {
	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed opening db [%w]", err)
	}
	return nil
}

func getQuery(name string) (string, error) {
	query, err := queries.ReadFile(fmt.Sprintf("queries/%s.sql", name))
	if err != nil {
		return "", fmt.Errorf("failed reading query: [%w]", err)
	}
	return string(query), nil
}

func createDB() error {
	if db == nil {
		if err := openDB(); err != nil {
			return err
		}
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("database is not alive: %v\n", err)
	}

	query, err := getQuery(createDBQuery)
	if err != nil {
		return err
	}

	if _, err := db.Exec(query); err != nil {
		return err
	}

	return nil
}
