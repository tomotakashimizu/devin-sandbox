package todo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Common errors
var (
	ErrTodoNotFound = errors.New("todo not found")
	ErrEmptyTitle   = errors.New("todo title cannot be empty")
)

// Todo represents a task that needs to be completed
type Todo struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewTodo creates a new Todo entity
func NewTodo(title, description string) (*Todo, error) {
	if title == "" {
		return nil, ErrEmptyTitle
	}

	now := time.Now()
	return &Todo{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// MarkAsCompleted marks a todo as completed
func (t *Todo) MarkAsCompleted() {
	t.Completed = true
	t.UpdatedAt = time.Now()
}

// MarkAsIncomplete marks a todo as incomplete
func (t *Todo) MarkAsIncomplete() {
	t.Completed = false
	t.UpdatedAt = time.Now()
}

// Update updates the todo's title and description
func (t *Todo) Update(title, description string) error {
	if title == "" {
		return ErrEmptyTitle
	}

	t.Title = title
	t.Description = description
	t.UpdatedAt = time.Now()
	return nil
}
