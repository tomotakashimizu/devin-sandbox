package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/application"
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/domain/todo"
)

// TodoHandler handles HTTP requests for todo operations
type TodoHandler struct {
	todoService *application.TodoService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService *application.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// TodoResponse represents the HTTP response for a todo
type TodoResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// CreateTodoRequest represents the HTTP request for creating a todo
type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// UpdateTodoRequest represents the HTTP request for updating a todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// todoToResponse converts a todo to a response
func todoToResponse(t *todo.Todo) TodoResponse {
	return TodoResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt.Format(http.TimeFormat),
		UpdatedAt:   t.UpdatedAt.Format(http.TimeFormat),
	}
}

// RegisterRoutes registers the todo routes
func (h *TodoHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/todos", h.CreateTodo).Methods("POST")
	r.HandleFunc("/api/todos", h.GetAllTodos).Methods("GET")
	r.HandleFunc("/api/todos/{id}", h.GetTodo).Methods("GET")
	r.HandleFunc("/api/todos/{id}", h.UpdateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{id}", h.DeleteTodo).Methods("DELETE")
	r.HandleFunc("/api/todos/{id}/complete", h.CompleteTodo).Methods("PATCH")
	r.HandleFunc("/api/todos/{id}/incomplete", h.IncompleteTodo).Methods("PATCH")
}

// CreateTodo handles the creation of a new todo
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := application.CreateTodoDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	t, err := h.todoService.Create(dto)
	if err != nil {
		if err == todo.ErrEmptyTitle {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todoToResponse(t))
}

// GetAllTodos handles retrieving all todos
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.todoService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		response = append(response, todoToResponse(t))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTodo handles retrieving a todo by ID
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := h.todoService.GetByID(id)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoToResponse(t))
}

// UpdateTodo handles updating a todo
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto := application.UpdateTodoDTO{
		Title:       req.Title,
		Description: req.Description,
	}

	t, err := h.todoService.Update(id, dto)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == todo.ErrEmptyTitle {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoToResponse(t))
}

// DeleteTodo handles deleting a todo
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.todoService.Delete(id)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CompleteTodo handles marking a todo as completed
func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := h.todoService.MarkAsCompleted(id)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoToResponse(t))
}

// IncompleteTodo handles marking a todo as incomplete
func (h *TodoHandler) IncompleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := h.todoService.MarkAsIncomplete(id)
	if err != nil {
		if err == todo.ErrTodoNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoToResponse(t))
}
