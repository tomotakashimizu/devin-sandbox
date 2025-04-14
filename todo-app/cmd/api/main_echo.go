package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/infrastructure/repository/memory"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/interfaces/api"
)

func main() {
	todoRepo := memory.NewTodoRepository()

	todoService := application.NewTodoService(todoRepo)

	todoAdapter := api.NewTodoEchoAdapter(todoService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %v", err)
	}
	swagger.Servers = nil // Clear servers to use request host

	api.RegisterHandlers(e, todoAdapter)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message":   "Todo API is running",
			"endpoints": []string{"/api/todos"},
		})
	})

	port := ":8080"
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(e.Start(port))
}
