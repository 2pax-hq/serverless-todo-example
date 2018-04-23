package todo

import (
	"time"
)

// ValidationError represents an error while adding a task.
type ValidationError string

func (e ValidationError) Error() string { return "validation failed: " + string(e) }

// UnknownTaskError string represents a failure to find the specified task in
// the store.
type UnknownTaskError string

func (e UnknownTaskError) Error() string { return "unknown task: " + string(e) }

// Task represents a to-do item
type Task struct {
	ID        string    `json:"id"`
	Done      bool      `json:"done"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
