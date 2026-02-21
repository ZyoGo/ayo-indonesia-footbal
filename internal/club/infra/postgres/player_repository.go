package postgres

import (
	"context"
	"errors"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type playerRepository struct {
	db *pgxpool.Pool
}

func NewPlayerRepository(db *pgxpool.Pool) domain.PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) Create(ctx context.Context, player *domain.Player) error {
	_, err := r.db.Exec(ctx, queryInsertPlayer,
		player.ID,
		player.TeamID,
		player.Name,
		player.Height,
		player.Weight,
		player.Position.String(),
		player.JerseyNumber,
		player.CreatedAt,
		player.UpdatedAt,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to insert player")
	}
	return nil
}

func (r *playerRepository) FindByID(ctx context.Context, id string) (*domain.Player, error) {
	var player domain.Player
	var position string
	err := r.db.QueryRow(ctx, queryFindPlayerByID, id).Scan(
		&player.ID,
		&player.TeamID,
		&player.Name,
		&player.Height,
		&player.Weight,
		&position,
		&player.JerseyNumber,
		&player.CreatedAt,
		&player.UpdatedAt,
		&player.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, derrors.WrapErrorf(domain.ErrPlayerNotFound, derrors.ErrorCodeNotFound, "%s", domain.ErrPlayerNotFound.Error())
		}
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to find player")
	}
	pos, _ := domain.ParsePosition(position)
	player.Position = pos
	return &player, nil
}

func (r *playerRepository) FindByTeamID(ctx context.Context, teamID string) ([]domain.Player, error) {
	rows, err := r.db.Query(ctx, queryFindPlayersByTeamID, teamID)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query players by team")
	}
	defer rows.Close()

	var players []domain.Player
	for rows.Next() {
		var player domain.Player
		var position string
		if err := rows.Scan(
			&player.ID,
			&player.TeamID,
			&player.Name,
			&player.Height,
			&player.Weight,
			&position,
			&player.JerseyNumber,
			&player.CreatedAt,
			&player.UpdatedAt,
			&player.DeletedAt,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan player row")
		}
		pos, _ := domain.ParsePosition(position)
		player.Position = pos
		players = append(players, player)
	}

	return players, nil
}

func (r *playerRepository) Update(ctx context.Context, player *domain.Player) error {
	_, err := r.db.Exec(ctx, queryUpdatePlayer,
		player.Name,
		player.Height,
		player.Weight,
		player.Position.String(),
		player.JerseyNumber,
		player.UpdatedAt,
		player.ID,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to update player")
	}
	return nil
}

func (r *playerRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, querySoftDeletePlayer, id)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to soft delete player")
	}
	return nil
}

func (r *playerRepository) IsJerseyNumberTaken(ctx context.Context, teamID string, jerseyNumber int, excludePlayerID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, queryIsJerseyNumberTaken, teamID, jerseyNumber, excludePlayerID).Scan(&exists)
	if err != nil {
		return false, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to check jersey number uniqueness")
	}
	return exists, nil
}
