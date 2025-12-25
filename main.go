package main

import (
	"fmt"
	// "github.com/gorilla/mux"
	"encoding/json"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello mother father"))
	})

	r.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "get id : %v", 1)
	})

	r.HandleFunc("POST /test", func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "User created: %+v\n", user)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
