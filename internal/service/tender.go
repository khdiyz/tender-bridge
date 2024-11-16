package service

import (
	"errors"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type tenderService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func NewTenderService(repo *repository.Repository, logger *logger.Logger) *tenderService {
	return &tenderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *tenderService) CreateTender(request models.CreateTender) (uuid.UUID, error) {
	request.Status = config.TenderStatusOpen

	id, err := s.repo.Tender.Create(request)
	if err != nil {
		return uuid.Nil, serviceError(err, codes.Internal)
	}

	return id, nil
}

func (s *tenderService) GetTenders(filter models.TenderFilter) ([]models.Tender, int, error) {
	tenders, total, err := s.repo.Tender.GetList(filter)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	return tenders, total, nil
}

func (s *tenderService) GetTender(id uuid.UUID) (models.Tender, error) {
	tender, err := s.repo.Tender.GetById(id)
	if err != nil {
		return models.Tender{}, serviceError(err, codes.Internal)
	}

	return tender, nil
}

func (s *tenderService) UpdateTender(request models.UpdateTender) error {
	if request.Status != config.TenderStatusAwarded && request.Status != config.TenderStatusClosed && request.Status != config.TenderStatusOpen {
		return serviceError(errors.New("invalid tender status"), codes.InvalidArgument)
	}

	if err := s.repo.Tender.Update(request); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *tenderService) DeleteTender(id uuid.UUID) error {
	if err := s.repo.Tender.Delete(id); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}
