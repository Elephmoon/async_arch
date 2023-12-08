package models

import (
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	"task-tracker/schema_registry/task_lifecycle"
	"task-tracker/schema_registry/task_pricing"
)

func NewTaskPricingEvent(taskID, userID uuid.UUID, cost int64) (kgo.Record, error) {
	event := task_pricing.Event{
		TaskId:         taskID.String(),
		UserId:         userID.String(),
		Cost:           cost,
		IdempotencyKey: uuid.NewString(),
	}
	eventJson, err := protojson.Marshal(&event)
	if err != nil {
		return kgo.Record{}, err
	}

	return kgo.Record{
		Key:   taskID[:],
		Value: eventJson,
	}, nil
}

func NewTaskEvent(task Task) (kgo.Record, error) {
	var closedAt *timestamppb.Timestamp
	if task.ClosedAt != nil {
		closedAt = timestamppb.New(*task.ClosedAt)
	}
	event := task_lifecycle.Event{
		Id:          task.ID.String(),
		PublicId:    task.PublicID.String(),
		UserId:      task.UserID.String(),
		Name:        task.Name,
		JiraId:      nil,
		Fee:         task.Fee,
		Cost:        task.Cost,
		Description: task.Description,
		ClosedAt:    closedAt,
	}

	eventJson, err := protojson.Marshal(&event)
	if err != nil {
		return kgo.Record{}, err
	}

	return kgo.Record{
		Key:   task.ID[:],
		Value: eventJson,
	}, nil
}
