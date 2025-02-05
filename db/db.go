package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func Main() {
	db, err := sql.Open("sqlite", "/db/gotodo.db")
	if err != nil {
		fmt.Printf("failed opening db")
		return
	}
	defer db.Close()
}
