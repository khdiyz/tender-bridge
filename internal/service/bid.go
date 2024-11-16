package service

import (
	"database/sql"
	"errors"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type bidService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func NewBidService(repo *repository.Repository, logger *logger.Logger) *bidService {
	return &bidService{
		repo:   repo,
		logger: logger,
	}
}

func (s *bidService) CreateBid(request models.CreateBid) (uuid.UUID, error) {
	if request.Price <= 0 || request.DeliveryTime <= 0 || request.Comment == "" {
		return uuid.Nil, serviceError(errors.New("error: Invalid bid data"), codes.InvalidArgument)
	}

	tender, err := s.repo.Tender.GetById(request.TenderId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, serviceError(err, codes.Internal)
	} else if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, serviceError(errors.New("Tender not found"), codes.NotFound)
	}

	if tender.Status != config.TenderStatusOpen {
		return uuid.Nil, serviceError(errors.New("Tender is not open for bids"), codes.InvalidArgument)
	}

	request.Status = config.BidStatusPending

	id, err := s.repo.Bid.Create(request)
	if err != nil {
		return uuid.Nil, serviceError(err, codes.Internal)
	}

	return id, nil
}

func (s *bidService) GetBids(filter models.BidFilter) ([]models.Bid, int, error) {
	bids, total, err := s.repo.Bid.GetList(filter)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	return bids, total, nil
}

func (s *bidService) GetBid(id uuid.UUID) (models.Bid, error) {
	bid, err := s.repo.Bid.GetById(id)
	if err != nil {
		return models.Bid{}, serviceError(err, codes.Internal)
	}

	return bid, nil
}

func (s *bidService) UpdateBid(request models.UpdateBid) error {
	if err := s.repo.Bid.Update(request); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *bidService) DeleteContractorBid(contractorId, bidId uuid.UUID) error {
	bid, err := s.repo.Bid.GetById(bidId)
	if err != nil {
		return serviceError(errors.New("error: Bid not found or access denied"), codes.NotFound)
	}

	if bid.ContractorId != contractorId {
		return serviceError(errors.New("error: Bid not found or access denied"), codes.NotFound)
	}

	if err := s.repo.Bid.Delete(bidId); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *bidService) AwardBid(clientId, tenderId, bidId uuid.UUID) error {
	tender, err := s.repo.Tender.GetById(tenderId)
	if err != nil {
		return serviceError(errors.New("error: Tender not found or access denied"), codes.NotFound)
	}

	if tender.ClientId != clientId {
		return serviceError(errors.New("error: Tender not found or access denied"), codes.NotFound)
	}

	bid, err := s.repo.Bid.GetById(bidId)
	if err != nil {
		return serviceError(err, codes.Internal)
	}

	if err = s.repo.Tender.Update(models.UpdateTender{
		Id:          tenderId,
		ClientId:    tender.ClientId,
		Title:       tender.Title,
		Description: tender.Description,
		Deadline:    tender.Deadline,
		Budget:      tender.Budget,
		File:        tender.File,
		Status:      config.TenderStatusAwarded,
	}); err != nil {
		return serviceError(err, codes.Internal)
	}

	if err = s.repo.Bid.Update(models.UpdateBid{
		Id:           bidId,
		ContractorId: bid.ContractorId,
		TenderId:     tenderId,
		Price:        bid.Price,
		DeliveryTime: bid.DeliveryTime,
		Comment:      bid.Comment,
		Status:       config.BidStatusAwarded,
	}); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}