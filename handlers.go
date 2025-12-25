package main

import (
	"database/sql"
	"fmt"
)

func GetTodos(db *sql.DB) ([]Todo, error) {
	rows, err := db.Query("SELECT id, title FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodo(db *sql.DB, id string) (Todo, error) {
	var todo Todo

	if err := db.QueryRow("SELECT id, title FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title); err != nil {
		if err == sql.ErrNoRows {
			return Todo{}, fmt.Errorf("todo not found")
		}
		return Todo{}, err
	}

	return todo, nil
}

func AddTodo(db *sql.DB, title string) (int64, error) {
	result, err := db.Exec("INSERT INTO todos (title) VALUES( ? )", title)
	if err != nil {
		return 0, fmt.Errorf("add todo : %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add todo : %v", err)
	}

	return id, nil
}

func UpdateTodo(db *sql.DB, title string, id int64) (Todo, error) {
	result, err := db.Exec("UPDATE todos SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return Todo{}, fmt.Errorf("add todo : %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Todo{}, err
	}

	if rowsAffected == 0 {
		return Todo{}, sql.ErrNoRows
	}

	var todo Todo

	if err := db.QueryRow("SELECT id, title FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title); err != nil {
		if err == sql.ErrNoRows {
			return Todo{}, fmt.Errorf("todo not found")
		}
		return Todo{}, err
	}

	return todo, nil
}
