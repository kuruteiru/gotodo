package db

import (
	"database/sql"
	"embed"
	"fmt"

	_ "modernc.org/sqlite"
)

const (
	dsn       = "db/database.db"
	sqlDir	  = "sql"
)

type SqlFile string

const (
	schemaSqlFile SqlFile = sqlDir+"/schema.sql"
)

var (
	//go:embed sql/*.sql
	sqlFiles embed.FS
	db       *sql.DB
)

func openDB() error {
	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed opening db [%w]", err)
	}
	return nil
}

func loadSchema() error {
	if db == nil {
		if err := openDB(); err != nil {
			return err
		}
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("database is not alive: [%w]", err)
	}

	query, err := getQuery(schemaSqlFile)
	if err != nil {
		return err
	}

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed executing %v: [%w]", schemaSqlFile, err)
	}

	return nil
}

func getQuery(name SqlFile) (string, error) {
	query, err := sqlFiles.ReadFile(string(name))
	if err != nil {
		return "", fmt.Errorf("failed reading query %s: [%w]", name, err)
	}
	return string(query), nil
}
