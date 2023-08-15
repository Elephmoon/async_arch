package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type KafkaUserEvent struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type TaskPricingEvent struct {
	TaskID         uuid.UUID `json:"task_id"`
	IdempotencyKey uuid.UUID `json:"idempotency_key"`
	UserID         uuid.UUID `json:"user_id"`
	Cost           int64     `json:"cost"`
}

func NewTaskPricingEvent(taskID, userID uuid.UUID, cost int64) (kgo.Record, error) {
	event := TaskPricingEvent{
		TaskID:         taskID,
		IdempotencyKey: uuid.New(),
		UserID:         userID,
		Cost:           cost,
	}
	eventJson, err := json.Marshal(event)
	if err != nil {
		return kgo.Record{}, err
	}

	return kgo.Record{
		Key:   taskID[:],
		Value: eventJson,
	}, nil
}

type TaskEvent struct {
	TaskID      uuid.UUID  `json:"task_id"`
	UserID      uuid.UUID  `json:"user_id"`
	Name        string     `json:"name"`
	Fee         int64      `json:"fee"`
	Cost        int64      `json:"cost"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at,omitempty"`
}

func NewTaskEvent(task Task) (kgo.Record, error) {
	event := TaskEvent{
		TaskID:      task.ID,
		UserID:      task.UserID,
		Name:        task.Name,
		Fee:         task.Fee,
		Cost:        task.Cost,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		ClosedAt:    task.ClosedAt,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return kgo.Record{}, err
	}
	return kgo.Record{
		Key:   task.ID[:],
		Value: jsonEvent,
	}, nil
}
