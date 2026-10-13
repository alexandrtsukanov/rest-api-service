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
}
