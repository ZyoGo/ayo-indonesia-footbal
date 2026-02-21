package domain

import "context"

// TeamRepository defines the port for team persistence.
type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	FindByID(ctx context.Context, id string) (*Team, error)
	FindAll(ctx context.Context) ([]Team, error)
	Update(ctx context.Context, team *Team) error
	SoftDelete(ctx context.Context, id string) error
	ExistsByName(ctx context.Context, name string, excludeID string) (bool, error)
}

// PlayerRepository defines the port for player persistence.
type PlayerRepository interface {
	Create(ctx context.Context, player *Player) error
	FindByID(ctx context.Context, id string) (*Player, error)
	FindByTeamID(ctx context.Context, teamID string) ([]Player, error)
	Update(ctx context.Context, player *Player) error
	SoftDelete(ctx context.Context, id string) error
	IsJerseyNumberTaken(ctx context.Context, teamID string, jerseyNumber int, excludePlayerID string) (bool, error)
}
