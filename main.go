package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/glebarez/go-sqlite"
)

type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	db, err := sql.Open("sqlite", "./my.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	fmt.Println("sqlite db connected")
	r := http.NewServeMux()

	// A lock makes sure only ONE thing at a time can access some shared data.
	var mu sync.Mutex

	r.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		todos, err := GetTodos(db)
		if err != nil {
			http.Error(w, "failed to fetch todo", http.StatusInternalServerError)
		}

		json.NewEncoder(w).Encode(todos)
	})

	r.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		// With a lock:
		// one person writes
		// finishes
		// next person writes
		mu.Lock()
		defer mu.Unlock()

		todo, err := GetTodo(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(todo)
	})

	r.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		var todo Todo

		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()
		_, err = AddTodo(db, todo.Title)
		if err != nil {
			http.Error(w, "failed to create todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("DELETE /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		err = DeleteTodo(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to delete todo", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	r.HandleFunc("PUT /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}

		var input struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if input.Title == "" {
			http.Error(w, "Title is required", http.StatusBadRequest)
			return
		}

		todo, err := UpdateTodo(db, input.Title, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "todo not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to update todo", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todo)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
