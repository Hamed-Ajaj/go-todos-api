# Go Todos API

A minimal REST-style Todo API written in Go using `net/http` and SQLite. It stores todos in a local `my.db` file and exposes CRUD endpoints for listing, creating, and updating items.

## How It Works

- `main.go` starts an HTTP server on `:8080` and registers routes.
- `handlers.go` performs SQLite queries for todos.
- `db.go` contains a helper to create the `todos` table (not called automatically).
- A mutex guards database access to keep requests serialized.

## Requirements

- Go 1.25+ (see `go.mod`)
- SQLite database file `my.db`

## Setup

Create the database and table:

```bash
sqlite3 my.db "CREATE TABLE IF NOT EXISTS todos (id INTEGER PRIMARY KEY, title TEXT NOT NULL);"
```

## Run Locally

```bash
go run .
```

Server listens on `http://localhost:8080`.

## API Endpoints

```text
GET    /todos
GET    /todos/{id}
POST   /todos
PUT    /todos/{id}
```

### Example Requests

Create a todo:

```bash
curl -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy milk"}'
```

List todos:

```bash
curl http://localhost:8080/todos
```

Update a todo:

```bash
curl -X PUT http://localhost:8080/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Buy oat milk"}'
```
