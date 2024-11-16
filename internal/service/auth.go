package service

import (
	"database/sql"
	"errors"
	"tender-bridge/config"
	"tender-bridge/internal/models"
	"tender-bridge/internal/repository"
	"tender-bridge/pkg/helper"
	"tender-bridge/pkg/logger"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

type authService struct {
	repo   *repository.Repository
	logger *logger.Logger
	cfg    *config.Config
}

func NewAuthService(repo *repository.Repository, logger *logger.Logger, cfg *config.Config) *authService {
	return &authService{
		repo:   repo,
		logger: logger,
		cfg:    cfg,
	}
}

type jwtCustomClaim struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	Type   string    `json:"type"`
}

func (s *authService) CreateToken(user models.User, tokenType string, expiresAt time.Time) (*models.Token, error) {
	claims := &jwtCustomClaim{
		UserId: user.Id,
		Role:   user.Role,
		Type:   tokenType,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiresAt.Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return nil, serviceError(err, codes.Internal)
	}

	return &models.Token{
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GenerateTokens(user models.User) (*models.Token, *models.Token, error) {
	accessExpiresAt := time.Now().Add(time.Duration(s.cfg.JWTAccessExpirationHours) * time.Hour)
	refreshExpiresAt := time.Now().Add(time.Duration(s.cfg.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := s.CreateToken(user, config.TokenTypeAccess, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := s.CreateToken(user, config.TokenTypeRefresh, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) ParseToken(token string) (*jwtCustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(config.GetConfig().JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}

func (s *authService) Login(request models.Login) (*models.Token, *models.Token, error) {
	user, err := s.repo.User.GetByUsername(request.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, serviceError(errors.New("User not found"), codes.NotFound)
		}
		return nil, nil, serviceError(err, codes.Internal)
	}

	hashPassword, err := helper.GenerateHash(request.Password)
	if err != nil {
		return nil, nil, serviceError(err, codes.Internal)
	}

	if user.Password != hashPassword {
		return nil, nil, serviceError(errors.New("error: Invalid username or password"), codes.Unauthenticated)
	}

	return s.GenerateTokens(user)
}

func (s *authService) Register(request models.Register) (*models.Token, *models.Token, error) {
	// Check if the email already exists
	_, err := s.repo.User.GetByEmail(request.Email) // Ensure GetByEmail belongs to s.repo.User
	if err == nil {
		return nil, nil, serviceError(errors.New("this Email already exists"), codes.InvalidArgument)
	} else if err != sql.ErrNoRows {
		return nil, nil, serviceError(err, codes.Internal)
	}

	// Check if the username already exists
	_, err = s.repo.User.GetByUsername(request.Username)
	if err == nil {
		return nil, nil, serviceError(errors.New("username already exists"), codes.InvalidArgument)
	} else if err != sql.ErrNoRows {
		return nil, nil, serviceError(err, codes.Internal)
	}

	// Hash the password
	request.Password, err = helper.GenerateHash(request.Password)
	if err != nil {
		return nil, nil, serviceError(err, codes.Internal)
	}

	// Validate role
	if request.Role != config.RoleClient && request.Role != config.RoleContractor {
		return nil, nil, serviceError(errors.New("invalid role"), codes.InvalidArgument)
	}

	// Create user
	userId, err := s.repo.User.Create(models.CreateUser(request))
	if err != nil {
		return nil, nil, serviceError(err, codes.Internal)
	}

	// Generate tokens for the new user
	return s.GenerateTokens(models.User{
		Id:       userId,
		Role:     request.Role,
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})
}
