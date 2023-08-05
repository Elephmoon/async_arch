package handlers

import (
	"auth/internal/pkg/repositories"
	"auth/internal/pkg/users"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
	"net/http"
	"os"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, name, role, pass string) (string, error)
	LoginUser(ctx context.Context, name, pass string) (string, error)
}

type Service struct {
	userUseCase UserUseCase
}

func NewService(db *pgx.Conn, kafkaClient *kgo.Client) *Service {
	usersRepo := repositories.NewUser(db)

	return &Service{
		userUseCase: users.New(usersRepo, kafkaClient),
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
