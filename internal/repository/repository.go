package repository

import (
	"tender-bridge/internal/models"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) *Repository {
	return &Repository{
		User: NewUserRepo(db, logger),
	}
}

type User interface {
	Create(request models.CreateUser) (uuid.UUID, error)
	GetList(filter models.UserFilter) ([]models.User, int, error)
	GetById(id uuid.UUID) (models.User, error)
	Update(request models.UpdateUser) error
	Delete(id uuid.UUID) error

	GetByUsername(username string) (models.User, error)
}
