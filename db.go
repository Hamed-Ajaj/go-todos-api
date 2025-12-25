package main

import "database/sql"

func CreateTodoTable(db *sql.DB) (sql.Result, error) {
	sql := `CREATE TABLE IF NOT EXISTS todos (
        id INTEGER PRIMARY KEY,
        title     TEXT NOT NULL
    );`

	return db.Exec(sql)
}
