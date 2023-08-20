package close_day

import (
	"accounting/internal/pkg/models"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
)

type CloseDayUseCase struct {
	db          *pgx.Conn
	kafkaClient *kgo.Client
}

func NewCloseDayUseCase(db *pgx.Conn, kafkaClient *kgo.Client) *CloseDayUseCase {
	return &CloseDayUseCase{
		db:          db,
		kafkaClient: kafkaClient,
	}
}

func (h *CloseDayUseCase) CloseDay(ctx context.Context) error {
	// todo когда нашим таск трекером начнёт пользоваться 100500 тыщ попугов стоит
	// добавить какой нибудь батчинг для адекватного закрытия дня

	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// лочим счета с положительными балансами
	sqlQuery := `select number, balance from accounts where balance > 0 for update`
	rows, err := tx.Query(ctx, sqlQuery)
	if err != nil {
		return err
	}
	var accounts []models.Account
	for rows.Next() {
		account := models.Account{}
		err = rows.Scan(&account.Number, account.Balance)
		if err != nil {
			return err
		}
		accounts = append(accounts, account)
	}
	defer rows.Close()

	var events []*kgo.Record
	for _, account := range accounts {
		changeAmount := -account.Balance
		finalBalance := account.Balance + changeAmount
		sqlQuery = `update accounts set balance = $1 where number = $2`
		_, err = tx.Exec(ctx, sqlQuery, finalBalance, account.Number)
		if err != nil {
			return err
		}
		auditLog := models.AuditLog{
			ID:            uuid.New(),
			Event:         models.AuditLogEventPaymentMade,
			AccountNumber: account.Number,
			ChangeAmount:  changeAmount,
		}
		sqlQuery = `insert into audit_log (id, account_number, change_amount, event, idempotency_key)
			values ($1, $2, $3, $4, $5)
		`
		_, err = tx.Exec(ctx, sqlQuery, auditLog.ID, auditLog.AccountNumber, auditLog.ChangeAmount, auditLog.Event,
			uuid.NewString())
		if err != nil {
			return err
		}
		auditLogEvent, err2 := models.NewAuditLogKafkaEvent(auditLog)
		if err2 != nil {
			return err2
		}
		events = append(events, &auditLogEvent)
		// TODO планировать отправку email с выплатой
	}

	err = h.kafkaClient.ProduceSync(ctx, events...).FirstErr()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
