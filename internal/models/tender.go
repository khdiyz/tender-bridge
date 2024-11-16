package models

import (
	"time"

	"github.com/google/uuid"
)

type Tender struct {
	Id          uuid.UUID `json:"id"`
	ClientId    uuid.UUID `json:"client_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      int64     `json:"budget"`
	File        string    `json:"file"`
	Status      string    `json:"status"`
}

type CreateTender struct {
	ClientId    uuid.UUID `json:"-"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Deadline    string    `json:"deadline" validate:"required"`
	Budget      int64     `json:"budget"`
	File        string    `json:"file"`
	Status      string    `json:"-"`
}

type UpdateTender struct {
	Id          uuid.UUID `json:"-"`
	ClientId    uuid.UUID `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	Budget      int64     `json:"budget"`
	File        string    `json:"file"`
	Status      string    `json:"status"`
}

type UpdateTenderStatus struct {
	Id     uuid.UUID `json:"-"`
	Status string    `json:"status"`
}

type TenderFilter struct {
	Search string
	Limit  int
	Offset int
}
