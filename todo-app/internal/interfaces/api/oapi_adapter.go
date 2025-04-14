package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/domain/todo"
)

type TodoAPIAdapter struct {
	todoService *application.TodoService
}

func NewTodoAPIAdapter(todoService *application.TodoService) *TodoAPIAdapter {
	return &TodoAPIAdapter{
		todoService: todoService,
	}
}

func (a *TodoAPIAdapter) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := a.todoService.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		response = append(response, todoToResponse(t))
	}

	writeJSON(w, http.StatusOK, response)
}

func (a *TodoAPIAdapter) GetTodoById(w http.ResponseWriter, r *http.Request, id string) {
	t, err := a.todoService.GetByID(id)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, todoToResponse(t))
}

func (a *TodoAPIAdapter) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dto := application.CreateTodoDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	t, err := a.todoService.Create(dto)
	if err != nil {
		if errors.Is(err, todo.ErrEmptyTitle) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, todoToResponse(t))
}

func (a *TodoAPIAdapter) UpdateTodo(w http.ResponseWriter, r *http.Request, id string) {
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dto := application.UpdateTodoDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	t, err := a.todoService.Update(id, dto)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, todo.ErrEmptyTitle) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, todoToResponse(t))
}

func (a *TodoAPIAdapter) DeleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	err := a.todoService.Delete(id)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a *TodoAPIAdapter) CompleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	t, err := a.todoService.MarkAsCompleted(id)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, todoToResponse(t))
}

func (a *TodoAPIAdapter) IncompleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	t, err := a.todoService.MarkAsIncomplete(id)
	if err != nil {
		if errors.Is(err, todo.ErrTodoNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, todoToResponse(t))
}


func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

var _ ServerInterface = (*TodoAPIAdapter)(nil)

type ServerInterface interface {
	GetAllTodos(w http.ResponseWriter, r *http.Request)
	GetTodoById(w http.ResponseWriter, r *http.Request, id string)
	CreateTodo(w http.ResponseWriter, r *http.Request)
	UpdateTodo(w http.ResponseWriter, r *http.Request, id string)
	DeleteTodo(w http.ResponseWriter, r *http.Request, id string)
	CompleteTodo(w http.ResponseWriter, r *http.Request, id string)
	IncompleteTodo(w http.ResponseWriter, r *http.Request, id string)
}
