package models

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	Cost        int64
	Fee         int64
	CreatedAt   time.Time
	ClosedAt    *time.Time
}
