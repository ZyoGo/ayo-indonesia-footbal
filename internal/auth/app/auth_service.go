package app

import (
	"context"
	"errors"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/auth/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

const tokenExpiry = 24 * time.Hour

// AuthServicePort defines the auth service interface.
type AuthServicePort interface {
	Login(ctx context.Context, username, password string) (string, error)
}

// AuthService handles authentication logic.
type AuthService struct {
	userRepo   domain.UserRepository
	jwtService *jwt.Service
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo domain.UserRepository, jwtService *jwt.Service) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

// Login authenticates a user and returns a JWT token.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" {
		return "", derrors.NewErrorf(derrors.ErrorCodeBadRequest, "username is required")
	}
	if password == "" {
		return "", derrors.NewErrorf(derrors.ErrorCodeBadRequest, "password is required")
	}

	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", derrors.NewErrorf(derrors.ErrorCodeUnauthorized, "invalid credentials")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", derrors.NewErrorf(derrors.ErrorCodeUnauthorized, "invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(jwt.JwtAttr{Email: user.Username}, tokenExpiry)
	if err != nil {
		return "", derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to generate token")
	}

	return token, nil
}
