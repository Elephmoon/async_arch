package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Role      string
	CreatedAt time.Time
	DeletedAt *time.Time
	Password  string
}

type KafkaUserEvent struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func NewKafkaUserEvent(user User) (kgo.Record, error) {
	event := KafkaUserEvent{
		ID:        user.ID,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		DeletedAt: user.DeletedAt,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return kgo.Record{}, err
	}

	return kgo.Record{
		Key:   user.ID[:],
		Value: payload,
		Topic: "auth-users-info",
	}, nil
}
