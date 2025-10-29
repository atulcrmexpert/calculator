package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
    "sync"
)

// Simple in-memory TODO app (no external dependencies).
// Endpoints:
// GET  /todos           - list all todos
// POST /todos           - create a new todo (JSON body: {"title":"...","completed":false})
// GET  /todos/{id}      - get a todo by id
// DELETE /todos/{id}    - delete a todo by id

type Todo struct {
    ID        int    `json:"id"`
    Title     string `json:"title"`
    Completed bool   `json:"completed"`
}

var (
    mu    sync.Mutex
    todos = map[int]Todo{}
    nextID = 1
)

func main() {
    http.HandleFunc("/todos", todosHandler)
    http.HandleFunc("/todos/", todoHandler) // for paths with id
    addr := ":8080"
    fmt.Println("Starting server at", addr)
    fmt.Println("Open http://localhost:8080/todos")
    log.Fatal(http.ListenAndServe(addr, nil))
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        listTodos(w, r)
    case http.MethodPost:
        createTodo(w, r)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
    // URL expected: /todos/{id}
    parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
    if len(parts) != 2 {
        http.NotFound(w, r)
        return
    }
    idStr := parts[1]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "invalid id", http.StatusBadRequest)
        return
    }

    switch r.Method {
    case http.MethodGet:
        getTodo(w, r, id)
    case http.MethodDelete:
        deleteTodo(w, r, id)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func listTodos(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()
    arr := make([]Todo, 0, len(todos))
    for _, t := range todos {
        arr = append(arr, t)
    }
    writeJSON(w, arr)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
    var t Todo
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
        return
    }
    if strings.TrimSpace(t.Title) == "" {
        http.Error(w, "title required", http.StatusBadRequest)
        return
    }
    mu.Lock()
    t.ID = nextID
    nextID++
    todos[t.ID] = t
    mu.Unlock()
    w.WriteHeader(http.StatusCreated)
    writeJSON(w, t)
}

func getTodo(w http.ResponseWriter, r *http.Request, id int) {
    mu.Lock()
    t, ok := todos[id]
    mu.Unlock()
    if !ok {
        http.NotFound(w, r)
        return
    }
    writeJSON(w, t)
}

func deleteTodo(w http.ResponseWriter, r *http.Request, id int) {
    mu.Lock()
    _, ok := todos[id]
    if ok {
        delete(todos, id)
    }
    mu.Unlock()
    if !ok {
        http.NotFound(w, r)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    enc := json.NewEncoder(w)
    enc.SetIndent("", "  ")
    if err := enc.Encode(v); err != nil {
        http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
    }
}
