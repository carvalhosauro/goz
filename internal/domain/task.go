package domain

import "errors"

var (
	ErrEmptyText       = errors.New("task text cannot be empty")
	ErrInvalidPriority = errors.New("invalid priority")
)

// Task is the central domain entity. Phase 0 keeps Estimate and Due as free-form
// strings to mirror the design mock; later phases will turn them into typed values.
type Task struct {
	ID        string
	Text      string
	Done      bool
	Priority  Priority
	Estimate  string
	Due       string
	Tag       Tag
	ParentID  string
	SortOrder int
	Overdue   bool
}

func (t Task) Validate() error {
	if t.Text == "" {
		return ErrEmptyText
	}
	if !t.Priority.Valid() {
		return ErrInvalidPriority
	}
	return nil
}
