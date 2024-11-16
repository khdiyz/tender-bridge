package service

import (
	"tender-bridge/config"
	"tender-bridge/internal/repository"
	"tender-bridge/internal/storage"
	"tender-bridge/pkg/logger"
)

type Service struct {
}

func NewService(repos *repository.Repository, storage *storage.Storage, cfg *config.Config, loggers *logger.Logger) *Service {
	return &Service{
	}
}
