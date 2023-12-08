package kafka_handlers

import (
	"accounting/schema_registry/user_lifecycle"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/encoding/protojson"
)

type UsersKafkaHandler struct {
	db *pgx.Conn
}

func NewUsersKafkaHandler(db *pgx.Conn) *UsersKafkaHandler {
	return &UsersKafkaHandler{db: db}
}

func (h *UsersKafkaHandler) Handle(ctx context.Context, msg []byte) error {
	event := user_lifecycle.Event{}
	// todo если мы не распарсили сообщение стоит положить его в dlq,
	// т.к скорее всего проблема в схеме сообщения и нет смысла его ретраить
	// и стопать консюминг, проще положить в dlq и проалёртить
	// чтобы человек посмотрел и разобрался что пошло не так
	err := protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(msg, &event)
	if err != nil {
		return err
	}
	return h.createUserAccount(ctx, &event)
}

func (h *UsersKafkaHandler) createUserAccount(ctx context.Context, event *user_lifecycle.Event) error {
	const sqlQuery = `
		insert into accounts(number, public_id, user_id, user_public_id)
		values ($1, $2, $3, $4)
		on conflict (user_id, user_public_id)
		do nothing;
	`
	accountNumber := uuid.NewString()
	accountPublicID := uuid.NewString()
	_, err := h.db.Exec(ctx, sqlQuery, accountNumber, accountPublicID, event.GetId(), event.GetPublicId())

	return err
}
