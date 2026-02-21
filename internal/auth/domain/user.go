package domain

import (
	"context"
	"time"

	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
)

var ErrUserNotFound = derrors.NewErrorf(derrors.ErrorCodeNotFound, "user not found")

// User represents an authenticated user (admin).
type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

// UserRepository defines the persistence interface for users.
type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
}
