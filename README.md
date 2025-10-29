# Simple Go TODO REST API

## What
A minimal TODO REST API written in Go using only the standard library. Stores data in memory (no database).
Endpoints:
- `GET  /todos` - list all todos
- `POST /todos` - create a todo (JSON body: {"title":"...", "completed":false})
- `GET  /todos/{id}` - get todo by id
- `DELETE /todos/{id}` - delete todo by id

## Requirements
- Go 1.16+ installed

## Run
1. Open a terminal in the project folder.
2. Run:
   ```
   go run .
   ```
3. Server listens on `:8080`.

## Examples (using curl)
Create a todo:
```
curl -s -X POST http://localhost:8080/todos -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","completed":false}'
```

List todos:
```
curl -s http://localhost:8080/todos
```

Get todo with id 1:
```
curl -s http://localhost:8080/todos/1
```

Delete todo with id 1:
```
curl -i -X DELETE http://localhost:8080/todos/1
```

## Notes
- Data is stored in-memory and lost on restart.
- No external dependencies; compile with `go build` if desired.
