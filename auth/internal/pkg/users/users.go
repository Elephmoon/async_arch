package users

import (
	jwt2 "auth/internal/pkg/jwt"
	"auth/internal/pkg/models"
	"auth/internal/pkg/password"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"time"
)

type Repo interface {
	Insert(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, userName string) (*models.User, error)
}

func New(repo Repo, kafkaClient *kgo.Client) *UseCase {
	return &UseCase{
		repo:        repo,
		kafkaClient: kafkaClient,
	}
}

type UseCase struct {
	repo        Repo
	kafkaClient *kgo.Client
}

func (h *UseCase) CreateUser(ctx context.Context, name, role, pass string) (string, error) {
	passEncrypt := password.NewDefaultPassEncrypt()
	pass, err := passEncrypt.EncryptPassword(pass)
	if err != nil {
		return "", fmt.Errorf("gen pass error")
	}

	createdUser, err := h.repo.Insert(ctx, models.User{
		ID:        uuid.New(),
		PublicID:  uuid.New(),
		Name:      name,
		Role:      role,
		Password:  pass,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", fmt.Errorf("cant create user %w", err)
	}

	event, err := models.NewKafkaUserEvent(createdUser)
	if err != nil {
		return "", fmt.Errorf("cant create event %w", err)
	}

	err = h.kafkaClient.ProduceSync(ctx, &event).FirstErr()
	if err != nil {
		return "", fmt.Errorf("cant publish event %w", err)
	}

	token, err := jwt2.NewJWTForUser(createdUser)
	if err != nil {
		return "", fmt.Errorf("cant geenerate token %w", err)
	}

	return token, nil
}

func (h *UseCase) LoginUser(ctx context.Context, name, pass string) (string, error) {
	user, err := h.repo.GetUser(ctx, name)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", nil
	}

	passEncrypt := password.NewDefaultPassEncrypt()
	valid, err := passEncrypt.ValidatePass(pass, user.Password)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	if !valid {
		return "", fmt.Errorf("invalid password")
	}

	token, err := jwt2.NewJWTForUser(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}
