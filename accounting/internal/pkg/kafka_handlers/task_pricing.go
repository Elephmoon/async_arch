package kafka_handlers

import (
	"accounting/internal/pkg/models"
	"accounting/schema_registry/task_pricing"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/encoding/protojson"
)

type TaskPricingKafkaHandler struct {
	db          *pgx.Conn
	kafkaClient *kgo.Client
}

func NewTaskPricingKafkaHandler(db *pgx.Conn, kafkaClient *kgo.Client) *TaskPricingKafkaHandler {
	return &TaskPricingKafkaHandler{
		db:          db,
		kafkaClient: kafkaClient,
	}
}

func (h *TaskPricingKafkaHandler) Handle(ctx context.Context, msg []byte) error {
	event := task_pricing.Event{}
	// todo если мы не распарсили сообщение стоит положить его в dlq,
	// т.к скорее всего проблема в схеме сообщения и нет смысла его ретраить
	// и стопать консюминг, проще положить в dlq и проалёртить
	// чтобы человек посмотрел и разобрался что пошло не так
	err := protojson.Unmarshal(msg, &event)
	if err != nil {
		return err
	}

	return h.processEvent(ctx, &event)
}

func (h *TaskPricingKafkaHandler) processEvent(ctx context.Context, event *task_pricing.Event) error {
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// todo если мы не распарсили task_id стоит положить всё сообщение в dlq,
	// т.к скорее всего проблема в нарушении контракта и нет смысла ретраить сообщение
	// и стопать консюминг, проще положить в dlq и проалёртить
	// чтобы человек посмотрел и разобрался что пошло не так
	taskID, err := uuid.Parse(event.GetTaskId())
	if err != nil {
		return err
	}

	// лочим счёт
	sqlQuery := `select number, balance from accounts where user_id = $1 for update`
	var account models.Account
	row := tx.QueryRow(ctx, sqlQuery, event.GetUserId())
	err = row.Scan(&account.Number, account.Balance)
	if err != nil {
		return err
	}

	// проверяем не обрабатывали ли мы этот запрос раньше
	sqlQuery = `select count(*) from audit_log where idempotency_key = $1`
	var idempotencyCount int64
	row = tx.QueryRow(ctx, sqlQuery, event.GetIdempotencyKey())
	err = row.Scan(&idempotencyCount)
	if idempotencyCount != 0 {
		return tx.Commit(ctx)
	}

	// изменяем баланс
	finalBalance := account.Balance + event.GetCost()
	sqlQuery = `update accounts set balance = $1 where number = $2`
	_, err = tx.Exec(ctx, sqlQuery, finalBalance, account.Number)

	auditLog := models.AuditLog{
		ID:            uuid.New(),
		Event:         models.AuditLogEventBalanceChanged,
		AccountNumber: account.Number,
		TaskID:        &taskID,
		ChangeAmount:  event.GetCost(),
	}

	// сохраняем аудит лог
	sqlQuery = `insert into audit_log (id, account_number, task_id, change_amount, event, idempotency_key)
		values ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, sqlQuery, auditLog.ID, auditLog.AccountNumber, auditLog.TaskID,
		auditLog.ChangeAmount, auditLog.Event, event.GetIdempotencyKey())
	if err != nil {
		return err
	}

	// паблишим аудит лог в кафку
	auditLogEvent, err := models.NewAuditLogKafkaEvent(auditLog)
	if err != nil {
		return err
	}

	err = h.kafkaClient.ProduceSync(ctx, &auditLogEvent).FirstErr()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
