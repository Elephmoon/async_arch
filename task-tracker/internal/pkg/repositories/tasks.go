package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"task-tracker/internal/pkg/models"
)

type Tasks struct {
	db *pgx.Conn
}

func NewTasks(db *pgx.Conn) *Tasks {
	return &Tasks{db: db}
}

func (h *Tasks) GetOpenedTasks(ctx context.Context) ([]models.Task, error) {
	const sqlQuery = `
		select id, name, description, price, fee, created_at
		from tasks	
		where closed_at is null
	`
	rows, err := h.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Cost, &task.Fee, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (h *Tasks) AssigneeTask(ctx context.Context, tx pgx.Tx, taskID, userID uuid.UUID) error {
	const sqlQuery = `update tasks set user_id = $1 where id = $2`
	_, err := tx.Exec(ctx, sqlQuery, userID, taskID)
	return err
}

func (h *Tasks) GetTaskExecutor(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error) {
	const sqlQuery = `select user_id from tasks where id = $1`
	row := h.db.QueryRow(ctx, sqlQuery, taskID)
	var executorID uuid.UUID
	err := row.Scan(&executorID)
	return executorID, err
}

func (h *Tasks) CloseTask(ctx context.Context, taskID uuid.UUID) (models.Task, error) {
	const sqlQuery = `
		update tasks set closed_at = now() where id = $1 and closed_at is null
		returning id, name, description, user_id, price, fee, created_at, now()
	`
	row := h.db.QueryRow(ctx, sqlQuery, taskID)
	var task models.Task
	err := row.Scan(&task.ID,
		&task.Name,
		&task.Description,
		&task.UserID,
		&task.Cost,
		&task.Fee,
		&task.CreatedAt,
		&task.ClosedAt,
	)
	return task, err
}

func (h *Tasks) CreateTask(ctx context.Context, task models.Task) error {
	const sqlQuery = `
		insert into tasks (id, name, description, user_id, price, fee, public_id)
		values ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := h.db.Exec(ctx, sqlQuery, task.ID, task.Name, task.Description, task.UserID, task.Cost, task.Fee, task.PublicID)

	return err
}

func (h *Tasks) GetUserTasks(ctx context.Context, userID uuid.UUID) ([]models.Task, error) {
	const sqlQuery = `
		select id, name, description, user_id, price, fee, created_at
		from tasks	
		where closed_at is null and user_id = $1
	`
	rows, err := h.db.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.UserID, &task.Cost, &task.Fee, &task.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
