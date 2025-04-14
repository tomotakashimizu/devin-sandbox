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
	todoRepo := memory.NewTodoRepository()

	todoService := application.NewTodoService(todoRepo)

	todoAdapter := api.NewTodoAPIAdapter(todoService)

	router := mux.NewRouter()

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}
	swagger.Servers = nil // Clear servers to use request host

	router.Use(loggingMiddleware)
	router.Use(api.OapiRequestValidator(swagger))

	api.HandlerFromMux(todoAdapter, router)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Todo API is running", "endpoints": ["/api/todos"]}`)
	})

	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./api/swagger-ui"))))

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
