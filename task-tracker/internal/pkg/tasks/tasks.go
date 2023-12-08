package tasks

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twmb/franz-go/pkg/kgo"
	"math/rand"
	"task-tracker/internal/pkg/models"
	"task-tracker/internal/pkg/repositories"
	"time"
)

const (
	minTaskFee   = -20
	maxTaskFee   = -10
	minTaskPrice = 20
	maxTaskPrice = 40
)

const (
	taskPricingTopic = "task-pricing"
	tasksTopic       = "tasks"
)

type UsersRepo interface {
	GetActiveUserIDs(ctx context.Context) ([]uuid.UUID, error)
	GetUserByName(ctx context.Context, name string) (models.User, error)
}

type RepoTasks interface {
	GetOpenedTasks(ctx context.Context) ([]models.Task, error)
	AssigneeTask(ctx context.Context, tx pgx.Tx, taskID, userID uuid.UUID) error
	GetTaskExecutor(ctx context.Context, taskID uuid.UUID) (uuid.UUID, error)
	CloseTask(ctx context.Context, taskID uuid.UUID) (models.Task, error)
	CreateTask(ctx context.Context, task models.Task) error
	GetUserTasks(ctx context.Context, userID uuid.UUID) ([]models.Task, error)
}

type UseCase struct {
	db          *pgx.Conn
	kafkaClient *kgo.Client
	usersRepo   UsersRepo
	repoTasks   RepoTasks
}

func New(db *pgx.Conn, kafkaClient *kgo.Client) *UseCase {
	return &UseCase{
		db:          db,
		kafkaClient: kafkaClient,
		usersRepo:   repositories.NewUsers(db),
		repoTasks:   repositories.NewTasks(db),
	}
}

func (h *UseCase) CreateTask(ctx context.Context, userName, taskName, description string) (uuid.UUID, error) {
	user, err := h.usersRepo.GetUserByName(ctx, userName)
	if err != nil {
		return uuid.UUID{}, err
	}
	if user.Role == "manager" || user.Role == "admin" {
		return uuid.UUID{}, fmt.Errorf("invalid user for task assignee")
	}

	taskID := uuid.New()
	fee, price := calcTaskFeePrice()
	task := models.Task{
		ID:          taskID,
		PublicID:    uuid.New(),
		UserID:      user.ID,
		Name:        taskName,
		Description: description,
		Cost:        price,
		Fee:         fee,
		CreatedAt:   time.Now(),
	}

	err = h.repoTasks.CreateTask(ctx, task)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = h.publishAssigneeTaskEvents(ctx, task)
	if err != nil {
		return uuid.UUID{}, err
	}

	return taskID, nil
}

func (h *UseCase) CloseTask(ctx context.Context, taskID, userID uuid.UUID) error {
	// в реальном проекте действие по закрытию задачи нужно совершать транзакционно с локом задачи
	// чтобы параллельно задача не попала в шафлинг от попуга менеджера с горящими глазами
	executorID, err := h.repoTasks.GetTaskExecutor(ctx, taskID)
	if err != nil {
		return err
	}
	// чтобы попуг не отъехал в дурку из за того что кто то закрывает за него задачи,
	// провалидируем что закрывать можно только свои задачи
	// вероятно в будущем нужно будет разрешить закрывать таски попугу менеджеру с горящими глазами
	if executorID != userID {
		return fmt.Errorf("this user cant close this task")
	}

	task, err := h.repoTasks.CloseTask(ctx, taskID)
	if err != nil {
		return err
	}

	return h.publishClosingTaskEvents(ctx, task)
}

func (h *UseCase) Shuffle(ctx context.Context, userRole string) error {
	if userRole != "admin" && userRole != "manager" {
		return fmt.Errorf("invalid role for shuffling")
	}
	userIDs, err := h.usersRepo.GetActiveUserIDs(ctx)
	if err != nil {
		return err
	}
	tasks, err := h.repoTasks.GetOpenedTasks(ctx)
	if err != nil {
		return err
	}

	return h.shuffle(ctx, userIDs, tasks)
}

func (h *UseCase) GetUserTasks(ctx context.Context, userID uuid.UUID) ([]models.Task, error) {
	return h.repoTasks.GetUserTasks(ctx, userID)
}

func (h *UseCase) shuffle(ctx context.Context, userIDs []uuid.UUID, tasks []models.Task) error {
	if len(userIDs) == 0 {
		return nil
	}

	var events []*kgo.Record
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i := range tasks {
		shuffleID := rand.Intn(len(userIDs))
		newUserID := userIDs[shuffleID]
		tasks[i].UserID = newUserID

		err2 := h.repoTasks.AssigneeTask(ctx, tx, tasks[i].ID, newUserID)
		if err2 != nil {
			return err2
		}

		taskEvent, err2 := models.NewTaskEvent(tasks[i])
		if err2 != nil {
			return err2
		}
		taskEvent.Topic = tasksTopic

		events = append(events, &taskEvent)
		tasPricingEvent, err2 := models.NewTaskPricingEvent(tasks[i].ID, newUserID, tasks[i].Fee)
		if err2 != nil {
			return err2
		}
		tasPricingEvent.Topic = taskPricingTopic
		events = append(events, &tasPricingEvent)
	}

	// не стоит ходить в реальных проектах в открытой транзакции во внешний сервис
	err = h.kafkaClient.ProduceSync(ctx, events...).FirstErr()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func calcTaskFeePrice() (int64, int64) {
	taskFee := rand.Intn(maxTaskFee-minTaskFee) + minTaskFee
	taskPrice := rand.Intn(maxTaskPrice-minTaskPrice) + minTaskPrice

	return int64(taskFee), int64(taskPrice)
}

func (h *UseCase) publishAssigneeTaskEvents(ctx context.Context, task models.Task) error {
	taskEvent, err := models.NewTaskEvent(task)
	if err != nil {
		return err
	}
	taskEvent.Topic = tasksTopic

	taskPricingEvent, err := models.NewTaskPricingEvent(task.ID, task.UserID, task.Fee)
	if err != nil {
		return err
	}
	taskPricingEvent.Topic = taskPricingTopic

	return h.kafkaClient.ProduceSync(ctx, &taskEvent, &taskPricingEvent).FirstErr()
}

func (h *UseCase) publishClosingTaskEvents(ctx context.Context, task models.Task) error {
	taskEvent, err := models.NewTaskEvent(task)
	if err != nil {
		return err
	}
	taskEvent.Topic = tasksTopic

	taskPricingEvent, err := models.NewTaskPricingEvent(task.ID, task.UserID, task.Cost)
	if err != nil {
		return err
	}
	taskPricingEvent.Topic = taskPricingTopic

	return h.kafkaClient.ProduceSync(ctx, &taskEvent, &taskPricingEvent).FirstErr()
}
