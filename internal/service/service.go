package service

import (
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/internal/storage"
	"tender-bridge/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	User
	Authorization
}

func NewService(repos *repository.Repository, storage *storage.Storage, cfg *config.Config, loggers *logger.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos, loggers, cfg),
		User:          NewUserService(repos, loggers),
	}
}

type User interface {
	CreateUser(request models.CreateUser) (uuid.UUID, error)
	GetUsers(filter models.UserFilter) ([]models.User, int, error)
	GetUser(id uuid.UUID) (models.User, error)
	UpdateUser(request models.UpdateUser) error
	DeleteUser(id uuid.UUID) error
}

type Authorization interface {
	CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error)
	GenerateTokens(user models.User) (*models.Token, *models.Token, error)
	ParseToken(token string) (*jwtCustomClaim, error)
	Login(request models.Login) (*models.Token, *models.Token, error)
	Register(request models.Register) (*models.Token, *models.Token, error)
}
