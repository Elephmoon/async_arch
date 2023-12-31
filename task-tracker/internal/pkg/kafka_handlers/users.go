package kafka_handlers

import (
	"context"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/encoding/protojson"
	"task-tracker/schema_registry/user_lifecycle"
	"time"
)

type UsersKafkaHandler struct {
	db *pgx.Conn
}

func NewUsersKafkaHandler(db *pgx.Conn) *UsersKafkaHandler {
	return &UsersKafkaHandler{db: db}
}

func (h *UsersKafkaHandler) Handle(ctx context.Context, msg []byte) error {
	event := user_lifecycle.Event{}
	// todo если мы не распарсили сообщение стоит пололожить его в dlq,
	// т.к скорее всего проблема в схеме сообщения и нет смысла его ретраить
	// и стопать консюминг, проще положить в dlq и проалёртить
	// чтобы человек посмотрел и разобрался что пошло не так
	err := protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(msg, &event)
	if err != nil {
		return err
	}

	return h.upsertUser(ctx, &event)
}

func (h *UsersKafkaHandler) upsertUser(ctx context.Context, userEvent *user_lifecycle.Event) error {
	var deletedAt *time.Time
	if userEvent.GetDeletedAt() != nil {
		tm := userEvent.GetDeletedAt().AsTime()
		deletedAt = &tm
	}

	const sqlQuery = `
		insert into users (id, public_id, name, role, created_at, deleted_at) 
		values ($1, $2, $3, $4, now(), $5)
		on conflict (name) 
		do update set 
			role = $3,
			deleted_at = $6 
		where users.id = $1 and users.name = $3 and public_id = $2
		`
	_, err := h.db.Exec(ctx, sqlQuery,
		userEvent.GetId(), userEvent.GetPublicId(), userEvent.GetName(),
		userEvent.GetRole(), deletedAt)

	return err
}
