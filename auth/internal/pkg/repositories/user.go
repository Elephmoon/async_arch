package repositories

import (
	"auth/internal/pkg/models"
	"context"
	"github.com/jackc/pgx/v5"
)

type User struct {
	db *pgx.Conn
}

func NewUser(db *pgx.Conn) *User {
	return &User{db: db}
}

func (h *User) Insert(ctx context.Context, user models.User) (models.User, error) {
	const sqlQuery = `insert into users (id, name, password, role, created_at) values($1, $2, $3, $4, $5)`

	_, err := h.db.Exec(ctx, sqlQuery, user.ID, user.Name, user.Password, user.Role, user.CreatedAt)

	return user, err
}

func (h *User) GetUser(ctx context.Context, userName string) (*models.User, error) {
	const sqlQuery = `select id, name, role, password from users where name = $1 and deleted_at is null`

	var user models.User
	row := h.db.QueryRow(ctx, sqlQuery, userName)
	err := row.Scan(&user.ID, &user.Name, &user.Role, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
