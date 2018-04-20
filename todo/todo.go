package todo

import (
	"time"
)

type ValidationError string

func (e ValidationError) Error() string { return "validation error: " + string(e) }

// Task represents a to-do item
type Task struct {
	ID        string    `json:"id"`
	Done      bool      `json:"done"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
