package service

import (
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/logger"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type userService struct {
	repo   *repository.Repository
	logger *logger.Logger
}

func NewUserService(repo *repository.Repository, logger *logger.Logger) *userService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) CreateUser(request models.CreateUser) (uuid.UUID, error) {
	id, err := s.repo.User.Create(request)
	if err != nil {
		return uuid.Nil, serviceError(err, codes.Internal)
	}

	return id, nil
}

func (s *userService) GetUsers(filter models.UserFilter) ([]models.User, int, error) {
	users, total, err := s.repo.User.GetList(filter)
	if err != nil {
		return nil, 0, serviceError(err, codes.Internal)
	}

	return users, total, nil
}

func (s *userService) GetUser(id uuid.UUID) (models.User, error) {
	user, err := s.repo.User.GetById(id)
	if err != nil {
		return models.User{}, serviceError(err, codes.Internal)
	}

	return user, nil
}

func (s *userService) UpdateUser(request models.UpdateUser) error {
	if err := s.repo.User.Update(request); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	if err := s.repo.User.Delete(id); err != nil {
		return serviceError(err, codes.Internal)
	}

	return nil
}
