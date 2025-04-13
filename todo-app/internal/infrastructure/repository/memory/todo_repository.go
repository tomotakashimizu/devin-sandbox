package memory

import (
	"sync"

	"github.com/tomotakashimizu/devin-sandbox/todo-app/internal/domain/todo"
)

// TodoRepository implements the todo.Repository interface with in-memory storage
type TodoRepository struct {
	todos map[string]*todo.Todo
	mutex *sync.RWMutex
}

// NewTodoRepository creates a new in-memory todo repository
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		todos: make(map[string]*todo.Todo),
		mutex: &sync.RWMutex{},
	}
}

// Save adds a todo to the repository
func (r *TodoRepository) Save(t *todo.Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.todos[t.ID] = t
	return nil
}

// GetByID retrieves a todo by its ID
func (r *TodoRepository) GetByID(id string) (*todo.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	t, exists := r.todos[id]
	if !exists {
		return nil, todo.ErrTodoNotFound
	}
	return t, nil
}

// GetAll retrieves all todos
func (r *TodoRepository) GetAll() ([]*todo.Todo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	todos := make([]*todo.Todo, 0, len(r.todos))
	for _, t := range r.todos {
		todos = append(todos, t)
	}
	return todos, nil
}

// Update updates a todo in the repository
func (r *TodoRepository) Update(t *todo.Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.todos[t.ID]; !exists {
		return todo.ErrTodoNotFound
	}

	r.todos[t.ID] = t
	return nil
}

// Delete removes a todo from the repository
func (r *TodoRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.todos[id]; !exists {
		return todo.ErrTodoNotFound
	}

	delete(r.todos, id)
	return nil
}
