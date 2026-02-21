package postgres

import (
	"context"
	"errors"

	"github.com/ZyoGo/ayo-indonesia-footbal/internal/match/domain"
	"github.com/ZyoGo/ayo-indonesia-footbal/pkg/derrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type matchRepository struct {
	db *pgxpool.Pool
}

func NewMatchRepository(db *pgxpool.Pool) domain.MatchRepository {
	return &matchRepository{db: db}
}

func (r *matchRepository) Create(ctx context.Context, match *domain.Match) error {
	_, err := r.db.Exec(ctx, queryInsertMatch,
		match.ID,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.MatchTime,
		match.Stadium,
		match.CreatedAt,
		match.UpdatedAt,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to insert match")
	}
	return nil
}

func (r *matchRepository) FindByID(ctx context.Context, id string) (*domain.Match, error) {
	var match domain.Match
	err := r.db.QueryRow(ctx, queryFindMatchByID, id).Scan(
		&match.ID,
		&match.HomeTeamID,
		&match.AwayTeamID,
		&match.MatchDate,
		&match.MatchTime,
		&match.Stadium,
		&match.HomeTeamName,
		&match.AwayTeamName,
		&match.CreatedAt,
		&match.UpdatedAt,
		&match.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, derrors.WrapErrorf(domain.ErrMatchNotFound, derrors.ErrorCodeNotFound, "%s", domain.ErrMatchNotFound.Error())
		}
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to find match")
	}
	return &match, nil
}

func (r *matchRepository) FindAll(ctx context.Context) ([]domain.Match, error) {
	rows, err := r.db.Query(ctx, queryFindAllMatches)
	if err != nil {
		return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to query matches")
	}
	defer rows.Close()

	var matches []domain.Match
	for rows.Next() {
		var match domain.Match
		if err := rows.Scan(
			&match.ID,
			&match.HomeTeamID,
			&match.AwayTeamID,
			&match.MatchDate,
			&match.MatchTime,
			&match.Stadium,
			&match.HomeTeamName,
			&match.AwayTeamName,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.DeletedAt,
		); err != nil {
			return nil, derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to scan match row")
		}
		matches = append(matches, match)
	}

	return matches, nil
}

func (r *matchRepository) Update(ctx context.Context, match *domain.Match) error {
	_, err := r.db.Exec(ctx, queryUpdateMatch,
		match.HomeTeamID,
		match.AwayTeamID,
		match.MatchDate,
		match.MatchTime,
		match.Stadium,
		match.UpdatedAt,
		match.ID,
	)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to update match")
	}
	return nil
}

func (r *matchRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	// Soft delete goals associated with the result of this match
	if _, err := tx.Exec(ctx, querySoftDeleteGoalsByMatchID, id); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to soft delete goals")
	}

	// Soft delete the match result
	if _, err := tx.Exec(ctx, querySoftDeleteResultByMatchID, id); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to soft delete match result")
	}

	// Soft delete the match itself
	if _, err := tx.Exec(ctx, queryDeleteMatch, id); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to soft delete match")
	}

	if err := tx.Commit(ctx); err != nil {
		return derrors.WrapErrorf(err, derrors.ErrorCodeInternal, "failed to commit transaction")
	}

	return nil
}
