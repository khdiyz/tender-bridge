package service

import (
	"database/sql"
	"errors"
	"tender-bridge/config"
	"tender-bridge/internal/cache"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/internal/ws"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

var (
	errBidNotFound    = errors.New("error: Bid not found or access denied")
	errTenderNotFound = errors.New("error: Tender not found or access denied")
)

type bidService struct {
	repo   *repository.Repository
	cache  *cache.RedisCache
	logger *logger.Logger
}

func NewBidService(repo *repository.Repository, cache *cache.RedisCache, logger *logger.Logger) *bidService {
	return &bidService{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *bidService) SubmitBid(request models.CreateBid) (uuid.UUID, error) {
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

	go func() {
		ws.BroadcastNotification(tender.ClientId.String(), "Submitted new bid")
	}()

	return id, nil
}

func (s *bidService) GetBids(filter models.BidFilter) ([]models.Bid, int, error) {
	bids, total, err := s.repo.Bid.GetList(filter)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	tenderIds := make([]uuid.UUID, len(bids))
	for i := range bids {
		tenderIds[i] = bids[i].TenderId
	}

	tenders, err := s.repo.Tender.GetByIds(tenderIds)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	tendersMap := make(map[uuid.UUID]models.Tender, len(tenders))
	for i := range tenders {
		tendersMap[tenders[i].Id] = tenders[i]
	}

	for i := range bids {
		bids[i].Tender = tendersMap[bids[i].TenderId]
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
		return serviceError(errBidNotFound, codes.NotFound)
	}

	if bid.ContractorId != contractorId {
		return serviceError(errBidNotFound, codes.NotFound)
	}

	if err := s.repo.Bid.Delete(bidId); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *bidService) AwardBid(clientId, tenderId, bidId uuid.UUID) error {
	tender, err := s.repo.Tender.GetById(tenderId)
	if err != nil {
		return serviceError(errTenderNotFound, codes.NotFound)
	}

	if tender.ClientId != clientId {
		return serviceError(errTenderNotFound, codes.NotFound)
	}

	if tender.Status != config.TenderStatusOpen {
		return serviceError(errors.New("the tender is not open"), codes.InvalidArgument)
	}

	bid, err := s.repo.Bid.GetById(bidId)
	if err != nil {
		return serviceError(err, codes.Internal)
	}

	if bid.Status != config.BidStatusPending {
		return serviceError(errors.New("the bid is not pending"), codes.InvalidArgument)
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

	go func() {
		ws.BroadcastNotification(bid.ContractorId.String(), "Your bid awarded")
	}()

	return nil
}
