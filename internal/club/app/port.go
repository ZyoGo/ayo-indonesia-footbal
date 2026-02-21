package app

import (
	"context"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
)

//go:generate mockgen -source=port.go -destination=mock/mock_service.go -package=mock

// TeamServicePort defines the contract for team business operations.
type TeamServicePort interface {
	Create(ctx context.Context, team *domain.Team) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	GetAll(ctx context.Context) ([]domain.Team, error)
	Update(ctx context.Context, id string, team *domain.Team) error
	Delete(ctx context.Context, id string) error
}

// PlayerServicePort defines the contract for player business operations.
type PlayerServicePort interface {
	Create(ctx context.Context, player *domain.Player) (string, error)
	GetByID(ctx context.Context, id string) (*domain.Player, error)
	GetByTeamID(ctx context.Context, teamID string) ([]domain.Player, error)
	Update(ctx context.Context, id string, player *domain.Player) error
	Delete(ctx context.Context, id string) error
}
