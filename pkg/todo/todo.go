package todo

import (
	"time"
)

type ValidationError string

func (e ValidationError) Error() string { return "validation error: " + string(e) }

type Task struct {
	ID        string    `json:"id"`
	Checked   bool      `json:"checked"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
