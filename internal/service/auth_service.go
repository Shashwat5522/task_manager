package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vedologic/task-manager/internal/domain"
	"github.com/vedologic/task-manager/internal/dto"
	"github.com/vedologic/task-manager/internal/repository"
	"github.com/vedologic/task-manager/pkg/utils"
)

// authService implements AuthService interface with business logic
type authService struct {
	userRepo       repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiryHours int) AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register registers a new user
func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	user := &domain.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := utils.GenerateToken(fmt.Sprintf("%d", user.ID), user.Email, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    fmt.Sprintf("%d", user.ID),
			Email: user.Email,
		},
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := utils.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(fmt.Sprintf("%d", user.ID), user.Email, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    fmt.Sprintf("%d", user.ID),
			Email: user.Email,
		},
	}, nil
}
