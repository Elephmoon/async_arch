package models

import (
	"accounting/schema_registry/accounting_audit_log"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
)

type AuditLogEvent string

const (
	// AuditLogEventBalanceChanged событие изменения баланса
	AuditLogEventBalanceChanged AuditLogEvent = "balance_changed"
	// AuditLogEventPaymentMade событие выплаты в конце дня
	AuditLogEventPaymentMade AuditLogEvent = "payment_made"
)

type AuditLog struct {
	ID            uuid.UUID
	Event         AuditLogEvent
	AccountNumber string
	TaskID        *uuid.UUID
	ChangeAmount  int64
}

func NewAuditLogKafkaEvent(auditLog AuditLog) (kgo.Record, error) {
	var taskID *string
	if auditLog.TaskID != nil {
		tskID := auditLog.TaskID.String()
		taskID = &tskID
	}
	event := accounting_audit_log.Event{
		Id:            auditLog.ID.String(),
		EventType:     string(auditLog.Event),
		AccountNumber: auditLog.AccountNumber,
		TaskId:        taskID,
		ChangeAmount:  auditLog.ChangeAmount,
	}

	jsonPayload, err := protojson.Marshal(&event)
	if err != nil {
		return kgo.Record{}, nil
	}

	return kgo.Record{
		Key:   auditLog.ID[:],
		Value: jsonPayload,
		Topic: "accounting-audit-log",
	}, nil
}
