package store

import (
	"errors"

	"goz/internal/domain"
)

var ErrNotFound = errors.New("task not found")

// Store is the persistence boundary. Phase 0 has only an in-memory implementation;
// Phase 1 swaps in SQLite without changing the rest of the app.
type Store interface {
	List() []domain.Task
	Get(id string) (domain.Task, bool)
	Add(text string) (domain.Task, error)
	Toggle(id string) error
	Delete(id string) error
	Move(id string, delta int) error
}
