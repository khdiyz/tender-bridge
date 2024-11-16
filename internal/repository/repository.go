package repository

import (
	"tender-bridge/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) *Repository {
	return &Repository{}
}
