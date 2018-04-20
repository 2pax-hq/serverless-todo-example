package todo

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type MockStore struct{}

func (s MockStore) Add(note string) (Task, error) {

	if note == "" {
		return Task{}, ValidationError(fmt.Sprintf("invalid note value: `%s`", note))
	}
	now := time.Now().UTC()
	return Task{
		ID:        uuid.NewV4().String(),
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
