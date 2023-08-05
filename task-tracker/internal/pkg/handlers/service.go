package handlers

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
	"net/http"
	"os"
	"task-tracker/internal/pkg/models"
	"task-tracker/internal/pkg/tasks"
)

type TasksUseCase interface {
	CreateTask(ctx context.Context, userName, taskName, description string) (uuid.UUID, error)
	CloseTask(ctx context.Context, taskID, userID uuid.UUID) error
	Shuffle(ctx context.Context, userRole string) error
	GetUserTasks(ctx context.Context, userID uuid.UUID) ([]models.Task, error)
}

type Service struct {
	tasksUseCase TasksUseCase
}

func NewService(db *pgx.Conn, kafkaClient *kgo.Client) *Service {
	return &Service{
		tasksUseCase: tasks.New(db, kafkaClient),
	}
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
