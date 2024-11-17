package service

import (
	"tender-bridge/config"
	"tender-bridge/internal/cache"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	User
	Authorization
	Tender
	Bid
}

func NewService(repos *repository.Repository, cache *cache.RedisCache, cfg *config.Config, loggers *logger.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos, loggers, cfg),
		User:          NewUserService(repos, loggers),
		Tender:        NewTenderService(repos, cache, loggers),
		Bid:           NewBidService(repos, loggers),
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

type Tender interface {
	CreateTender(request models.CreateTender) (uuid.UUID, error)
	GetTenders(filter models.TenderFilter) ([]models.Tender, int, error)
	GetTender(id uuid.UUID) (models.Tender, error)
	UpdateTender(request models.UpdateTender) error
	DeleteTender(id uuid.UUID) error
	UpdateTenderStatus(request models.UpdateTenderStatus) error
}

type Bid interface {
	SubmitBid(request models.CreateBid) (uuid.UUID, error)
	GetBids(filter models.BidFilter) ([]models.Bid, int, error)
	GetBid(id uuid.UUID) (models.Bid, error)
	UpdateBid(request models.UpdateBid) error
	DeleteContractorBid(contractorId, bidId uuid.UUID) error
	AwardBid(clientId, tenderId, bidId uuid.UUID) error
}
