package models

import "github.com/google/uuid"

type TaskResponse struct {
	TaskID uuid.UUID `json:"task_id"`
}
