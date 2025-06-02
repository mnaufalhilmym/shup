package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID          uuid.UUID
	FileName    string
	ContentType string
	Size        int64
	CreatedAt   time.Time
	ExpiredAt   time.Time
}
