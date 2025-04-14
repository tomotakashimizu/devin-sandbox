package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/domain/todo"
)

type TodoEchoAdapter struct {
	todoService *application.TodoService
}

func NewTodoEchoAdapter(todoService *application.TodoService) *TodoEchoAdapter {
	return &TodoEchoAdapter{
		todoService: todoService,
	}
}

func (a *TodoEchoAdapter) GetAllTodos(ctx echo.Context) error {
	todos, err := a.todoService.GetAll()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	response := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		response = append(response, domainToResponse(t))
	}

	return ctx.JSON(http.StatusOK, response)
}

func (a *TodoEchoAdapter) GetTodoById(ctx echo.Context, id types.UUID) error {
	t, err := a.todoService.GetByID(string(id))
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, domainToResponse(t))
}

func (a *TodoEchoAdapter) CreateTodo(ctx echo.Context) error {
	var req CreateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
	}

	dto := application.CreateTodoDTO{
		Title:       req.Title,
		Description: stringPtrToString(req.Description),
	}

	t, err := a.todoService.Create(dto)
	if err != nil {
		if errors.Is(err, todo.ErrEmptyTitle) {
			return ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, domainToResponse(t))
}

func (a *TodoEchoAdapter) UpdateTodo(ctx echo.Context, id types.UUID) error {
	var req UpdateTodoRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
	}

	dto := application.UpdateTodoDTO{
		Title:       req.Title,
		Description: stringPtrToString(req.Description),
	}

	t, err := a.todoService.Update(string(id), dto)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
		}
		if errors.Is(err, todo.ErrEmptyTitle) {
			return ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, domainToResponse(t))
}

func (a *TodoEchoAdapter) DeleteTodo(ctx echo.Context, id types.UUID) error {
	err := a.todoService.Delete(string(id))
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (a *TodoEchoAdapter) CompleteTodo(ctx echo.Context, id types.UUID) error {
	t, err := a.todoService.MarkAsCompleted(string(id))
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, domainToResponse(t))
}

func (a *TodoEchoAdapter) IncompleteTodo(ctx echo.Context, id types.UUID) error {
	t, err := a.todoService.MarkAsIncomplete(string(id))
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			return ctx.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, domainToResponse(t))
}


func domainToResponse(t *todo.Todo) TodoResponse {
	var description *string
	if t.Description != "" {
		description = &t.Description
	}

	return TodoResponse{
		Id:          types.UUID(t.ID),
		Title:       t.Title,
		Description: description,
		Completed:   t.Completed,
		CreatedAt:   parseTime(t.CreatedAt),
		UpdatedAt:   parseTime(t.UpdatedAt),
	}
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC1123, timeStr)
	return t
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var _ ServerInterface = (*TodoEchoAdapter)(nil)
