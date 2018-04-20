package todo

import (
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

// MockStore allows to add tasks but does not store them anywhere.
type MockStore struct{}

// Add creates a new task.
func (s MockStore) Add(note string) (Task, error) {
	if note == "" {
		return Task{}, ValidationError(fmt.Sprintf("invalid note value: `%s`", note))
	}
	now := time.Now().UTC()
	id, _ := uuid()

	return Task{
		ID:        id,
		Note:      note,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// curtesy of https://play.golang.org/p/4FkNSiUDMg
func uuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
