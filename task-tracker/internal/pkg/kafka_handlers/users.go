package kafka_handlers

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"task-tracker/internal/pkg/models"
	"time"
)

type UsersKafkaHandler struct {
	db *pgx.Conn
}

func NewUsersKafkaHandler(db *pgx.Conn) *UsersKafkaHandler {
	return &UsersKafkaHandler{db: db}
}

func (h *UsersKafkaHandler) Handle(ctx context.Context, msg []byte) error {
	userEvent := models.KafkaUserEvent{}
	err := json.Unmarshal(msg, &userEvent)
	if err != nil {
		return err
	}
	return h.upsertUser(ctx, userEvent)
}

func (h *UsersKafkaHandler) upsertUser(ctx context.Context, userEvent models.KafkaUserEvent) error {
	// todo добавить механизм ордеринга событий
	const sqlQuery = `
		insert into users (id, name, role, created_at, synced_at, deleted_at) 
		values ($1, $2, $3, $4, $5, $6)
		on conflict (name) 
		do update set 
		    synced_at = now(),
			role = $3,
			deleted_at = $6 
		where users.id = $1 and users.name = $2
		`
	_, err := h.db.Exec(ctx, sqlQuery,
		userEvent.ID, userEvent.Name, userEvent.Role, userEvent.CreatedAt,
		time.Now(), userEvent.DeletedAt)

	return err
}
