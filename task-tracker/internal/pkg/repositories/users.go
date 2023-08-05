package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"task-tracker/internal/pkg/models"
)

func NewUsers(db *pgx.Conn) *Users {
	return &Users{db: db}
}

type Users struct {
	db *pgx.Conn
}

func (h *Users) GetActiveUserIDs(ctx context.Context) ([]uuid.UUID, error) {
	const sqlQuery = `select id from users 
          where role != 'admin' and role != 'manager' and deleted_at is null`

	rows, err := h.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

func (h *Users) GetUserByName(ctx context.Context, name string) (models.User, error) {
	sqlQuery := `select id, role, name from users where name = $1 and deleted_at is null`
	res := h.db.QueryRow(ctx, sqlQuery, name)
	var user models.User
	err := res.Scan(&user.ID, &user.Role, &user.Name)

	return user, err
}
