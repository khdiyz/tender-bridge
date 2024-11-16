package repository

import (
	"tender-bridge/internal/models"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User
	Tender
	Bid
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) *Repository {
	return &Repository{
		User:   NewUserRepo(db, logger),
		Tender: NewTenderRepo(db, logger),
		Bid:    NewBidRepo(db, logger),
	}
}

type User interface {
	Create(request models.CreateUser) (uuid.UUID, error)
	GetList(filter models.UserFilter) ([]models.User, int, error)
	GetById(id uuid.UUID) (models.User, error)
	Update(request models.UpdateUser) error
	Delete(id uuid.UUID) error
	GetByUsername(username string) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetByIds(ids []uuid.UUID) ([]models.User, error)
}

type Tender interface {
	Create(request models.CreateTender) (uuid.UUID, error)
	GetList(filter models.TenderFilter) ([]models.Tender, int, error)
	GetById(id uuid.UUID) (models.Tender, error)
	Update(request models.UpdateTender) error
	Delete(id uuid.UUID) error
}

type Bid interface {
	Create(request models.CreateBid) (uuid.UUID, error)
	GetList(filter models.BidFilter) ([]models.Bid, int, error)
	GetById(id uuid.UUID) (models.Bid, error)
	Update(request models.UpdateBid) error
	Delete(id uuid.UUID) error
}
