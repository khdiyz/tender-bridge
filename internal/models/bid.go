package models

import (
	"github.com/google/uuid"
)

type Bid struct {
	Id           uuid.UUID `json:"id"`
	ContractorId uuid.UUID `json:"contractor_id"`
	TenderId     uuid.UUID `json:"tender_id"`
	Price        int64     `json:"price"`
	DeliveryTime int       `json:"delivery_time"`
	Comment      string    `json:"comments"`
	Status       string    `json:"status"`
}

type CreateBid struct {
	ContractorId uuid.UUID `json:"-"`
	TenderId     uuid.UUID `json:"-"`
	Price        int64     `json:"price"`
	DeliveryTime int       `json:"delivery_time"`
	Comment      string    `json:"comments"`
	Status       string    `json:"-"`
}

type UpdateBid struct {
	Id           uuid.UUID `json:"-"`
	ContractorId uuid.UUID `json:"-"`
	TenderId     uuid.UUID `json:"-"`
	Price        int64     `json:"price"`
	DeliveryTime int       `json:"delivery_time"`
	Comment      string    `json:"comments"`
	Status       string    `json:"status"`
}

type BidFilter struct {
	Search       string
	FromPrice    int64
	ToPrice      int64
	TenderId     uuid.UUID
	ContractorId uuid.UUID
	Limit        int
	Offset       int
}

type BidNotification struct {
	OrderID string `json:"tender_id"`
	BidID   string `json:"bid_id"`
	Message string `json:"message"`
}
