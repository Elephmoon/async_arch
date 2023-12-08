package handlers

import (
	"accounting/internal/pkg/close_day"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
)

type DayCloser interface {
	CloseDay(ctx context.Context) error
}

type Service struct {
	dayCloser DayCloser
}

func NewService(db *pgx.Conn, kafkaClient *kgo.Client) *Service {
	return &Service{
		dayCloser: close_day.NewCloseDayUseCase(db, kafkaClient),
	}
}
