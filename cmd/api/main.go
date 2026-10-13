package main

import (
	"log"
	"net/http"
	"os"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/handlers"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5432/tasksdb?schema=public"
	}

	serverPort := os.Getenv("SERVER_PORT")

	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Server started on port %s", serverPort)

	db, err := database.Connect(databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database successfully connected")

	taskService := database.NewTaskService(db)
	handlers := handlers.NewHandlers(taskService)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", methodHandler(handlers.GetAllTasks, "GET"))
	mux.HandleFunc("/tasks/create", methodHandler(handlers.CreateTask, "POST"))
	mux.HandleFunc("/tasks/:id", taskIdHandler(handlers))
}

func methodHandler(handler http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}

		handler(w, r)
	}
}

func taskIdHandler(handler *handlers.Handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTask(w, r)
		case http.MethodPut:
			handler.UpdateTask(w, r)
		case http.MethodDelete:
			handler.DeleteTask(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
