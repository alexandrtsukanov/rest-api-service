package main

import (
	"log"
	"net/http"
	"os"
	"restapi-tasks/internal/database"
	"restapi-tasks/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = "postgres://taskuser:taskpass@localhost:5433/tasksdb?sslmode=prefer"
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

	mux.HandleFunc("/tasks", methodHandler(handlers.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handlers.CreateTask, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIdHandler(handlers))

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort

	globErr := http.ListenAndServe(serverAddr, loggedMux)

	if globErr != nil {
		log.Fatal(globErr)
	}
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"Method: %s, url: %s, remote address: %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
	})
}
