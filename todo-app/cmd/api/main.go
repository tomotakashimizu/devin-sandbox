package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/infrastructure/repository/memory"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/interfaces/api"
)

func main() {
	// Setup repositories
	todoRepo := memory.NewTodoRepository()

	// Setup services
	todoService := application.NewTodoService(todoRepo)

	// Setup HTTP handlers
	todoHandler := api.NewTodoHandler(todoService)

	// Setup router
	router := mux.NewRouter()
	todoHandler.RegisterRoutes(router)

	// Add middleware for logging
	router.Use(loggingMiddleware)

	// Root handler for API info
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Todo API is running", "endpoints": ["/api/todos"]}`)
	})

	// Start server
	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
