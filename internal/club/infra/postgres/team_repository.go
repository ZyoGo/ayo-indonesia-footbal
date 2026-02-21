package postgres

import (
	"context"
	"errors"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/club/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type teamRepository struct {
	db *pgxpool.Pool
}

func NewTeamRepository(db *pgxpool.Pool) domain.TeamRepository {
	return &teamRepository{db: db}
}

func (r *teamRepository) Create(ctx context.Context, team *domain.Team) error {
	_, err := r.db.Exec(ctx, queryInsertTeam,
		team.ID,
		team.Name,
		team.LogoURL,
		team.YearFounded,
		team.Address,
		team.City,
		team.CreatedAt,
		team.UpdatedAt,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to insert team")
	}
	return nil
}

func (r *teamRepository) FindByID(ctx context.Context, id string) (*domain.Team, error) {
	var team domain.Team
	err := r.db.QueryRow(ctx, queryFindTeamByID, id).Scan(
		&team.ID,
		&team.Name,
		&team.LogoURL,
		&team.YearFounded,
		&team.Address,
		&team.City,
		&team.CreatedAt,
		&team.UpdatedAt,
		&team.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, derrors.WrapErrorf(domain.ErrTeamNotFound, derrors.ErrorCodeNotFound, "%s", domain.ErrTeamNotFound.Error())
		}
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to find team")
	}
	return &team, nil
}

func (r *teamRepository) FindAll(ctx context.Context) ([]domain.Team, error) {
	rows, err := r.db.Query(ctx, queryFindAllTeams)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query teams")
	}
	defer rows.Close()

	var teams []domain.Team
	for rows.Next() {
		var team domain.Team
		if err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.LogoURL,
			&team.YearFounded,
			&team.Address,
			&team.City,
			&team.CreatedAt,
			&team.UpdatedAt,
			&team.DeletedAt,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan team row")
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *teamRepository) Update(ctx context.Context, team *domain.Team) error {
	_, err := r.db.Exec(ctx, queryUpdateTeam,
		team.Name,
		team.LogoURL,
		team.YearFounded,
		team.Address,
		team.City,
		team.UpdatedAt,
		team.ID,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to update team")
	}
	return nil
}

func (r *teamRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, querySoftDeleteTeam, id)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to soft delete team")
	}
	return nil
}
