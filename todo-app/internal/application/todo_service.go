package application

import (
	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/domain/todo"
)

// TodoService defines the application service for todo operations
type TodoService struct {
	repo todo.Repository
}

// NewTodoService creates a new todo service
func NewTodoService(repo todo.Repository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// CreateTodoDTO represents data for creating a todo
type CreateTodoDTO struct {
	Title       string
	Description string
}

// UpdateTodoDTO represents data for updating a todo
type UpdateTodoDTO struct {
	Title       string
	Description string
}

// Create creates a new todo
func (s *TodoService) Create(dto CreateTodoDTO) (*todo.Todo, error) {
	t, err := todo.NewTodo(dto.Title, dto.Description)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(t); err != nil {
		return nil, err
	}

	return t, nil
}

// GetByID retrieves a todo by its ID
func (s *TodoService) GetByID(id string) (*todo.Todo, error) {
	return s.repo.GetByID(id)
}

// GetAll retrieves all todos
func (s *TodoService) GetAll() ([]*todo.Todo, error) {
	return s.repo.GetAll()
}

// Update updates a todo
func (s *TodoService) Update(id string, dto UpdateTodoDTO) (*todo.Todo, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := t.Update(dto.Title, dto.Description); err != nil {
		return nil, err
	}

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}

	return t, nil
}

// Delete deletes a todo
func (s *TodoService) Delete(id string) error {
	return s.repo.Delete(id)
}

// MarkAsCompleted marks a todo as completed
func (s *TodoService) MarkAsCompleted(id string) (*todo.Todo, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	t.MarkAsCompleted()

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}

	return t, nil
}

// MarkAsIncomplete marks a todo as incomplete
func (s *TodoService) MarkAsIncomplete(id string) (*todo.Todo, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	t.MarkAsIncomplete()

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}

	return t, nil
}
