package todo

// Repository defines the interface for todo data access
type Repository interface {
	Save(*Todo) error
	GetByID(id string) (*Todo, error)
	GetAll() ([]*Todo, error)
	Update(*Todo) error
	Delete(id string) error
}
