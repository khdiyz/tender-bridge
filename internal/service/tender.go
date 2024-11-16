package service

import (
	"database/sql"
	"errors"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/logger"
	"time"

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
	deadlineTime, err := time.Parse(time.RFC3339, request.Deadline)
	if err != nil {
		return uuid.Nil, serviceError(err, codes.Internal)
	}

	isDeadlineBefore := deadlineTime.Before(time.Now())
	isBudgetNegative := request.Budget < 0

	if isDeadlineBefore || isBudgetNegative {
		return uuid.Nil, serviceError(errors.New("error: Invalid tender data"), codes.InvalidArgument)
	}

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
	_, err := s.repo.Tender.GetById(id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return serviceError(err, codes.Internal)
	} else if errors.Is(err, sql.ErrNoRows) {
		return serviceError(errors.New("Tender not found or access denied"), codes.NotFound)
	}

	if err := s.repo.Tender.Delete(id); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *tenderService) UpdateTenderStatus(request models.UpdateTenderStatus) error {
	if request.Status != config.TenderStatusAwarded && request.Status != config.TenderStatusClosed && request.Status != config.TenderStatusOpen {
		return serviceError(errors.New("error: Invalid tender status"), codes.InvalidArgument)
	}

	tender, err := s.repo.Tender.GetById(request.Id)
	if err != nil {
		return serviceError(err, codes.InvalidArgument)
	}

	if err := s.repo.Tender.Update(models.UpdateTender{
		Id:          request.Id,
		ClientId:    tender.ClientId,
		Title:       tender.Title,
		Description: tender.Description,
		Deadline:    tender.Deadline,
		Budget:      tender.Budget,
		File:        tender.File,
		Status:      request.Status,
	}); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}
