package models

import (
	"auth/schema_registry/user_lifecycle"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const userLifecycleTopicName = "user-lifecycle"

type User struct {
	ID        uuid.UUID
	PublicID  uuid.UUID
	Name      string
	Role      string
	CreatedAt time.Time
	DeletedAt *time.Time
	Password  string
}

func NewKafkaUserEvent(user User) (kgo.Record, error) {
	var deletedAt *timestamppb.Timestamp
	if user.DeletedAt != nil {
		deletedAt = timestamppb.New(*user.DeletedAt)
	}
	userEvent := user_lifecycle.Event{
		Id:        user.ID.String(),
		PublicId:  user.PublicID.String(),
		Name:      user.Name,
		Role:      user.Role,
		DeletedAt: deletedAt,
	}

	pl, err := protojson.Marshal(&userEvent)
	if err != nil {
		return kgo.Record{}, err
	}

	return kgo.Record{
		Key:   user.ID[:],
		Value: pl,
		Topic: userLifecycleTopicName,
	}, nil
}
